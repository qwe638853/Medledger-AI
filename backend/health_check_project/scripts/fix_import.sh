#!/bin/bash
echo "修復 Python import 問題..."

# 1. 確保 generated 資料夾存在
mkdir -p generated

# 2. 編譯 proto
python -m grpc_tools.protoc \
    -I ./proto \
    -I ./proto/google/api \
    --python_out=./generated \
    --grpc_python_out=./generated \
    ./proto/hyperledger.proto

# 3. 加入 __init__.py
touch generated/__init__.py

# 4. 強制覆蓋 test.py
cat > test.py << 'EOF'
import grpc
import json
from generated import hyperledger_pb2, hyperledger_pb2_grpc

HYPERLEDGER_GRPC_URL = "localhost:50051"  # 改成你的 IP

def main():
    print("正在連線 Hyperledger gRPC 服務...")
    channel = grpc.insecure_channel(HYPERLEDGER_GRPC_URL)
    stub = hyperledger_pb2_grpc.HealthServiceStub(channel)

    data = {"height": 170, "weight": 70, "bmi": 24.2}
    request = hyperledger_pb2.AnalyzeHealthReportRequest(
        report_id="test_001",
        patient_id="user_123",
        test_results_json=json.dumps(data)
    )

    try:
        response = stub.AnalyzeHealthReportForUser(request, timeout=10)
        print(f"健康分數: {response.health_score}")
        print(f"摘要: {response.summary}")
    except grpc.RpcError as e:
        print(f"gRPC 錯誤: {e.code().name} - {e.details()}")
        if e.code() == grpc.StatusCode.UNAVAILABLE:
            print("提示: Hyperledger 服務未啟動或 IP 錯誤")
    except Exception as e:
        print(f"其他錯誤: {e}")

if __name__ == "__main__":
    main()
EOF

echo "修復完成！執行: python test.py"
