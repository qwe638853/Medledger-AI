#!/bin/bash
# 啟動 Python gRPC 服務器

cd "$(dirname "$0")"

# 激活虛擬環境
if [ ! -d "venv" ]; then
    echo "錯誤: 虛擬環境不存在，請先運行 wsl_setup.sh"
    exit 1
fi

source venv/bin/activate

# 檢查依賴
python -c "import grpc" 2>/dev/null || {
    echo "錯誤: grpc 模塊未安裝，正在安裝..."
    pip install -r requirements.txt
}

# 設置環境變數（如果需要）
export PYTHON_BACKEND_GRPC_PORT=${PYTHON_BACKEND_GRPC_PORT:-50052}
export PYTHON_BACKEND_GRPC_HOST=${PYTHON_BACKEND_GRPC_HOST:-"[::]"}

# 啟動服務器
echo "========================================="
echo "啟動 Python gRPC 服務器"
echo "監聽端口: $PYTHON_BACKEND_GRPC_PORT"
echo "監聽地址: $PYTHON_BACKEND_GRPC_HOST"
echo ""
echo "注意:"
echo "  - Go Server gRPC 運行在 :50051"
echo "  - Python Backend gRPC 運行在 :$PYTHON_BACKEND_GRPC_PORT"
echo "  - Go Server HTTP Gateway 運行在 :8080"
echo "========================================="
echo ""
python test.py
