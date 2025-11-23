import grpc
import data_pb2, data_pb2_grpc
import logging
import json
import os
import sys
from google.protobuf.json_format import MessageToDict

logging.basicConfig(level=logging.INFO, format='%(levelname)s:%(name)s:%(message)s')
logger = logging.getLogger(__name__)

def run():
    try:
        # 連線到 Python Backend 的端口 (50052)
        channel = grpc.insecure_channel('localhost:50052')
        stub = data_pb2_grpc.HealthServiceStub(channel)

        report_id = "json_test_001"
        test_results_json = ""

        # 1. 檢查是否有傳入檔案 (支援 PDF/Word 解析)
        if len(sys.argv) > 1:
            file_path = sys.argv[1]
            if os.path.exists(file_path):
                logger.info(f"📂 正在請求伺服器解析文件: {file_path}")
                
                # 讀取檔案二進制內容
                with open(file_path, "rb") as f:
                    file_content = f.read()
                
                file_type = "pdf" if file_path.lower().endswith(".pdf") else "docx"

                # 呼叫 gRPC ParseDocument
                parse_response = stub.ParseDocument(data_pb2.ParseDocumentRequest(
                    file_content=file_content,
                    file_type=file_type
                ))

                if parse_response.success:
                    logger.info("✅ 文件解析成功！")
                    test_results_json = parse_response.result_json
                    # 這裡只印出 log，不印出內容干擾最後的 JSON
                    logger.info(f"解析到的原始數據: {test_results_json}")
                else:
                    logger.error(f"❌ 解析失敗: {parse_response.error_message}")
                    return
            else:
                logger.error(f"找不到檔案: {file_path}")
                return
        else:
            # 2. 如果沒有檔案，使用內建測試數據
            logger.info("⚡ 未指定檔案，使用內建測試數據...")
            test_results_json = json.dumps({
                "Glu-AC": "92 mg/dL",
  "HbA1c": "5.2 %",
  "Glu-PC": "115 mg/dL",
  "Alb": "4.6 g/dL",
  "TP": "7.2 g/dL",
  "AST(GOT)": "22 U/L",
  "ALT(GPT)": "18 U/L",
  "D-Bil": "0.1 mg/dL",
  "ALP": "65 U/L",
  "T-Bil": "0.8 mg/dL",
  "UN": "14 mg/dL",
  "CRE": "0.9 mg/dL",
  "U.A": "5.5 mg/dL",
  "T-CHO": "175 mg/dL",
  "LDL-C": "95 mg/dL",
  "HDL-C": "62 mg/dL",
  "TG": "98 mg/dL",
  "Hb": "14.5 g/dL",
  "Hct": "44.0 %",
  "PLT": "250 x10^3/uL",
  "WBC": "6.50 x10^3/uL",
  "RBC": "4.80 x10^6/uL",
  "hsCRP": "0.12 mg/dL",
  "AFP": "3.5 ng/mL",
  "CEA": "1.2 ng/mL",
  "CA-125": "15 U/mL",
  "CA19-9": "10 U/mL",
  "BP": "118/76 mmHg",
  "MCV": "90.0 fL",
  "MCH": "30.0 pg",
  "MCHC": "33.5 g/dL",
  "PT": "11.0 sec",
  "aPTT": "30.0 sec",
  "ESR": "8 mm/hr",
  "RDW-CV": "13.0 %",
  "Specific Gravity": "1.020",
  "PH": "6.5",
  "Protein (Dipstick)": "-",
  "Glucose (Dipstick)": "-",
  "Bilirubin (Dipstick)": "-",
  "Urobilinogen (Dipstick)": "0.2 mg/dL",
  "RBC (Urine)": "0 /HPF",
  "WBC (Urine)": "1 /HPF",
  "Epithelial Cells": "1 /HPF",
  "Casts": "0 /LPF",
  "Ketone": "-",
  "Crystal": "None",
  "Bacteria": "-",
  "Albumin (Dipstick)": "5 mg/L",
  "Creatinine (Dipstick)": "100 mg/dL",
  "Alb/CRE Ratio": "10",
  "Nitrite": "-",
  "Occult Blood": "-",
  "WBC Esterase": "-"
            })

        logger.info("🚀 開始進行健康分析並生成 JSON...")

        # --- 用戶分析 ---
        user_response = stub.AnalyzeHealthReportForUser(data_pb2.AnalyzeHealthReportRequest(
            report_id=report_id, patient_id="P_JSON", test_results_json=test_results_json))
        
        # [修正] 移除有問題的 always_print_primitive_fields 參數
        # 使用最基本的轉換，保證相容性
        user_data = MessageToDict(
            user_response, 
            preserving_proto_field_name=True
        )

        # --- 保險分析 ---
        ins_response = stub.AnalyzeHealthReportForInsurer(data_pb2.AnalyzeHealthReportRequest(
            report_id=report_id, patient_id="P_JSON", test_results_json=test_results_json))
        
        # [修正] 移除有問題的參數
        insurer_data = MessageToDict(
            ins_response, 
            preserving_proto_field_name=True
        )

        # --- 整合並輸出最終 JSON ---
        final_output = {
            "report_id": report_id,
            "timestamp": "2025-11-23T12:00:00Z", # 模擬時間戳
            "data": {
                "user_analysis": user_data,
                "insurer_analysis": insurer_data
            }
        }

        print("\n" + "="*20 + " JSON OUTPUT START " + "="*20)
        # ensure_ascii=False 讓中文能正常顯示
        print(json.dumps(final_output, indent=2, ensure_ascii=False))
        print("="*20 + " JSON OUTPUT END " + "="*20 + "\n")

        channel.close()

    except grpc.RpcError as e:
        logger.error(f"gRPC 錯誤: {e.details()}")
    except Exception as e:
        logger.error(f"執行失敗: {str(e)}")

if __name__ == "__main__":
    run()