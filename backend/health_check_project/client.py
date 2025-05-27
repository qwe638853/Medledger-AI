import grpc
import data_pb2
import data_pb2_grpc
import logging
import json

# 設定日誌
logging.basicConfig(level=logging.INFO, format='%(levelname)s:%(name)s:%(message)s')
logger = logging.getLogger(__name__)

def run():
    # 連接到 gRPC 服務器
    try:
        channel = grpc.insecure_channel('localhost:50051')
        stub = data_pb2_grpc.HealthServiceStub(channel)

        # 測試報告 ID
        report_id = "report001"
        test_results_json = '''
        {
            "Glu-AC": "89 mg/dL",
            "HbA1c": "4.1 %",
            "Glu-PC": "124 mg/dL",
            "Alb": "4.5 g/dL",
            "TP": "6.5 g/dL",
            "AST（GOT）": "27 U/L",
            "ALT（GPT）": "10 U/L",
            "D-Bil": "0.03 mg/dL",
            "ALP": "74 U/L",
            "T-Bil": "0.7 mg/dL",
            "UN": "23 mg/dL",
            "CRE": "1.2 mg/dL",
            "U.A": "4.9 mg/dL",
            "T-CHO": "164 mg/dL",
            "LDL-C": "128 mg/dL",
            "HDL-C": "54 mg/dL",
            "TG": "143 mg/dL",
            "Hb": "13.3 g/dL",
            "Hct": "49.7 %",
            "PLT": "286 x10^3/uL",
            "WBC": "4.04 x10^3/uL",
            "RBC": "4.56 x10^6/uL",
            "hsCRP": "0.38 mg/dL",
            "AFP": "14 ng/mL",
            "CEA": "2.8 ng/mL",
            "CA-125": "28 U/mL",
            "CA19-9": "29 U/mL",
            "BP": "127/61 mmHg",
            "MCV": "96.2 fL",
            "MCH": "26.3 pg",
            "MCHC": "33.8 g/dL",
            "PT": "10.4 sec",
            "aPTT": "27.8 sec",
            "ESR": "5 mm/hr",
            "RDW-CV": "12.5 %",
            "Specific Gravity": "1.033",
            "PH": "6.0",
            "Protein (Dipstick)": "-",
            "Glucose (Dipstick)": "-",
            "Bilirubin (Dipstick)": "-",
            "Urobilinogen (Dipstick)": "0.4 mg/dL",
            "RBC (Urine)": "0 /HPF",
            "WBC (Urine)": "5 /HPF",
            "Epithelial Cells": "5 /HPF",
            "Casts": "2 /LPF",
            "Ketone": "-",
            "Crystal": "None",
            "Bacteria": "-",
            "Albumin (Dipstick)": "10 mg/L",
            "Creatinine (Dipstick)": "N/A",
            "Alb/CRE Ratio": "10",
            "Nitrite": "-",
            "Occult Blood": "-",
            "WBC Esterase": "-"
        }
        '''

        logger.info(f"\n=== 處理報告 ID: {report_id} ===")

        # 用戶分析
        user_response = stub.AnalyzeHealthReportForUser(data_pb2.AnalyzeHealthReportRequest(
            report_id=report_id,
            patient_hash="patient001",
            test_results_json=test_results_json
        ))
        logger.info(f"👤 用戶分析")
        logger.info(f"👤 總結: {user_response.summary}")
        logger.info(f"👤 建議: {user_response.advice}")
        policy_map = {
            "standard health insurance plan": "標準健康保單",
            "高風險健康保單": "高風險健康保單"
        }
        policy = policy_map.get(user_response.recommended_policy.lower() if user_response.recommended_policy else "", user_response.recommended_policy)
        logger.info(f"👤 推薦保單: {policy if policy else '無推薦保單'}")

        # 保險分析
        insurer_response = stub.AnalyzeHealthReportForInsurer(data_pb2.AnalyzeHealthReportRequest(
            report_id=report_id,
            patient_hash="patient001",
            test_results_json=test_results_json
        ))
        logger.info(f"🏢 保險分析")
        logger.info(f"🏢 總結: {insurer_response.summary}")
        logger.info(f"🏢 指標: {insurer_response.metrics}")
        logger.info(f"🏢 保單類型: {insurer_response.policy_type}")
        if insurer_response.risks:
            risks_str = [f"{risk.disease} ({risk.impact}): {risk.description}" for risk in insurer_response.risks]
            logger.info(f"🏢 疾病風險: {risks_str}")

        channel.close()

    except grpc.RpcError as e:
        logger.error(f"gRPC 錯誤: {e.details()}")
    except Exception as e:
        logger.error(f"執行失敗: {str(e)}")

if __name__ == "__main__":
    run()