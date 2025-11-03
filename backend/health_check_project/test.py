import grpc
from concurrent import futures
import logging
import json
import re
import os
import socket
from pydantic import BaseModel, ValidationError
from typing import List, Optional
from ollama import Client
from langchain_chroma import Chroma
from langchain_huggingface import HuggingFaceEmbeddings
import data_pb2, data_pb2_grpc
import time

# 配置日誌
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# 環境配置
OLLAMA_HOST = "http://localhost:11434"
MODEL_NAME = "meditron:7b"

# Pydantic 模型
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
    disease_risks: List[DiseaseRisk]
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
    disease_risk_evaluation: List[DiseaseRiskForInsurer]
    success: bool

class HealthAnalysisServicer(data_pb2_grpc.HealthServiceServicer):
    def __init__(self):
        try:
            logger.info("開始初始化服務器依賴")
            chroma_dir = "D:/gg/chroma_db"
            if not os.path.exists(chroma_dir):
                os.makedirs(chroma_dir)
                logger.info(f"創建 Chroma 資料夾: {chroma_dir}")
            start_time = time.time()
            self.embedding = HuggingFaceEmbeddings(model_name="sentence-transformers/all-MiniLM-L6-v2")
            logger.info(f"嵌入模型初始化完成，耗時 {time.time() - start_time:.2f} 秒")
            start_time = time.time()
            self.vectorstore = Chroma(
                persist_directory=chroma_dir,
                embedding_function=self.embedding,
                collection_name="health_knowledge"
            )
            logger.info(f"向量資料庫加載完成，耗時 {time.time() - start_time:.2f} 秒")
            self.retriever = self.vectorstore.as_retriever(search_kwargs={"k": 2})
            logger.info("Retriever 配置完成")
            if not self._check_ollama_connection():
                raise RuntimeError(f"Ollama 服務不可用: {OLLAMA_HOST}")
            self.ollama_client = Client(host=OLLAMA_HOST)
            logger.info(f"Ollama 客戶端初始化，目標地址: {OLLAMA_HOST}")
        except Exception as e:
            logger.error(f"服務器初始化失敗: {str(e)}")
            raise

    def _check_ollama_connection(self, max_retries=3, delay=2) -> bool:
        try:
            host, port = OLLAMA_HOST.replace("http://", "").split(":")
            for attempt in range(max_retries):
                with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
                    s.settimeout(5)
                    result = s.connect_ex((host, int(port)))
                    if result == 0:
                        logger.info(f"Ollama 連線檢查成功，嘗試次數: {attempt + 1}")
                        return True
                    logger.warning(f"Ollama 連線檢查失敗，嘗試 {attempt + 1}/{max_retries}，代碼 {result}")
                    time.sleep(delay)
            logger.error("Ollama 連線檢查失敗，所有重試均失敗")
            return False
        except Exception as e:
            logger.error(f"Ollama 連線檢查失敗: {str(e)}")
            return False

    def _clean_doc_content(self, text: str) -> str:
        return re.sub(r'問題:.*\n回答:', '', text, flags=re.DOTALL).strip()

    def _fetch_context(self, test_results: dict) -> str:
        query_categories = {
            "blood_sugar": ["Glu-AC", "HbA1c", "Glu-PC"],
            "lipid": ["LDL-C", "HDL-C", "TG", "T-CHO"],
            "liver": ["ALT(GPT)", "AST(GOT)", "ALP", "T-Bil", "D-Bil"],
            "kidney": ["CRE", "UN", "Alb/CRE Ratio"],
            "general": ["Hb", "Hct", "PLT", "WBC", "RBC", "hsCRP"]
        }
        context_parts = []
        for category, keys in query_categories.items():
            query_items = [f"{k}: {test_results.get(k, 'N/A')}" for k in keys if k in test_results]
            if query_items:
                query_text = "\n".join(query_items)
                docs = self.retriever.invoke(query_text)  # 改用 invoke
                cleaned_docs = [self._clean_doc_content(doc.page_content) for doc in docs if hasattr(doc, 'page_content') and doc.page_content]
                context_parts.extend(cleaned_docs)
        context = "\n".join(context_parts) if context_parts else "無相關醫學參考資料"
        logger.info(f"Multi-Query 上下文 (前 200 字): {context[:200]}...")
        return context

    def _generate_hypothetical_doc(self, query_text: str) -> str:
        hyde_prompt = f"""根據以下健康數據生成假設性醫學摘要（繁體中文）：
{query_text}
摘要應包含關鍵指標分析與潛在風險，但不得包含具體數值。"""
        try:
            response = self.ollama_client.chat(
                model=MODEL_NAME,
                messages=[{"role": "user", "content": hyde_prompt}]
            )
            hyde_doc = response['message']['content'].strip()
            logger.info(f"HyDE 假設性文檔 (前 200 字): {hyde_doc[:200]}...")
            return hyde_doc
        except Exception as e:
            logger.error(f"HyDE 生成失敗: {str(e)}")
            return ""

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
            multi_query_context = self._fetch_context(test_results)
            hypothetical_doc = self._generate_hypothetical_doc(query_text)
            hyde_docs = self.retriever.invoke(hypothetical_doc)  # 改用 invoke
            hyde_context = "\n".join([
                self._clean_doc_content(doc.page_content)
                for doc in hyde_docs
                if hasattr(doc, 'page_content') and doc.page_content
            ])
            full_context = f"醫學參考資料:\n{multi_query_context}\n\n假設性分析:\n{hyde_context}"
            logger.info(f"完整上下文 (前 200 字): {full_context[:200]}...")

            prompt = f"""
你是一位AI健康分析專家，請根據以下資料生成 JSON 格式分析報告：
=== 健康檢查數據 ===
{query_text}

=== 相關醫學知識 ===
{full_context}

=== 輸出格式（嚴格遵守）===
{{
  "health_score": 70,
  "summary": "您的健康狀況為中度健康，但存在一些異常指標需要注意。",
  "personal_analysis": "血糖、腎功能與血脂偏高，建議進一步檢查。",
  "disease_risks": [
    {{
      "disease": "糖尿病", "risk_percent": 60, "risk_level": "中風險",
      "main_factors": ["HbA1c", "Glu-AC"], "advice": "控制飲食，定期追蹤血糖"
    }},
    {{
      "disease": "腎病", "risk_percent": 50, "risk_level": "中風險",
      "main_factors": ["CRE"], "advice": "腎臟科追蹤，減少蛋白質攝取"
    }}
  ],
  "insurance_recommendation": ["健康險", "重疾險"],
  "protection_plan": ["每3個月回診", "低糖飲食", "規律運動"],
  "success": true
}}
**僅輸出 JSON！main_factors 只能使用輸入數據中的項目！**
"""
            result = self._call_ollama_json(prompt, HealthAnalysis)
            if not result:
                raise ValueError("LLM 回傳無效格式")

            disease_risks_proto = [
                data_pb2.DiseaseRisk(
                    disease=dr.disease,
                    risk_percent=dr.risk_percent,
                    risk_level=dr.risk_level,
                    main_factors=[f for f in dr.main_factors if f in available_keys],
                    advice=dr.advice
                ) for dr in result.disease_risks if dr.disease
            ]

            return data_pb2.UserHealthAnalysisResponse(
                health_score=result.health_score,
                summary=result.summary,
                personal_analysis=result.personal_analysis,
                disease_risks=disease_risks_proto,
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
        try:
            test_results = json.loads(request.test_results_json)
            available_keys = list(test_results.keys())
            query_text = "\n".join([f"{k}: {v}" for k, v in test_results.items()])
            multi_query_context = self._fetch_context(test_results)
            hypothetical_doc = self._generate_hypothetical_doc(query_text)
            hyde_docs = self.retriever.invoke(hypothetical_doc)
            hyde_context = "\n".join([
                self._clean_doc_content(doc.page_content)
                for doc in hyde_docs
                if hasattr(doc, 'page_content') and doc.page_content
            ])
            full_context = f"醫學參考資料:\n{multi_query_context}\n\n假設性分析:\n{hyde_context}"

            prompt = f"""
你是一位保險核保分析師，請根據以下資料生成 JSON 格式分析報告：
=== 健康檢查數據 ===
{query_text}

=== 相關醫學知識 ===
{full_context}

=== 輸出格式（嚴格遵守）===
{{
  "overall_risk_score": 75,
  "risk_level_label": "中風險",
  "summary": "存在糖尿病與腎病風險，建議人工審核。",
  "core_recommendation": ["人工審核", "要求近3個月報告"],
  "disease_risk_evaluation": [
    {{"disease": "糖尿病", "risk_score": 65, "risk_level": "中風險", "main_factors": ["HbA1c"], "advice": "加收保費"}},
    {{"disease": "腎病", "risk_score": 55, "risk_level": "中風險", "main_factors": ["CRE"], "advice": "排除腎臟相關理賠"}}
  ],
  "success": true
}}
**僅輸出 JSON！main_factors 必須來自輸入數據！**
"""
            result = self._call_ollama_json(prompt, InsurerAnalysis)
            if not result:
                raise ValueError("LLM 回傳無效格式")

            disease_risk_proto = [
                data_pb2.DiseaseRiskForInsurer(
                    disease=dr.disease,
                    risk_score=dr.risk_score,
                    risk_level=dr.risk_level,
                    main_factors=[f for f in dr.main_factors if f in available_keys],
                    advice=dr.advice
                ) for dr in result.disease_risk_evaluation if dr.disease
            ]

            return data_pb2.InsurerHealthAnalysisResponse(
                overall_risk_score=result.overall_risk_score,
                risk_level_label=result.risk_level_label,
                summary=result.summary,
                core_recommendation=result.core_recommendation,
                disease_risk_evaluation=disease_risk_proto,
                success=result.success
            )

        except Exception as e:
            logger.error(f"保險分析失敗: {str(e)}")
            return data_pb2.InsurerHealthAnalysisResponse(
                overall_risk_score=0, risk_level_label="", summary="分析失敗",
                core_recommendation=[], disease_risk_evaluation=[], success=False
            )

def serve():
    try:
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
        data_pb2_grpc.add_HealthServiceServicer_to_server(HealthAnalysisServicer(), server)
        server.add_insecure_port('[::]:50051')
        logger.info("gRPC 服務器啟動，監聽端口 50051")
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
