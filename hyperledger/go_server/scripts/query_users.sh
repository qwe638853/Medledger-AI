#!/bin/bash
# 查詢 SQLite 資料庫中的用戶資料

DB_PATH="database/user_data.sqlite"

if [ ! -f "$DB_PATH" ]; then
    echo "❌ 資料庫文件不存在: $DB_PATH"
    exit 1
fi

echo "=========================================="
echo "📋 查詢 users 表的所有資料"
echo "=========================================="
echo ""

# 檢查 sqlite3 是否安裝
if ! command -v sqlite3 &> /dev/null; then
    echo "❌ sqlite3 未安裝，請先安裝："
    echo "   Ubuntu/Debian: sudo apt-get install sqlite3"
    echo "   macOS: brew install sqlite3"
    exit 1
fi

echo "【方法 1】查看所有欄位（表格格式）："
echo "----------------------------------------"
sqlite3 -header -column "$DB_PATH" "SELECT username, name, date, email, phone FROM users;"

echo ""
echo "【方法 2】查看所有欄位（含密碼雜湊）："
echo "----------------------------------------"
sqlite3 -header -column "$DB_PATH" "SELECT * FROM users;"

echo ""
echo "【方法 3】統計用戶數量："
echo "----------------------------------------"
sqlite3 "$DB_PATH" "SELECT COUNT(*) as user_count FROM users;"

echo ""
echo "【方法 4】查看表結構："
echo "----------------------------------------"
sqlite3 "$DB_PATH" ".schema users"

echo ""
echo "⚠️  注意：username 欄位存的是 SHA256 雜湊值，不是原始用戶ID"
echo "   若要查詢特定用戶，需要先計算雜湊值"

