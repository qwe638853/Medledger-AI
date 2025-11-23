#!/bin/bash

echo "🧹 開始整理專案目錄..."

# 1. 建立新資料夾
mkdir -p data
mkdir -p scripts
mkdir -p docs

# 2. 移動數據文件 (如果檔案存在才移動)
echo "📦 移動數據文件..."
[ -f "health_rules.txt" ] && mv health_rules.txt data/
[ -f "my_test_report.docx" ] && mv my_test_report.docx data/

# 3. 移動文檔
echo "📄 移動文檔..."
[ -f "ARCHITECTURE.md" ] && mv ARCHITECTURE.md docs/

# 4. 移動工具腳本 (Python & Shell)
echo "🛠️ 移動工具腳本..."
[ -f "add_data.py" ] && mv add_data.py scripts/
[ -f "add_knowledge.py" ] && mv add_knowledge.py scripts/
[ -f "create_collection.py" ] && mv create_collection.py scripts/
[ -f "create_docx.py" ] && mv create_docx.py scripts/
[ -f "debug.py" ] && mv debug.py scripts/
[ -f "check_model.sh" ] && mv check_model.sh scripts/
[ -f "fix_import.sh" ] && mv fix_import.sh scripts/
[ -f "start_grpc_server.sh" ] && mv start_grpc_server.sh scripts/
[ -f "start_ollama.sh" ] && mv start_ollama.sh scripts/
[ -f "wsl_setup.sh" ] && mv wsl_setup.sh scripts/
[ -f "setup_new_pc.sh" ] && mv setup_new_pc.sh scripts/
[ -f "clean_project.sh" ] && cp clean_project.sh scripts/ # 備份自己

# 5. 改名主程式 (test.py -> server.py)
if [ -f "test.py" ]; then
    echo "🔄 重新命名主程式: test.py -> server.py"
    mv test.py server.py
fi

# 6. 清理垃圾檔案/資料夾
echo "🗑️ 清理垃圾檔案..."
rm -rf "D:"         
rm -rf "first"      
rm -rf ".venv-1"    
rm -rf "__pycache__"
rm -rf "generated" 

echo "✨ 整理完成！"