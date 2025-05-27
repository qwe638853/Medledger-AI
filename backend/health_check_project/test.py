import grpc
from concurrent import futures
import logging
import json
import re
from langchain.prompts import PromptTemplate
from langchain.schema.runnable import RunnablePassthrough
from langchain_chroma import Chroma
from langchain_huggingface import HuggingFaceEmbeddings
from langchain_ollama import OllamaLLM
import data_pb2
import data_pb2_grpc

# 設置日誌
logging.basicConfig(level=logging.INFO, format='%(levelname)s:%(name)s:%(message)s')
logger = logging.getLogger(__name__)

class HealthAnalysisServicer(data_pb2_grpc.HealthServiceServicer):
    def __init__(self):
        embedding = HuggingFaceEmbeddings(model_name="sentence-transformers/all-MiniLM-L6-v2")
        self.vectorstore = Chroma(
            persist_directory="D:/gg/chroma_db",
            embedding_function=embedding,
            collection_name="health_knowledge"
        )
        self.retriever = self.vectorstore.as_retriever(search_kwargs={"k": 3})
        self.llm = OllamaLLM(model="llama3:8b", base_url="http://localhost:11434")
        self.translations = {
            "Glu-AC": "飯前血糖", "HbA1c": "糖化血紅蛋白", "Glu-PC": "飯後血糖", "Alb": "白蛋白", "TP": "總蛋白",
            "AST(GOT)": "天門冬氨酸轉氨酶", "ALT(GPT)": "丙氨酸轉氨酶", "D-Bil": "直接膽紅素", "ALP": "鹼性磷酸酶",
            "T-Bil": "總膽紅素", "UN": "尿素氮", "CRE": "肌酐", "U.A": "尿酸", "T-CHO": "總膽固醇",
            "LDL-C": "低密度脂蛋白膽固醇", "HDL-C": "高密度脂蛋白膽固醇", "TG": "三酸甘油酯", "Hb": "血紅蛋白",
            "Hct": "紅細胞壓積", "PLT": "血小板", "WBC": "白細胞", "RBC": "紅細胞", "hsCRP": "高敏感C反應蛋白",
            "AFP": "甲胎蛋白", "CEA": "癌胚抗原", "CA-125": "癌症抗原125", "CA19-9": "癌症抗原19-9", "BP": "血壓",
            "MCV": "平均紅細胞體積", "MCH": "平均紅細胞血紅蛋白量", "MCHC": "平均紅細胞血紅蛋白濃度", "PT": "凝血酶原時間",
            "aPTT": "活化部分凝血活酶時間", "ESR": "血沉", "RDW-CV": "紅細胞分佈寬度", "Specific Gravity": "尿液比重",
            "PH": "尿液酸鹼值", "Protein (Dipstick)": "蛋白質 (試紙)", "Glucose (Dipstick)": "葡萄糖 (試紙)", "Bilirubin (Dipstick)": "膽紅素 (試紙)",
            "Urobilinogen (Dipstick)": "尿膽原 (試紙)", "RBC (Urine)": "紅細胞 (尿液)", "WBC (Urine)": "白細胞 (尿液)", "Epithelial Cells": "上皮細胞",
            "Casts": "管型", "Ketone": "酮體", "Crystal": "結晶", "Bacteria": "細菌", "Albumin (Dipstick)": "白蛋白 (試紙)",
            "Creatinine (Dipstick)": "肌酐 (試紙)", "Alb/CRE Ratio": "白蛋白/肌酐比率", "Nitrite": "亞硝酸鹽", "Occult Blood": "潛血",
            "WBC Esterase": "白細胞酯酶"
        }
        self.policy_translations = {
            "LIFE INSURANCE": "壽險",
            "HEALTH INSURANCE": "健康險",
            "Health Insurance": "健康險",
            "Life Insurance": "壽險"
        }
        logger.info("🚀 初始化完成")

    def clean_json(self, text):
        if not text or not isinstance(text, str):
            return None
        # 移除無效字符和額外文字
        text = re.sub(r'[\x00-\x1F\x7F-\x9F]', '', text)
        text = re.sub(r'\s+', ' ', text)
        # 移除 JSON 以外的內容（如 "Note:" 或其他評論）
        text = re.sub(r'}[^}]*$', '}', text)
        json_match = re.search(r'\{(?:[^{}]|(?:\{[^{}]*\}))*\}', text, re.DOTALL)
        if json_match:
            json_str = json_match.group(0)
            try:
                return json.loads(json_str)
            except json.JSONDecodeError:
                return None
        return None

    def extract_json(self, text):
        cleaned = self.clean_json(text)
        return json.dumps(cleaned, ensure_ascii=False) if cleaned else None

    def clean_metrics(self, metrics, original_metrics):
        if not isinstance(metrics, dict):
            return original_metrics
        cleaned_metrics = {}
        for key, value in original_metrics.items():
            translated_key = self.translations.get(key, key)
            # 使用原始數據，確保不被 LLM 修改
            if isinstance(value, (int, float)):
                cleaned_metrics[f'"{translated_key}"'] = value
            elif isinstance(value, str):
                match = re.match(r'^-?\d+(\.\d+)?', value)
                if match:
                    cleaned_metrics[f'"{translated_key}"'] = float(match.group(0))
                else:
                    cleaned_metrics[f'"{translated_key}"'] = f'"{value.strip()}"'
            else:
                cleaned_metrics[f'"{translated_key}"'] = f'"{str(value).strip()}"'
        return cleaned_metrics

    def translate_text(self, text):
        if not isinstance(text, str):
            return text
        # 替換指標名稱
        for eng, chi in self.translations.items():
            text = text.replace(eng, chi)
        # 替換常見英文詞彙
        text = text.replace("slightly high", "略高")
        text = text.replace("blood glucose testing", "血糖檢查")
        return text

    def translate_policies(self, policies):
        if not isinstance(policies, list):
            return policies
        return [self.policy_translations.get(policy, policy) for policy in policies]

    def translate_risks(self, risks):
        if not isinstance(risks, list):
            return risks
        translated_risks = []
        for risk in risks:
            translated_risk = risk.copy()
            disease = risk.get("disease", "未知")
            if not disease:
                translated_risk["disease"] = "潛在代謝疾病"
                translated_risk["impact"] = "中度"
                translated_risk["description"] = "本次體檢結果顯示低密度脂蛋白膽固醇和尿素氮值偏高，可能表明有心血管或腎功能風險。"
            else:
                description = risk.get("description", "無描述")
                for eng, chi in self.translations.items():
                    description = description.replace(eng, chi)
                translated_risk["description"] = description
            translated_risks.append(translated_risk)
        return translated_risks

    def generate_hypothetical_doc(self, query_text):
        hyde_prompt = PromptTemplate.from_template("""
            根據以下健康檢查資料，生成一段假設性健康總結（使用繁體中文）：
            {query}
            例如：這位患者的血糖控制良好，但低密度脂蛋白膽固醇偏高，可能有心血管風險。
            **注意**：
            - 必須完全使用繁體中文，不得包含英文或其他語言詞彙。
            - 指標名稱必須使用以下中文名稱：飯前血糖 (Glu-AC), 糖化血紅蛋白 (HbA1c), 飯後血糖 (Glu-PC), 總膽固醇 (T-CHO), 低密度脂蛋白膽固醇 (LDL-C), 高密度脂蛋白膽固醇 (HDL-C), 三酸甘油酯 (TG), 尿素氮 (UN), 高敏感C反應蛋白 (hsCRP), 血壓 (BP)。
            - 僅基於提供的健康檢查資料進行總結，不得假設或添加未提供的數據（如腦中風或虛構指標）。
            - 總結應簡潔並聚焦於主要指標。
        """)
        chain = hyde_prompt | self.llm
        hypothetical_doc = chain.invoke({"query": query_text})
        logger.info(f"HyDE 假設性文件：{hypothetical_doc}")
        return hypothetical_doc.strip()

    def get_multi_query_context(self, test_results):
        query_categories = {
            "blood_sugar": ["Glu-AC", "HbA1c", "Glu-PC"],
            "lipid": ["LDL-C", "HDL-C", "TG", "T-CHO"],
            "liver": ["ALT(GPT)", "AST(GOT)", "ALP", "T-Bil", "D-Bil"],
            "kidney": ["CRE", "UN", "Alb/CRE Ratio"],
            "general": ["Hb", "Hct", "PLT", "WBC", "RBC", "hsCRP"]
        }
        all_docs = []
        default_docs = {
            "blood_sugar": "血糖正常範圍：飯前血糖 70-100 mg/dL，糖化血紅蛋白 4%-6%。",
            "lipid": "血脂正常值：總膽固醇 < 200 mg/dL，低密度脂蛋白膽固醇 < 120 mg/dL，高密度脂蛋白膽固醇 > 40 mg/dL三酸甘油酯 < 150 mg/dL。",
            "liver": "肝功能正常範圍：丙氨酸轉氨酶 0-50 U/L，天門冬氨酸轉氨酶 0-50 U/L，總膽紅素 0.3-1.2 mg/dL。",
            "kidney": "腎功能正常範圍：肌酐 0.6-1.2 mg/dL，尿素氮 7-20 mg/dL，白蛋白/肌酐比率 < 30 mg/g。",
            "general": "血液常規正常範圍：血紅蛋白 12-16 g/dL，白細胞 4-10 x10^3/uL，血小板 150-450 x10^3/uL，高敏感C反應蛋白 < 1 mg/dL。"
        }
        for category, keys in query_categories.items():
            category_query = "\n".join([f"{k}: {v}" for k, v in test_results.items() if k in keys])
            if category_query:
                logger.info(f"Multi-Query 子查詢 ({category})：{category_query}")
                docs_with_scores = self.vectorstore.similarity_search_with_score(category_query, k=3)
                filtered_docs = [doc for doc, score in docs_with_scores if score < 0.7]
                cleaned_docs = [re.sub(r'問題:.*\n回答:', '', doc.page_content.strip(), flags=re.DOTALL) for doc in filtered_docs if doc.page_content]
                if not cleaned_docs:
                    cleaned_docs = [default_docs[category]]
                logger.info(f"Multi-Query 子查詢 ({category}) 結果：{cleaned_docs}")
                all_docs.extend(cleaned_docs)
        context = "\n".join(all_docs) if all_docs else "無相關參考資料"
        return context

    def AnalyzeHealthReportForUser(self, request, context):
        logger.info(f"👤 用戶健康報告分析：報告 ID {request.report_id}")
        try:
            test_results = json.loads(request.test_results_json)
            query_text = "\n".join([f"{k}: {v}" for k, v in test_results.items()])

            multi_query_context = self.get_multi_query_context(test_results)
            logger.info(f"Multi-Query 檢索結果：{multi_query_context}")

            hypothetical_doc = self.generate_hypothetical_doc(query_text)
            hyde_docs_with_scores = self.vectorstore.similarity_search_with_score(hypothetical_doc, k=3)
            filtered_hyde_docs = [doc for doc, score in hyde_docs_with_scores if score < 0.7]
            hyde_context = "\n".join([re.sub(r'問題:.*\n回答:', '', doc.page_content.strip(), flags=re.DOTALL) for doc in filtered_hyde_docs if doc.page_content]) if filtered_hyde_docs else "無相關參考資料"
            logger.info(f"HyDE 檢索結果：{hyde_context}")

            context_text = f"{multi_query_context}\n{hyde_context}".strip() or "無相關參考資料"
            logger.info(f"合併上下文：{context_text}")

            user_prompt = PromptTemplate.from_template("""
                你是一位醫療助理，以下是用戶的健康檢查資料：
                {query}
                參考上下文（若無則忽略）：
                {context}
                根據這些數據，請提供一個詳細的健康總結（逐項分析每個主要指標，與上下文中的正常範圍比較，說明是否異常及潛在影響），具體的改善建議（針對異常指標提供飲食、運動、醫療監測建議），並推薦至少兩種合適的保單類型。
                所有內容必須使用繁體中文，並以 **完整且嚴格的 JSON 格式** 回應，包含以下字段：
                {{"summary": "...", "advice": "...", "recommended_policies": ["...", "..."]}}
                **醫療背景**（僅供參考，優先使用上下文）：
                - 低密度脂蛋白膽固醇（LDL-C）正常範圍 < 120 mg/dL，偏高可能增加心血管疾病風險。
                - 高密度脂蛋白膽固醇（HDL-C）正常範圍 > 40 mg/dL，偏低可能影響心血管健康。
                - 飯前血糖（Glu-AC）正常範圍 70-100 mg/dL，偏高可能提示糖尿病風險。
                - 糖化血紅蛋白（HbA1c）正常範圍 4%-6%，偏高可能表示長期血糖控制問題。
                - 血壓正常範圍：收縮壓 < 120 mmHg，舒張壓 < 80 mmHg。
                - 尿素氮（UN）正常範圍 7-20 mg/dL，偏高可能提示腎功能問題。
                - 高敏感C反應蛋白（hsCRP）正常範圍 < 1 mg/dL，偏高可能提示炎症或心血管風險。
                **注意**：
                - **僅輸出 JSON 內容**，不得包含任何前綴、後綴、說明文字或 Markdown 格式。
                - 確保 JSON 格式完整，包含所有括號和逗號。
                - **summary 和 advice 必須完全使用繁體中文**，不得包含英文或其他語言詞彙。
                - 指標名稱必須使用以下中文名稱：飯前血糖 (Glu-AC), 糖化血紅蛋白 (HbA1c), 飯後血糖 (Glu-PC), 總膽固醇 (T-CHO), 低密度脂蛋白膽固醇 (LDL-C), 高密度脂蛋白膽固醇 (HDL-C), 三酸甘油酯 (TG), 尿素氮 (UN), 高敏感C反應蛋白 (hsCRP), 血壓 (BP)。
                - **所有引號必須使用英文雙引號（"）**，不得使用中文引號（「」）。
                - **recommended_policies 必須是一個字符串數組**，僅包含保單名稱，不得包含描述性文字，至少包含兩種保單類型，且與健康檢查結果相關。
                - **僅使用提供的健康檢查資料進行分析**，不得假設或添加未提供的數據或診斷（如腦中風）。
                - 飲食建議需具體，例如每日鹽分攝入量應少於 5 克。
                - 運動建議需具體，例如每週至少 150 分鐘中等強度運動。
                - 醫療監測建議需具體，例如每三個月檢查一次血脂。
                - 不得包含任何額外文字，否則回應將被視為無效。
            """)

            chain = {"query": RunnablePassthrough(), "context": lambda _: context_text} | user_prompt | self.llm
            result = chain.invoke({"query": query_text})

            logger.info(f"Raw LLM result: {result}")
            json_str = self.extract_json(result)
            if not json_str:
                logger.warning("無法提取有效 JSON，嘗試修復")
                default_response = {
                    "summary": "您的健康數據分析中，部分指標需要進一步確認。",
                    "advice": "建議定期進行健康檢查，並諮詢專業醫師以獲得更詳細的建議。",
                    "recommended_policies": ["健康險", "壽險"]
                }
                json_str = json.dumps(default_response, ensure_ascii=False)

            result_json = json.loads(json_str)
            result_json["summary"] = self.translate_text(result_json.get("summary", "無法分析"))
            result_json["advice"] = self.translate_text(result_json.get("advice", "無建議"))
            recommended_policies = self.translate_policies(result_json.get("recommended_policies", []))
            result_json["recommended_policies"] = recommended_policies if recommended_policies else ["健康險", "壽險"]

            return data_pb2.UserHealthAnalysisResponse(
                summary=result_json.get("summary", "無法分析"),
                advice=result_json.get("advice", "無建議"),
                recommended_policy=", ".join(recommended_policies) if recommended_policies else "無推薦保單",
                success=True
            )

        except Exception as e:
            logger.error(f"用戶健康報告分析失敗：{e}")
            return data_pb2.UserHealthAnalysisResponse(
                summary="分析失敗",
                advice="請稍後重試",
                recommended_policy="無推薦保單",
                success=False
            )

    def AnalyzeHealthReportForInsurer(self, request, context):
        logger.info(f"🏢 保險公司健康報告分析：報告 ID {request.report_id}")
        try:
            test_results = json.loads(request.test_results_json)
            query_text = "\n".join([f"{k}: {v}" for k, v in test_results.items()])

            multi_query_context = self.get_multi_query_context(test_results)
            logger.info(f"Multi-Query 檢索結果：{multi_query_context}")

            hypothetical_doc = self.generate_hypothetical_doc(query_text)
            hyde_docs_with_scores = self.vectorstore.similarity_search_with_score(hypothetical_doc, k=3)
            filtered_hyde_docs = [doc for doc, score in hyde_docs_with_scores if score < 0.7]
            hyde_context = "\n".join([re.sub(r'問題:.*\n回答:', '', doc.page_content.strip(), flags=re.DOTALL) for doc in filtered_hyde_docs if doc.page_content]) if filtered_hyde_docs else "無相關參考資料"
            logger.info(f"HyDE 檢索結果：{hyde_context}")

            context_text = f"{multi_query_context}\n{hyde_context}".strip() or "無相關參考資料"
            logger.info(f"合併上下文：{context_text}")

            insurer_prompt = PromptTemplate.from_template("""
                作為保險公司分析師，你收到以下體檢資料：
                {query}
                參考上下文（若無則忽略）：
                {context}
                請分析所有指標（逐項與正常範圍比較，說明是否異常及潛在影響），評估潛在風險疾病（包括詳細描述和長期影響），建議至少兩種對應保單種類，並輸出 **完整且嚴格的 JSON 格式**：
                {{"summary": "...", "metrics": {{...}}, "policy_types": ["...", "..."], "risks": [{{"disease": "...", "impact": "...", "description": "..."}}, ...], "insurance_suitability": "..."}}
                **醫療背景**（僅供參考，優先使用上下文）：
                - 低密度脂蛋白膽固醇（LDL-C）正常範圍 < 120 mg/dL，偏高可能增加心血管疾病風險。
                - 高密度脂蛋白膽固醇（HDL-C）正常範圍 > 40 mg/dL，偏低可能影響心血管健康。
                - 飯前血糖（Glu-AC）正常範圍 70-100 mg/dL，偏高可能提示糖尿病風險。
                - 糖化血紅蛋白（HbA1c）正常範圍 4%-6%，偏高可能表示長期血糖控制問題。
                - 血壓正常範圍：收縮壓 < 120 mmHg，舒張壓 < 80 mmHg。
                - 尿素氮（UN）正常範圍 7-20 mg/dL，偏高可能提示腎功能問題。
                - 高敏感C反應蛋白（hsCRP）正常範圍 < 1 mg/dL，偏高可能提示炎症或心血管風險。
                **注意**：
                - **僅輸出 JSON 內容**，不得包含任何前綴、後綴、說明文字或 Markdown 格式。
                - 確保 JSON 格式完整，包含所有括號和逗號。
                - **metrics 的鍵必須以雙引號包裹**，並使用提供的指標名稱和數值，嚴格基於輸入數據，不得硬編碼或改變數值。
                - **metrics 的值必須是數字或簡單字符串**，例如 89 或 "127"。
                - **metrics 必須包含所有提供的指標**，不得省略。
                - **summary、policy_types、risks 中的 disease、impact、description 以及 insurance_suitability 必須完全使用繁體中文**，不得包含英文。
                - 指標名稱必須使用以下中文名稱：飯前血糖 (Glu-AC), 糖化血紅蛋白 (HbA1c), 飯後血糖 (Glu-PC), 總膽固醇 (T-CHO), 低密度脂蛋白膽固醇 (LDL-C), 高密度脂蛋白膽固醇 (HDL-C), 三酸甘油酯 (TG), 尿素氮 (UN), 高敏感C反應蛋白 (hsCRP), 血壓 (BP)。
                - **所有引號必須使用英文雙引號（"）**，不得使用中文引號（「」）。
                - **policy_types 必須是一個字符串數組**，僅包含保單名稱，不得包含描述性文字。
                - **risks 不得為空**，必須提供至少一個具體的疾病風險，針對異常指標。
                - **summary 和 insurance_suitability 不得為空**，必須提供具體內容。
                - **僅使用提供的健康檢查資料進行分析**，不得假設或添加未提供的數據（如腦中風）。
                - 不得包含任何額外文字，否則回應將被視為無效。
            """)

            chain = {"query": RunnablePassthrough(), "context": lambda _: context_text} | insurer_prompt | self.llm
            result = chain.invoke({"query": query_text})
            logger.info(f"Raw LLM result: {result}")

            json_str = self.extract_json(result)
            if not json_str:
                raise ValueError("無法從回應中提取有效 JSON")

            result_json = json.loads(json_str)

            metrics = self.clean_metrics(result_json.get("metrics", test_results), test_results)
            metrics_str = ", ".join([f"{k}: {v}" for k, v in metrics.items()]) if isinstance(metrics, dict) else str(metrics)

            risks = self.translate_risks(result_json.get("risks", []))
            if not risks:
                risks = [{
                    "disease": "潛在代謝疾病",
                    "impact": "中度",
                    "description": "本次體檢結果顯示低密度脂蛋白膽固醇和尿素氮值偏高，可能表明有心血管或腎功能風險。"
                }]
            risks_proto = [
                data_pb2.Risk(
                    disease=r.get("disease", "未知"),
                    impact=r.get("impact", "無"),
                    description=r.get("description", "無描述")
                ) for r in risks
            ]

            summary = self.translate_text(result_json.get("summary", "").strip() or "無摘要")
            policy_types = self.translate_policies(result_json.get("policy_types", []))
            if not policy_types or len(policy_types) < 2:
                policy_types = ["健康險", "壽險"]
            insurance_suitability = self.translate_text(result_json.get("insurance_suitability", "").strip() or "請人工審核")

            return data_pb2.InsurerHealthAnalysisResponse(
                summary=summary,
                metrics=metrics_str,
                policy_type=", ".join(policy_types) if policy_types else "無保單",
                risks=risks_proto,
                insurance_suitability=insurance_suitability,
                success=True
            )

        except Exception as e:
            logger.error(f"保險公司健康報告分析失敗：{e}")
            return data_pb2.InsurerHealthAnalysisResponse(
                summary="分析失敗",
                metrics="異常",
                policy_type="無法推薦保單",
                risks=[data_pb2.Risk(disease="未知", impact="無", description="分析失敗")],
                insurance_suitability="請人工審核",
                success=False
            )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    data_pb2_grpc.add_HealthServiceServicer_to_server(HealthAnalysisServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    logger.info("🚀 服務器已啟動，監聽端口 50051")
    server.wait_for_termination()

if __name__ == "__main__":
    serve()