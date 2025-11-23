#!/bin/bash

echo "🚀 開始安裝 Medledger-AI 開發環境..."

# 1. 安裝系統依賴 (Python, Pip, Venv)
echo "📦 正在安裝系統套件..."
sudo apt update
sudo apt install -y python3 python3-pip python3-venv curl git

# 2. 安裝 Ollama (如果未安裝)
if ! command -v ollama &> /dev/null; then
    echo "🦙 正在安裝 Ollama..."
    curl -fsSL https://ollama.com/install.sh | sh
else
    echo "✅ Ollama 已安裝"
fi

# 3. 啟動 Ollama 服務 (背景執行)
echo "🔄 正在啟動 Ollama 服務..."
sudo systemctl start ollama

# 4. 下載 Llama 3 模型 (我們最終確認使用的模型)
echo "📥 正在下載 Llama 3:8b 模型 (這可能需要一點時間)..."
ollama pull llama3:8b

# 5. 進入專案目錄
PROJECT_DIR="$HOME/medledger-ai/Medledger-AI/backend/health_check_project"
if [ -d "$PROJECT_DIR" ]; then
    cd "$PROJECT_DIR"
    echo "📂 進入目錄: $PROJECT_DIR"
else
    echo "❌ 錯誤：找不到專案目錄 $PROJECT_DIR，請確認您已 Clone 專案"
    exit 1
fi

# 6. 建立並啟用虛擬環境
if [ ! -d ".venv" ]; then
    echo "🐍 正在建立 Python 虛擬環境 (.venv)..."
    python3 -m venv .venv
fi
source .venv/bin/activate

# 7. 安裝 Python 套件
echo "📦 正在安裝 Python 依賴套件..."
# 這裡列出了我們之前用到的所有關鍵套件
pip install --upgrade pip
pip install grpcio grpcio-tools protobuf
pip install ollama pydantic
pip install langchain-chroma langchain-huggingface langchain-community sentence-transformers

echo "🎉 安裝完成！"
echo "請執行 'source .venv/bin/activate' 來啟用環境，然後就可以執行 test.py 了！"
