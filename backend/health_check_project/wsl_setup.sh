#!/bin/bash
# wsl_setup.sh - 修正版：在 Linux 檔案系統執行

set -e

echo "=========================================="
echo "WSL 環境初始化（Linux 檔案系統）"
echo "=========================================="

# 1. 更新套件
sudo apt update -y

# 2. 安裝必要工具
sudo apt install -y python3 python3-pip python3-venv curl unzip

# 3. 建立虛擬環境（在 Linux 檔案系統）
python3 -m venv venv
source venv/bin/activate

# 4. 升級 pip
pip install --upgrade pip

# 5. 建立 requirements.txt
cat > requirements.txt << EOL
grpcio==1.66.0
grpcio-tools==1.66.0
pydantic==2.8.2
ollama==0.1.45
langchain-chroma==0.1.2
langchain-huggingface==0.0.3
chromadb==0.5.3
sentence-transformers==3.0.1
EOL

# 6. 安裝套件
pip install -r requirements.txt

# 7. 安裝 protoc
if ! command -v protoc &> /dev/null; then
    PROTOC_ZIP="protoc-27.2-linux-x86_64.zip"
    curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v27.2.0/$PROTOC_ZIP
    sudo unzip -o $PROTOC_ZIP -d /usr/local
    rm $PROTOC_ZIP
fi

# 8. 編譯 proto
mkdir -p generated
python -m grpc_tools.protoc \
    -I ./proto \
    --python_out=./generated \
    --grpc_python_out=./generated \
    ./proto/health.proto

# 9. 安裝 Ollama
if ! command -v ollama &> /dev/null; then
    curl -fsSL https://ollama.com/install.sh | sh
fi

# 10. 啟動 Ollama
ollama serve &
sleep 5

# 11. 下載模型
ollama pull llama3:8b

echo "=========================================="
echo "完成！現在執行："
echo "  source venv/bin/activate"
echo "  python test.py"
echo "=========================================="
