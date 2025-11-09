import grpc
from concurrent import futures
import logging
import json
import re
import os
import socket
from pydantic import BaseModel, ValidationError
from typing import List, Optional, Dict, Set
from ollama import Client
import data_pb2, data_pb2_grpc
import time

# 配置日誌
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# --- 環境配置 (使用 Llama 3) ---
OLLAMA_HOST = os.getenv("OLLAMA_HOST", "http://localhost:11434")
MODEL_NAME = os.getenv("OLLAMA_MODEL", "llama3:8b") # 我們切換回 Llama 3

# --- Pydantic 模型 (RAG 4.1：包含列表) ---
class DiseaseRisk(BaseModel):
    disease: str
    risk_percent: int
    risk_level: str
    main_factors: List[str]
    advice: str

class HealthAnalysis(BaseModel):
    health_score: int
    summary: str
    personal_analysis: str
    disease_risks: List[DiseaseRisk] # <-- 恢復列表
    insurance_recommendation: List[str]
    protection_plan: List[str]
    success: bool

class DiseaseRiskForInsurer(BaseModel):
    disease: str
    risk_score: int
    risk_level: str
    main_factors: List[str]
    advice: str

class InsurerAnalysis(BaseModel):
    overall_risk_score: int
    risk_level_label: str
    summary: str
    core_recommendation: List[str]
    disease_risk_evaluation: List[DiseaseRiskForInsurer] # <-- 恢復列表
    success: bool

class HealthAnalysisServicer(data_pb2_grpc.HealthServiceServicer):
    
    def _load_health_rules(self) -> (str, Dict[str, str]):
        """
        在服務啟動時一次性載入規則，並將其切割為 RAG 4.0 所需的知識塊。
        """
        rules_file = "health_rules.txt"
        if not os.path.exists(rules_file):
            logger.error(f"關鍵錯誤：知識文件 '{rules_file}' 不存在！")
            return "無可用的醫學參考資料。", {}
            
        try:
            with open(rules_file, 'r', encoding='utf-8') as f:
                full_text = f.read()
            logger.info(f"成功載入 '{rules_file}' (大小: {len(full_text)} 字節)")
            
            chunks = re.split(r'(### .*\n)', full_text)
            knowledge_chunks: Dict[str, str] = {}
            
            if not chunks[0].startswith("### "):
                knowledge_chunks["General_Header"] = chunks[0].strip()
                chunks = chunks[1:]

            for i in range(0, len(chunks), 2):
                if i + 1 < len(chunks):
                    header = chunks[i].strip()
                    content = chunks[i+1].strip()
                    knowledge_chunks[header] = content
                
            logger.info(f"(RAG 4.0) 知識庫已切割為 {len(knowledge_chunks)} 個主題塊。")
            return full_text, knowledge_chunks
            
        except Exception as e:
            logger.error(f"載入 '{rules_file}' 失敗: {str(e)}")
            return "讀取醫學參考資料失敗。", {}

    def _get_key_to_topic_map(self) -> Dict[str, List[str]]:
        """
        定義健檢指標 (Key) 應觸發哪個知識主題 (Topic)。
        """
        return {
            "血糖指標": ["Glu-AC", "HbA1c", "Glu-PC"],
            "肝功能指標": ["Alb", "TP", "AST", "ALT", "D-Bil", "ALP", "T-Bil"],
            "腎功能指標": ["UN", "CRE", "U.A", "Alb/CRE Ratio"],
            "血脂指標": ["T-CHO", "LDL-C", "HDL-C", "TG"],
            "全血球計數": ["Hb", "Hct", "PLT", "WBC", "RBC", "MCV", "MCH", "MCHC", "RDW-CV"],
            "凝血功能": ["PT", "aPTT"],
            "發炎指標": ["hsCRP", "ESR"],
            "癌症指標": ["AFP", "CEA", "CA-125", "CA19-9"],
            "尿液常規檢查": ["Specific Gravity", "PH", "Protein (Dipstick)", "Glucose (Dipstick)", 
                         "Bilirubin (Dipstick)", "Urobilinogen (Dipstick)", "RBC (Urine)", 
                         "WBC (Urine)", "Epithelial Cells", "Casts", "Ketone", 
                         "Crystal", "Bacteria", "Albumin (Dipstick)", "Creatinine (Dipstick)", 
                         "Nitrite", "Occult Blood", "WBC Esterase"],
            "其他": ["BP"],
            "核心評分準則": ["*ALWAYS_INCLUDE*"]
        }

    def _filter_context_for_rag(self, test_results: dict) -> str:
        """
        RAG 4.0 核心：根據輸入的指標，動態篩選相關的知識塊。
        """
        input_keys_raw = set(test_results.keys())
        relevant_topics: Set[str] = set() 
        topic_map = self._get_key_to_topic_map()
        
        for topic, indicators in topic_map.items():
            if "*ALWAYS_INCLUDE*" in indicators:
                relevant_topics.add(topic)
                continue
                
            for key in input_keys_raw:
                clean_key = key.split('(')[0].strip() 
                if clean_key in indicators:
                    relevant_topics.add(topic)
                    break 

        logger.info(f"(RAG 4.0) 偵測到相關主題: {relevant_topics}")

        context_parts: List[str] = []
        
        for header, content in self.knowledge_chunks.items():
            for topic in relevant_topics:
                if topic in header: 
                    context_parts.append(f"{header}\n{content}")
                    break

        if not context_parts:
            logger.warning("(RAG 4.0) 警告：未篩選到任何相關規則！")
            return "無相關醫學參考資料"
            
        final_context = "\n---\n".join(context_parts)
        logger.info(f"(RAG 4.0) 篩選後的上下文 (大小: {len(final_context)} 字節)")
        return final_context


    def __init__(self):
        try:
            logger.info("開始初始化服務器依賴")
            self.full_health_rules, self.knowledge_chunks = self._load_health_rules()
            
            if not self._check_ollama_connection():
                raise RuntimeError(f"Ollama 服務不可用: {OLLAMA_HOST}")
                
            self.ollama_client = Client(host=OLLAMA_HOST)
            logger.info(f"Ollama 客戶端初始化，目標地址: {OLLAMA_HOST}")
            
        except Exception as e:
            logger.error(f"服務器初始化失敗: {str(e)}")
            raise

    # ... ( _check_ollama_connection 和 _call_ollama_json 保持不變 ) ...
    def _check_ollama_connection(self, max_retries=3, delay=2) -> bool:
        try:
            host_port = OLLAMA_HOST.replace("http://", "").replace("https://", "")
            if ":" in host_port:
                host, port = host_port.split(":")
            else:
                host = host_port
                port = 11434
            logger.info(f"檢查 Ollama 連接: {host}:{port}")
            for attempt in range(max_retries):
                with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
                    s.settimeout(5)
                    result = s.connect_ex((host, int(port)))
                    if result == 0:
                        try:
                            import requests
                            response = requests.get(f"{OLLAMA_HOST}/api/tags", timeout=3)
                            if response.status_code == 200:
                                logger.info(f"Ollama 連線檢查成功（API 可用），嘗試次數: {attempt + 1}")
                                return True
                        except ImportError:
                            logger.info(f"Ollama 端口連接成功（跳過 API 檢查），嘗試次數: {attempt + 1}")
                            return True
                        except Exception as api_error:
                            logger.warning(f"Ollama API 檢查失敗: {str(api_error)}")
                            logger.info(f"Ollama 端口連接成功，嘗試次數: {attempt + 1}")
                            return True
                    logger.warning(f"Ollama 連線檢查失敗，嘗試 {attempt + 1}/{max_retries}，代碼 {result}")
                    if attempt < max_retries - 1:
                        time.sleep(delay)
            logger.error("Ollama 連線檢查失敗，所有重試均失敗")
            return False
        except Exception as e:
            logger.error(f"Ollama 連線檢查失敗: {str(e)}")
            return False

    def _call_ollama_json(self, prompt: str, schema: BaseModel) -> Optional[BaseModel]:
        logger.info(f"發送 Ollama prompt (前 200 字): {prompt[:200]}...")
        try:
            response = self.ollama_client.chat(
                model=MODEL_NAME,
                messages=[{"role": "user", "content": prompt}],
                format=schema.model_json_schema()
            )
            content = response['message']['content']
            logger.info(f"Ollama 回應完整內容: {content}")
            return schema.model_validate_json(content)
        except (ValidationError, json.JSONDecodeError) as e:
            logger.error(f"Ollama 解析失敗: {str(e)} - 回應內容: {content if 'content' in locals() else 'N/A'}")
            return None
        except Exception as e:
            logger.error(f"Ollama 其他錯誤: {str(e)}")
            return None

    def AnalyzeHealthReportForUser(self, request, context):
        logger.info(f"處理用戶健康報告，報告 ID: {request.report_id}")
        try:
            test_results = json.loads(request.test_results_json)
            available_keys = list(test_results.keys())
            query_text = "\n".join([f"{k}: {v}" for k, v in test_results.items()])
            
            full_context = self._filter_context_for_rag(test_results)

            # --- RAG 4.1：強化 Prompt 指令 (恢復列表) ---
            prompt = f"""
你是一位AI健康分析專家。請嚴格根據以下提供的「健康檢查數據」和「相關醫學知識」生成 JSON 格式分析報告。

=== 健康檢查數據 ===
{query_text}

=== 相關醫學知識 ===
{full_context}

=== 任務 ===
你必須分析「健康檢查數據」中的**每一個**指標，並嚴格比對「相關醫學知識」中的**所有**規則（特別是「核心評分準則」和各項指標的風險定義）。
你必須計算一個**最終的健康分數** (health_score)。
你必須在 `disease_risks` 列表中條列出**所有**偵測到的風險。

=== 輸出格式（嚴格遵守，僅輸出 JSON）===
{{
  "health_score": <一个 0-100 之間的分數 (0分代表極度危險, 100分代表非常健康)>,
  "summary": "<根據數據和知識生成的總結，必須提及所有高風險項>",
  "personal_analysis": "<對所有異常指標的詳細文字分析>",
  "disease_risks": [
    {{
      "disease": "<偵測到的疾病名稱>", 
      "risk_percent": <根據醫學知識計算出的風險百分比>, 
      "risk_level": "<高/中/低風險>",
      "main_factors": ["<必須填入觸發此風險的關鍵指標, 例如 'CRE: 6.0'>"], 
      "advice": "<具體的醫療或生活建議>"
    }}
  ],
  "insurance_recommendation": ["<推薦的保險類型>"],
  "protection_plan": ["<具體的健康防護計畫>"],
  "success": true
}}
**僅輸出 JSON！必須使用「相關醫學知識」區塊的規則來判斷風險和分數！不要幻想數據！**
"""
            result = self._call_ollama_json(prompt, HealthAnalysis)
            if not result:
                raise ValueError("LLM 回傳無效格式")

            disease_risks_proto = [
                data_pb2.DiseaseRisk(
                    disease=dr.disease,
                    risk_percent=dr.risk_percent,
                    risk_level=dr.risk_level,
                    main_factors=[f for f in dr.main_factors if f in available_keys or ':' in f], # 允許 'CRE: 6.0' 這種格式
                    advice=dr.advice
                ) for dr in result.disease_risks if dr.disease
            ]

            return data_pb2.UserHealthAnalysisResponse(
                health_score=result.health_score,
                summary=result.summary,
                personal_analysis=result.personal_analysis,
                disease_risks=disease_risks_proto, # <-- 恢復列表
                insurance_recommendation=result.insurance_recommendation,
                protection_plan=result.protection_plan,
                success=result.success
            )

        except Exception as e:
            logger.error(f"用戶分析失敗: {str(e)}")
            return data_pb2.UserHealthAnalysisResponse(
                health_score=0, summary="分析失敗", personal_analysis="請稍後重試",
                disease_risks=[], insurance_recommendation=[], protection_plan=[], success=False
            )

    def AnalyzeHealthReportForInsurer(self, request, context):
        logger.info(f"處理保險公司健康報告，報告 ID: {request.report_id}")
        
        disease_risk_proto_list = [] # (保持 Bug 修復)
        
        try:
            test_results = json.loads(request.test_results_json)
            available_keys = list(test_results.keys())
            query_text = "\n".join([f"{k}: {v}" for k, v in test_results.items()])
            
            full_context = self._filter_context_for_rag(test_results)

            # --- RAG 4.1：強化 Prompt 指令 (恢復列表) ---
            prompt = f"""
你是一位保險核保分析師。請嚴格根據以下提供的「健康檢查數據」和「相關醫學知識」生成 JSON 格式分析報告。

=== 健康檢查數據 ===
{query_text}

=== 相關醫學知識 ===
{full_context}

=== 任務 ===
你必須分析「健康檢查數據」中的**每一個**指標，並嚴格比對「相關醫學知識」中的**所有**規則（特別是「核心評分準則」）。
你必須計算一個**最終的風險分數** (overall_risk_score)。
你必須在 `disease_risk_evaluation` 列表中條列出**所有**偵測到的風險。

=== 輸出格式（嚴格遵守，僅輸出 JSON）===
{{
  "overall_risk_score": <一个 0-100 之間的分數 (0分代表最低風險, 100分代表極高風險)>,
  "risk_level_label": "<高/中/低風險>",
  "summary": "<專業的核保摘要，必須提及所有高風險項>",
  "core_recommendation": ["<具體的核保建議，例如 '拒保' 或 '加收保費'>"],
  "disease_risk_evaluation": [
    {{
      "disease": "<偵測到的疾病名稱>", 
      "risk_score": <根據醫學知識計算出的風險分數>, 
      "risk_level": "<高/中/低風險>", 
      "main_factors": ["<必須填入觸發此風險的關鍵指標, 例如 'LDL-C: 900'>"], 
      "advice": "<核保建議，如 '排除腎臟相關理賠'>"
    }}
  ],
  "success": true
}}
**僅輸出 JSON！必須使用「相關醫學知識」區塊的規則來判斷風險和分數！不要幻想數據！**
"""
            result = self._call_ollama_json(prompt, InsurerAnalysis)
            if not result:
                raise ValueError("LLM 回傳無效格式")

            disease_risk_proto_list = [
                data_pb2.DiseaseRiskForInsurer(
                    disease=dr.disease,
                    risk_score=dr.risk_score,
                    risk_level=dr.risk_level,
                    main_factors=[f for f in dr.main_factors if f in available_keys or ':' in f], 
                    advice=dr.advice
                ) for dr in result.disease_risk_evaluation if dr.disease
            ]

            return data_pb2.InsurerHealthAnalysisResponse(
                overall_risk_score=result.overall_risk_score,
                risk_level_label=result.risk_level_label,
                summary=result.summary,
                core_recommendation=result.core_recommendation,
                disease_risk_evaluation=disease_risk_proto_list, # <-- 恢復列表並修復 Bug
                success=result.success
            )

        except Exception as e:
            logger.error(f"保險分析失敗: {str(e)}")
            return data_pb2.InsurerHealthAnalysisResponse(
                overall_risk_score=100, risk_level_label="分析失敗", summary=f"分析失敗: {str(e)}",
                core_recommendation=[], disease_risk_evaluation=disease_risk_proto_list, success=False
            )

def serve():
    try:
        grpc_port = os.getenv("PYTHON_BACKEND_GRPC_PORT", "50052")
        grpc_host = os.getenv("PYTHON_BACKEND_GRPC_HOST", "[::]")
        
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
        data_pb2_grpc.add_HealthServiceServicer_to_server(HealthAnalysisServicer(), server)
        
        listen_addr = f"{grpc_host}:{grpc_port}"
        server.add_insecure_port(listen_addr)
        logger.info(f"Python gRPC 服務器啟動，監聽地址: {listen_addr}")
        logger.info(f"注意: Go Server 運行在 :50051，Python Backend 運行在 :{grpc_port}")
        
        server.start()
        try:
            while True:
                server.wait_for_termination(timeout=60)
        except KeyboardInterrupt:
            logger.info("服務器手動終止")
            server.stop(grace=None)
    except Exception as e:
        logger.error(f"服務器啟動失敗: {str(e)}")
        raise

if __name__ == "__main__":
    serve()