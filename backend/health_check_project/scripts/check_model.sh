#!/bin/bash
echo "檢查 Ollama 模型狀態..."
echo ""

# 檢查已安裝的模型
echo "已安裝的模型："
ollama list

echo ""
echo "測試 meditron:7b 模型..."
ollama run meditron:7b "請用繁體中文回答：健康檢查中的血糖正常值是多少？" 2>&1 | head -20

echo ""
echo "如果模型輸出異常，建議："
echo "1. 重新下載模型: ollama pull meditron:7b"
echo "2. 或嘗試其他模型: ollama pull llama3:8b 或 ollama pull qwen2.5:7b"
