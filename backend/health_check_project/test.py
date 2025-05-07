import grpc
from concurrent import futures
import data_pb2
import data_pb2_grpc

# 這裡我們自己實作 health.proto 定義的 HealthService
class HealthServiceServicer(data_pb2_grpc.HealthServiceServicer):
    # 收到上傳報告的請求
    def UploadReport(self, request, context):
        print(request) 
        print(f"📄 Upload Report - ID: {request.report_id}")
        print(f"🔒 Patient Hash: {request.patient_hash}")
        print(f"🧪 Test Results: {request.test_results_json}")
        #接test_results_json，把它做分析
        
        return data_pb2.UploadReportResponse(message="Upload successful!")

    def ClaimReport(self, request, context):
        print(f"✅ Claiming Report ID: {request.report_id}")
        return data_pb2.ClaimReportResponse(message="Claimed successfully!")
    
    def ReadReport(self, request, context):
        fake_report = '{"Glu-AC": "95 mg/dL", "HbA1c": "5.3%", "LDL-C": "125 mg/dL"}'
        print(f"📖 Reading Report ID: {request.report_id}")
        return data_pb2.ReadReportResponse(report_content=fake_report)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    data_pb2_grpc.add_HealthServiceServicer_to_server(HealthServiceServicer(), server)
    server.add_insecure_port('[::]:50051')  # 開在本機50051端口
    server.start()
    print("🚀 gRPC Server Started at port 50051...")
    server.wait_for_termination()

if __name__ == "__main__":
    serve()
