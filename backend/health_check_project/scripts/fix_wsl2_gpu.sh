#!/bin/bash

# 修復 WSL2 環境下 Ollama GPU 內存分配問題
# 問題：GPU 內存顯示為 0，但 Ollama 仍嘗試使用 GPU 導致分配失敗

echo "🔍 診斷 WSL2 GPU 環境..."

# 檢查是否在 WSL2 環境
if [ -f /proc/version ] && grep -qi microsoft /proc/version; then
    echo "✅ 檢測到 WSL2 環境"
    IS_WSL2=true
else
    echo "ℹ️  非 WSL2 環境"
    IS_WSL2=false
fi

# 檢查 CUDA 是否可用
echo ""
echo "1️⃣ 檢查 CUDA 可用性..."
if command -v nvidia-smi &> /dev/null; then
    echo "   NVIDIA-SMI 可用："
    nvidia-smi --query-gpu=name,memory.total,memory.used,memory.free --format=csv,noheader 2>/dev/null || echo "   無法查詢 GPU 信息"
else
    echo "   ⚠️  nvidia-smi 不可用"
fi

# 檢查 Ollama 是否正在運行
echo ""
echo "2️⃣ 檢查 Ollama 服務狀態..."
if pgrep -f ollama > /dev/null; then
    echo "   ✅ Ollama 正在運行 (PID: $(pgrep -f ollama))"
    echo "   停止 Ollama 服務..."
    pkill -f ollama
    sleep 3
else
    echo "   ℹ️  Ollama 未運行"
fi

# WSL2 特殊處理：使用分層 GPU 載入
echo ""
echo "3️⃣ 配置 WSL2 GPU 設置..."

if [ "$IS_WSL2" = true ]; then
    echo "   🔧 針對 WSL2 優化配置..."
    echo "   - 使用分層 GPU 載入（避免一次性分配大量內存）"
    echo "   - 設置較少的 GPU 層數以確保穩定"
    
    # WSL2 建議：使用較少的 GPU 層數，避免內存分配問題
    GPU_LAYERS=${OLLAMA_GPU_LAYERS:-15}  # 默認 15 層（對於 8GB GPU 較安全）
    NUM_GPU=${OLLAMA_NUM_GPU:-1}
    
    echo "   配置：OLLAMA_NUM_GPU=${NUM_GPU}, OLLAMA_GPU_LAYERS=${GPU_LAYERS}"
    
    # 設置額外的 WSL2 優化環境變數
    export OLLAMA_NUM_GPU=${NUM_GPU}
    export OLLAMA_GPU_LAYERS=${GPU_LAYERS}
    export OLLAMA_KEEP_ALIVE=5m  # 保持模型在內存中的時間
    
    # 啟動 Ollama
    echo ""
    echo "4️⃣ 啟動 Ollama 服務（WSL2 優化模式）..."
    OLLAMA_NUM_GPU=${NUM_GPU} \
    OLLAMA_GPU_LAYERS=${GPU_LAYERS} \
    OLLAMA_KEEP_ALIVE=5m \
    nohup ollama serve > /tmp/ollama.log 2>&1 &
    
else
    # 非 WSL2 環境：使用標準配置
    GPU_LAYERS=${OLLAMA_GPU_LAYERS:-20}
    NUM_GPU=${OLLAMA_NUM_GPU:-1}
    
    echo "   配置：OLLAMA_NUM_GPU=${NUM_GPU}, OLLAMA_GPU_LAYERS=${GPU_LAYERS}"
    
    export OLLAMA_NUM_GPU=${NUM_GPU}
    export OLLAMA_GPU_LAYERS=${GPU_LAYERS}
    
    echo ""
    echo "4️⃣ 啟動 Ollama 服務..."
    OLLAMA_NUM_GPU=${NUM_GPU} \
    OLLAMA_GPU_LAYERS=${GPU_LAYERS} \
    nohup ollama serve > /tmp/ollama.log 2>&1 &
fi

# 等待服務啟動
echo "   等待服務啟動..."
for i in {1..20}; do
    if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
        echo "   ✅ Ollama 服務已啟動"
        break
    fi
    if [ $i -eq 20 ]; then
        echo "   ⚠️  服務啟動可能失敗，請檢查 /tmp/ollama.log"
        echo "   查看日誌：tail -20 /tmp/ollama.log"
        exit 1
    fi
    sleep 1
done

# 驗證服務
echo ""
echo "5️⃣ 驗證服務狀態..."
if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
    echo "   ✅ Ollama 服務正常運行"
    echo ""
    echo "📋 可用模型："
    ollama list 2>/dev/null || echo "   無法獲取模型列表"
else
    echo "   ❌ Ollama 服務無法訪問"
    exit 1
fi

echo ""
echo "✅ 修復完成！"
echo ""
echo "📊 當前配置："
echo "   - OLLAMA_NUM_GPU=${NUM_GPU}"
echo "   - OLLAMA_GPU_LAYERS=${GPU_LAYERS}"
if [ "$IS_WSL2" = true ]; then
    echo "   - WSL2 優化模式：啟用"
fi
echo ""
echo "💡 如果仍有問題，可以嘗試："
echo "   1. 完全使用 CPU：export OLLAMA_NUM_GPU=0 && export OLLAMA_GPU_LAYERS=0"
echo "   2. 減少 GPU 層數：export OLLAMA_GPU_LAYERS=10"
echo "   3. 查看詳細日誌：tail -f /tmp/ollama.log"
echo ""
echo "📝 測試模型載入："
echo "   ollama run qwen2.5:14b '測試'"

