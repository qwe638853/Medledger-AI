import grpc
import data_pb2, data_pb2_grpc
import logging
import json

logging.basicConfig(level=logging.INFO, format='%(levelname)s:%(name)s:%(message)s')
logger = logging.getLogger(__name__)

def run():
    try:
        # 連線到 Python Backend 的端口 (50052)
        channel = grpc.insecure_channel('localhost:50052')
        stub = data_pb2_grpc.HealthServiceStub(channel)

        report_id = "report001"
        # 測試數據：極端異常值 (高糖尿病、高血脂、腎衰竭)
        test_results_json = json.dumps({
            "Glu-AC": 200,
            "HbA1c": 8.5,
            "LDL-C": 900,
            "HDL-C": 80,
            "TG": 150,
            "ALT(GPT)": 140,
            "CRE": 6.0,
            "Hb": 7.5,
            "WBC": 16.0
        })

        logger.info(f"\n=== 處理報告 ID: {report_id} ===")

        # --- 1. 用戶健康分析 ---
        user_response = stub.AnalyzeHealthReportForUser(data_pb2.AnalyzeHealthReportRequest(
            report_id=report_id,
            patient_id="patient001",
            test_results_json=test_results_json
        ))
        
        logger.info(f"👤 用戶健康分析 (總結)")
        logger.info(f"健康分數: {user_response.health_score} (0-100 分，0=極度危險)")
        logger.info(f"總結分析: {user_response.summary}")
        logger.info(f"個人詳細分析: {user_response.personal_analysis}")
        
        # 顯示推薦，但隱藏詳細的疾病列表
        logger.info(f"智能保險推薦: {list(user_response.insurance_recommendation)}")
        logger.info(f"健康防護計畫: {list(user_response.protection_plan)}")
        
        if user_response.disease_risks:
             logger.info(f"⚠️ 詳細疾病風險: 已在摘要中總結，共有 {len(user_response.disease_risks)} 項風險。")


        # --- 2. 保險業者分析 ---
        insurer_response = stub.AnalyzeHealthReportForInsurer(data_pb2.AnalyzeHealthReportRequest(
            report_id=report_id,
            patient_id="patient001",
            test_results_json=test_results_json
        ))
        
        logger.info(f"\n🏢 保險業者健康分析 (核保決策)")
        logger.info(f"整體風險分數: {insurer_response.overall_risk_score} (0-100 分，100=極高風險)")
        logger.info(f"風險分級標籤: {insurer_response.risk_level_label}")
        logger.info(f"核保摘要: {insurer_response.summary}")
        logger.info(f"專業核保建議: {list(insurer_response.core_recommendation)}")
        
        # 隱藏分類疾病列表
        if insurer_response.disease_risk_evaluation:
             logger.info(f"⚠️ 詳細核保評估: 已在摘要中總結，共有 {len(insurer_response.disease_risk_evaluation)} 項分類評估。")

        channel.close()

    except grpc.RpcError as e:
        logger.error(f"gRPC 錯誤: {e.details()}")
    except Exception as e:
        logger.error(f"執行失敗: {str(e)}")

if __name__ == "__main__":
    run()