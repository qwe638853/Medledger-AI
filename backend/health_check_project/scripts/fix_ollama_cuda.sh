#!/bin/bash

# 修復 Ollama CUDA 內存錯誤的腳本
# 通過強制使用 CPU 或限制 GPU 使用來解決內存分配問題

echo "🔧 修復 Ollama CUDA 內存錯誤..."

# 停止現有的 Ollama 服務
echo "1️⃣ 停止現有的 Ollama 服務..."
if pgrep -f ollama > /dev/null; then
    pkill -f ollama
    sleep 2
    echo "   ✅ Ollama 服務已停止"
else
    echo "   ℹ️  Ollama 服務未運行"
fi

# 設置環境變數（強制使用 CPU）
export OLLAMA_NUM_GPU=0
export OLLAMA_GPU_LAYERS=0

echo ""
echo "2️⃣ 配置 GPU 設置..."
echo "   OLLAMA_NUM_GPU=${OLLAMA_NUM_GPU} (0 = 使用 CPU)"
echo "   OLLAMA_GPU_LAYERS=${OLLAMA_GPU_LAYERS} (0 = 不使用 GPU 層)"

# 啟動 Ollama 服務
echo ""
echo "3️⃣ 啟動 Ollama 服務（CPU 模式）..."
OLLAMA_NUM_GPU=0 OLLAMA_GPU_LAYERS=0 nohup ollama serve > /tmp/ollama.log 2>&1 &

# 等待服務啟動
echo "   等待服務啟動..."
for i in {1..15}; do
    if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
        echo "   ✅ Ollama 服務已啟動（CPU 模式）"
        break
    fi
    if [ $i -eq 15 ]; then
        echo "   ⚠️  服務啟動可能失敗，請檢查 /tmp/ollama.log"
        exit 1
    fi
    sleep 1
done

echo ""
echo "4️⃣ 驗證服務狀態..."
if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
    echo "   ✅ Ollama 服務正常運行"
    echo ""
    echo "📋 可用模型："
    ollama list 2>/dev/null || echo "   無法獲取模型列表"
else
    echo "   ❌ Ollama 服務無法訪問"
    echo "   請檢查日誌: tail -f /tmp/ollama.log"
    exit 1
fi

echo ""
echo "✅ 修復完成！Ollama 現在運行在 CPU 模式下"
echo ""
echo "💡 提示："
echo "   - 如果以後想使用 GPU，可以設置環境變數："
echo "     export OLLAMA_NUM_GPU=1"
echo "     export OLLAMA_GPU_LAYERS=20  # 根據您的 GPU 內存調整"
echo "   - 然後重啟 Ollama 服務"
echo ""
echo "📝 查看日誌："
echo "   tail -f /tmp/ollama.log"


