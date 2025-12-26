#!/bin/bash

# 下載 Ollama 模型腳本

echo "📥 下載 Ollama 模型..."

# 檢查 Ollama 服務是否運行
if ! curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
    echo "❌ Ollama 服務未運行，請先啟動 Ollama"
    echo "   執行: ollama serve"
    exit 1
fi

echo "✅ Ollama 服務正在運行"
echo ""

# 檢查當前已安裝的模型
echo "📋 當前已安裝的模型："
ollama list
echo ""

# 模型選項
MODEL_OPTIONS=(
    "qwen2.5:7b - 較小，適合 8GB GPU（推薦）"
    "qwen2.5:14b - 較大，需要更多內存"
    "llama3:8b - 穩定，通用模型"
    "llama3.2:3b - 最小，快速"
    "自定義模型名稱"
)

echo "請選擇要下載的模型："
for i in "${!MODEL_OPTIONS[@]}"; do
    echo "  $((i+1)). ${MODEL_OPTIONS[$i]}"
done

read -p "請輸入選項 (1-5，默認 1): " choice
choice=${choice:-1}

case $choice in
    1)
        MODEL="qwen2.5:7b"
        ;;
    2)
        MODEL="qwen2.5:14b"
        ;;
    3)
        MODEL="llama3:8b"
        ;;
    4)
        MODEL="llama3.2:3b"
        ;;
    5)
        read -p "請輸入模型名稱（例如: qwen2.5:7b）: " MODEL
        ;;
    *)
        MODEL="qwen2.5:7b"
        echo "使用默認模型: $MODEL"
        ;;
esac

echo ""
echo "📥 正在下載模型: $MODEL"
echo "   （這可能需要幾分鐘，請耐心等待...）"
echo ""

# 下載模型
if ollama pull "$MODEL"; then
    echo ""
    echo "✅ 模型下載成功！"
    echo ""
    echo "📋 已安裝的模型："
    ollama list
    echo ""
    echo "💡 提示："
    echo "   如果您的 server.py 配置的模型名稱不同，請："
    echo "   1. 設置環境變數: export OLLAMA_MODEL=$MODEL"
    echo "   2. 或修改 server.py 中的 MODEL_NAME"
    echo ""
    echo "   當前 server.py 默認模型: qwen2.5:14b"
    echo "   下載的模型: $MODEL"
    echo ""
    
    if [ "$MODEL" != "qwen2.5:14b" ]; then
        echo "⚠️  模型名稱不匹配！"
        echo "   建議執行以下命令設置環境變數："
        echo "   export OLLAMA_MODEL=$MODEL"
        echo ""
        echo "   或修改 server.py 第 26 行："
        echo "   MODEL_NAME = os.getenv(\"OLLAMA_MODEL\", \"$MODEL\")"
    fi
else
    echo ""
    echo "❌ 模型下載失敗"
    echo "   請檢查："
    echo "   1. 網絡連接"
    echo "   2. 模型名稱是否正確"
    echo "   3. Ollama 服務是否正常運行"
    exit 1
fi

