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
MODEL_NAME = os.getenv("OLLAMA_MODEL", "qwen2.5:14b") # 我們切換回 Llama 3

# --- Pydantic 模型 (RAG 4.1：包含列表) ---
class DiseaseRisk(BaseModel):
    disease: str
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
**重要提醒**：
- 正常值（例如 Glu-AC: 89 mg/dL < 100，HbA1c: 4.1% < 5.7%）不應標記為風險
- 所有文字欄位都必須有內容，不得為空字串
- 除了欄位名稱外，其餘內容都必須使用繁體中文，特別是advice的內容確保是繁體中文
- 必須嚴格按照「相關醫學知識」中的數值範圍進行判斷
- 總結分析和個人詳細分析都必需嚴格執行字數要求

你是一位AI健康分析專家。請嚴格根據以下提供的「健康檢查數據」和「相關醫學知識」生成 JSON 格式分析報告。

=== 健康檢查數據 ===
{query_text}

=== 相關醫學知識 ===
{full_context}

=== 核心判斷規則（嚴格遵守）===
1. **數值範圍判斷（最重要）**：
   - 必須將每個指標的**實際數值**與「相關醫學知識」中的**正常範圍**進行嚴格比對
   - 例如：Glu-AC 正常 < 100 mg/dL，如果實際值是 89 mg/dL → **正常，不應標記風險**
   - 例如：HbA1c 正常 < 5.7%，如果實際值是 4.1% → **正常，不應標記風險**
   - 只有當數值**超出正常範圍**時，才標記為風險

2. **disease_risks 列表規則**：
   - **只包含異常指標**：只有當指標數值超出正常範圍時，才在 disease_risks 中列出
   - **不包含正常值**：如果所有指標都在正常範圍內，disease_risks 應該是空列表 []
   - **不包含空項目**：不要生成 disease 為空字串的項目

3. **健康分數計算**：
   - 如果所有指標正常 → health_score 應該在 80-100 分
   - 如果有輕微異常 → health_score 在 60-80 分
   - 如果有嚴重異常 → health_score 在 0-60 分
   - 參考「核心評分準則」進行計算

4. **內容完整性要求**：
   - summary：必須提供完整的健康總結（至少 50 字），說明整體健康狀況
   - personal_analysis：必須提供詳細的指標分析（至少 100 字），逐一分析異常指標
   - disease_risks 中的每個項目：
     * disease：必須是具體的疾病名稱（繁體中文）
     * main_factors：必須列出觸發此風險的具體指標和數值（例如 ["CRE: 6.0 mg/dL"]）
     * advice：必須提供具體的醫療或生活建議（至少 30 字）
   - insurance_recommendation：必須提供至少 2-3 項具體的保險建議
   - protection_plan：必須提供至少 3-5 項具體的健康防護計畫

=== 語言要求（嚴格遵守）===
- **JSON 欄位名稱**：必須使用英文（如 "summary", "personal_analysis", "disease", "advice" 等）
- **所有內容值**：必須使用繁體中文（Traditional Chinese）
  - 不得使用英文、簡體中文或其他語言
  - 醫學名詞必須使用繁體中文

=== 任務步驟 ===
1. 逐一檢查每個指標的數值是否在正常範圍內
2. 只將**超出正常範圍**的指標標記為風險
3. 根據異常指標的嚴重程度計算 health_score
4. 生成完整的 summary 和 personal_analysis（不得為空）
5. 只列出真正異常的疾病風險（正常值不應出現在 disease_risks 中）
6. 為每個風險提供具體的建議和防護計畫

=== 輸出格式（嚴格遵守，僅輸出 JSON，所有欄位都必須有內容）===
{{
  "health_score": <0-100 分數>,
  "summary": "<完整的健康總結，至少 50 字，繁體中文>",
  "personal_analysis": "<詳細的指標分析，至少 100 字，繁體中文>",
  "disease_risks": [
    {{
      "disease": "<具體疾病名稱，繁體中文>",  
      "risk_level": "<高/中/低風險，繁體中文>",
      "main_factors": ["<例如 'CRE: 6.0 mg/dL' 等>"], 
      "advice": "<具體建議，至少 30 字，繁體中文>"
    }}
  ],
  "insurance_recommendation": ["<至少 2-3 項具體建議，繁體中文>"],
  "protection_plan": ["<至少 3-5 項具體計畫，繁體中文>"],
  "success": true
}}

**重要提醒**：
- 正常值（例如 Glu-AC: 89 mg/dL < 100，HbA1c: 4.1% < 5.7%）不應標記為風險
- 所有文字欄位都必須有內容，不得為空字串
- 除了欄位名稱外，其餘內容都必須使用繁體中文，特別是advice的內容確保是繁體中文
- 必須嚴格按照「相關醫學知識」中的數值範圍進行判斷
- 總結分析和個人詳細分析都必需嚴格執行字數要求
"""
            result = self._call_ollama_json(prompt, HealthAnalysis)
            if not result:
                raise ValueError("LLM 回傳無效格式")

            disease_risks_proto = [
                data_pb2.DiseaseRisk(
                    disease=dr.disease,
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

=== 核心判斷規則（嚴格遵守）===
1. **數值範圍判斷（最重要）**：
   - 必須將每個指標的**實際數值**與「相關醫學知識」中的**正常範圍**進行嚴格比對
   - 例如：Glu-AC 正常 < 100 mg/dL，如果實際值是 89 mg/dL → **正常，不應標記風險**
   - 例如：HbA1c 正常 < 5.7%，如果實際值是 4.1% → **正常，不應標記風險**
   - 只有當數值**超出正常範圍**時，才標記為風險

2. **disease_risk_evaluation 列表規則**：
   - **只包含異常指標**：只有當指標數值超出正常範圍時，才在 disease_risk_evaluation 中列出
   - **不包含正常值**：如果所有指標都在正常範圍內，disease_risk_evaluation 應該是空列表 []
   - **不包含空項目**：不要生成 disease 為空字串的項目

3. **風險分數計算**：
   - 如果所有指標正常 → overall_risk_score 應該在 0-20 分（低風險）
   - 如果有輕微異常 → overall_risk_score 在 20-50 分（中風險）
   - 如果有嚴重異常 → overall_risk_score 在 50-100 分（高風險）
   - 參考「核心評分準則」進行計算

4. **內容完整性要求**：
   - summary：必須提供完整的核保摘要（至少 50 字），說明整體風險狀況
   - risk_level_label：必須是「高風險」、「中風險」或「低風險」（繁體中文）
   - core_recommendation：必須提供至少 2-3 項具體的核保建議（例如「標準承保」、「加收保費」、「部分排除」等）
   - disease_risk_evaluation 中的每個項目：
     * disease：必須是具體的疾病名稱（繁體中文）
     * main_factors：必須列出觸發此風險的具體指標和數值（例如 ["LDL-C: 900 mg/dL"]）
     * advice：必須提供具體的核保建議（至少 30 字，例如「排除腎臟相關理賠」）

=== 語言要求（嚴格遵守）===
- **JSON 欄位名稱**：必須使用英文（如 "summary", "risk_level_label", "disease", "advice" 等）
- **所有內容值**：必須使用繁體中文（Traditional Chinese）
  - 不得使用英文、簡體中文或其他語言
  - 醫學名詞必須使用繁體中文

=== 任務步驟 ===
1. 逐一檢查每個指標的數值是否在正常範圍內
2. 只將**超出正常範圍**的指標標記為風險
3. 根據異常指標的嚴重程度計算 overall_risk_score
4. 生成完整的 summary（不得為空）
5. 只列出真正異常的疾病風險（正常值不應出現在 disease_risk_evaluation 中）
6. 為每個風險提供具體的核保建議

=== 輸出格式（嚴格遵守，僅輸出 JSON，所有欄位都必須有內容）===
{{
  "overall_risk_score": <0-100 分數>,
  "risk_level_label": "<高/中/低風險，繁體中文>",
  "summary": "<完整的核保摘要，至少 50 字，繁體中文>",
  "core_recommendation": ["<至少 2-3 項具體核保建議，繁體中文>"],
  "disease_risk_evaluation": [
    {{
      "disease": "<具體疾病名稱，繁體中文>", 
      "risk_score": <風險分數>, 
      "risk_level": "<高/中/低風險，繁體中文>", 
      "main_factors": ["<具體指標和數值，例如 'LDL-C: 900 mg/dL'>"], 
      "advice": "<具體核保建議，至少 30 字，繁體中文>"
    }}
  ],
  "success": true
}}

**重要提醒**：
- 正常值（例如 Glu-AC: 89 mg/dL < 100，HbA1c: 4.1% < 5.7%）不應標記為風險
- 所有文字欄位都必須有內容，不得為空字串
- disease_risk_evaluation 中不應包含空項目
- 必須嚴格按照「相關醫學知識」中的數值範圍進行判斷
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