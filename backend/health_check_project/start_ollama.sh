#!/bin/bash

# Ollama 啟動腳本
# 用於在 WSL/Linux 環境中啟動 Ollama 服務

echo "🔍 檢查 Ollama 安裝狀態..."

# 檢查 Ollama 是否已安裝
if ! command -v ollama &> /dev/null; then
    echo "❌ Ollama 未安裝"
    echo ""
    echo "📦 安裝方法："
    echo "1. 使用 snap (推薦):"
    echo "   sudo snap install ollama"
    echo ""
    echo "2. 或使用官方安裝腳本:"
    echo "   curl -fsSL https://ollama.com/install.sh | sh"
    echo ""
    echo "3. 或從 GitHub 下載:"
    echo "   https://github.com/ollama/ollama/releases"
    echo ""
    exit 1
fi

echo "✅ Ollama 已安裝"

# 檢查 Ollama 是否正在運行
if pgrep -f ollama > /dev/null; then
    echo "✅ Ollama 服務正在運行"
    echo "   進程 ID: $(pgrep -f ollama)"
else
    echo "⚠️  Ollama 服務未運行，正在啟動..."
    
    # 嘗試啟動 Ollama
    if command -v systemctl &> /dev/null; then
        echo "   使用 systemctl 啟動..."
        sudo systemctl start ollama 2>/dev/null || {
            echo "   systemctl 啟動失敗，嘗試直接啟動..."
            ollama serve &
        }
    else
        echo "   直接啟動 Ollama 服務..."
        nohup ollama serve > /tmp/ollama.log 2>&1 &
        sleep 2
    fi
    
    # 等待服務啟動
    echo "   等待服務啟動..."
    for i in {1..10}; do
        if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
            echo "✅ Ollama 服務已啟動"
            break
        fi
        sleep 1
    done
fi

# 檢查服務是否可訪問
echo ""
echo "🔍 檢查 Ollama 服務連接..."
if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
    echo "✅ Ollama 服務可訪問 (http://localhost:11434)"
    echo ""
    echo "📋 可用模型列表:"
    ollama list 2>/dev/null || echo "   無法獲取模型列表"
else
    echo "❌ Ollama 服務無法訪問"
    echo ""
    echo "💡 可能的原因："
    echo "1. Ollama 服務未正確啟動"
    echo "2. 端口 11434 被占用"
    echo "3. 防火牆阻止連接"
    echo ""
    echo "🔧 手動檢查："
    echo "   netstat -tlnp | grep 11434"
    echo "   curl http://localhost:11434/api/tags"
    exit 1
fi

echo ""
echo "✅ Ollama 準備就緒！"

