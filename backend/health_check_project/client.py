import grpc
import data_pb2, data_pb2_grpc
import logging
import json

logging.basicConfig(level=logging.INFO, format='%(levelname)s:%(name)s:%(message)s')
logger = logging.getLogger(__name__)

def run():
    try:
        channel = grpc.insecure_channel('localhost:50051')
        stub = data_pb2_grpc.HealthServiceStub(channel)

        report_id = "report001"
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

        # 用戶分析
        user_response = stub.AnalyzeHealthReportForUser(data_pb2.AnalyzeHealthReportRequest(
            report_id=report_id,
            patient_id="patient001",
            test_results_json=test_results_json
        ))
        logger.info(f"👤 用戶健康分析")
        logger.info(f"健康分數: {user_response.health_score} (0-100 分，{user_response.health_score} 表示當前健康水平，建議根據異常指標進一步評估)")
        logger.info(f"摘要: {user_response.summary}")
        if user_response.disease_risks:
            for dr in user_response.disease_risks:
                logger.info(f"疾病風險: {dr.disease}")
                logger.info(f"  - 風險水平: {dr.risk_level} ({dr.risk_percent}%)")
                logger.info(f"  - 主要因素: {', '.join(dr.main_factors)}")
                logger.info(f"  - 建議: {dr.advice}")
                logger.info(f"  - 詳細解釋: 根據您的數據（如 {dr.main_factors[0] if dr.main_factors else '未知'}），可能需要立即諮詢醫生進行針對性檢查和治療。")
        else:
            logger.info("疾病風險: 無 (無明顯疾病風險，但建議定期檢查以確保健康)")
        logger.info(f"智能保險推薦: {list(user_response.insurance_recommendation)} (基於風險評估，建議考慮這些保險以覆蓋潛在醫療費用)")
        logger.info(f"健康防護計畫: {list(user_response.protection_plan)} (請嚴格遵循這些計畫，並每 3-6 個月進行一次全面健康檢查)")

        # 保險業者分析
        insurer_response = stub.AnalyzeHealthReportForInsurer(data_pb2.AnalyzeHealthReportRequest(
            report_id=report_id,
            patient_id="patient001",
            test_results_json=test_results_json
        ))
        logger.info(f"🏢 保險業者健康分析")
        logger.info(f"整體風險分數: {insurer_response.overall_risk_score} (0-100 分，高分表示高風險，需關注異常數據)")
        logger.info(f"風險分級標籤: {insurer_response.risk_level_label}")
        logger.info(f"核保摘要: {insurer_response.summary}")
        logger.info(f"專業核保建議: {list(insurer_response.core_recommendation)} (建議與醫生和保險專家討論這些建議)")
        if insurer_response.disease_risk_evaluation:
            for dr in insurer_response.disease_risk_evaluation:
                logger.info(f"分類疾病風險: {dr.disease}")
                logger.info(f"  - 風險分數: {dr.risk_score} 分")
                logger.info(f"  - 風險水平: {dr.risk_level}")
                logger.info(f"  - 主要因素: {', '.join(dr.main_factors)}")
                logger.info(f"  - 建議: {dr.advice}")
                logger.info(f"  - 詳細解釋: 您的 {dr.main_factors[0] if dr.main_factors else '未知'} 可能需要長期監控，請提供最新數據以重新評估保險資格。")
        else:
            logger.info("分類疾病風險: 無 (無明顯高風險疾病，但建議定期更新健康數據)")

        channel.close()

    except grpc.RpcError as e:
        logger.error(f"gRPC 錯誤: {e.details()}")
    except Exception as e:
        logger.error(f"執行失敗: {str(e)}")

if __name__ == "__main__":
    run()