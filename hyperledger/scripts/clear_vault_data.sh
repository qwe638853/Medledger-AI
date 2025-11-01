#!/usr/bin/env bash
set -euo pipefail

# 清空 Vault（KSM）的所有資料
# - 刪除所有 KV 儲存的用戶/診所/保險業者/平台憑證資料
# - 刪除所有 Transit 簽章金鑰和 wrap 金鑰
# ⚠️ 警告：此操作不可逆，請確認後執行

ROOT_DIR="/home/Medledger-AI/hyperledger"
INIT_DIR="$ROOT_DIR/vault/init"
VAULT_HTTP="http://127.0.0.1:8200"
VAULT_ADDR="${VAULT_ADDR:-$VAULT_HTTP}"

# 確認對話
if [ "${FORCE:-}" != "1" ]; then
    echo "⚠️  警告：此操作將刪除 Vault 中的所有資料！"
    echo "   - 所有用戶/診所/保險業者憑證（KV 儲存）"
    echo "   - 所有 Transit 簽章金鑰和 wrap 金鑰"
    echo ""
    read -p "確認要清空所有 Vault 資料嗎？(輸入 'yes' 確認): " confirm
    if [ "$confirm" != "yes" ]; then
        echo "[INFO] 操作已取消"
        exit 0
    fi
fi

# 檢查 Vault 是否運行
echo "[INFO] 檢查 Vault 狀態..."
if ! curl -sS "$VAULT_HTTP/v1/sys/health" -o /dev/null; then
    echo "[ERROR] Vault 未運行或無法連接。請先啟動 Vault：" >&2
    echo "   docker compose -f $ROOT_DIR/docker-compose.yaml up -d vault" >&2
    exit 1
fi

# 讀取 root token
if [ -z "${VAULT_TOKEN:-}" ]; then
    if [ ! -f "$INIT_DIR/root_token.txt" ]; then
        echo "[ERROR] 找不到 root token 檔案：$INIT_DIR/root_token.txt" >&2
        echo "   請手動設定 VAULT_TOKEN 環境變數，或確保 Vault 已初始化" >&2
        exit 1
    fi
    export VAULT_TOKEN=$(cat "$INIT_DIR/root_token.txt")
    echo "[INFO] 已從 $INIT_DIR/root_token.txt 讀取 root token"
fi

export VAULT_ADDR

# 使用容器內的 vault CLI（或本地 vault CLI）
VAULT_CMD="docker exec vault vault"
if command -v vault >/dev/null 2>&1; then
    VAULT_CMD="vault"
fi

echo "[INFO] 開始清空 Vault 資料..."

# 1. 刪除所有 KV 資料
echo "[INFO] 刪除 KV 儲存資料..."
for prefix in "users" "clinics" "insurers" "platform"; do
    echo "   清理 $prefix/ ..."
    if [ "$prefix" = "platform" ]; then
        # platform 是單一路徑
        $VAULT_CMD kv delete -mount=kv "platform" 2>/dev/null || true
    else
        # 列出所有鍵並逐一刪除
        keys=$($VAULT_CMD kv list -mount=kv "$prefix" 2>/dev/null || echo "")
        if [ -n "$keys" ]; then
            echo "$keys" | while IFS= read -r key; do
                if [ -n "$key" ]; then
                    $VAULT_CMD kv delete -mount=kv "$prefix/$key" 2>/dev/null || true
                    echo "     已刪除: $prefix/$key"
                fi
            done
        fi
    fi
done

# 2. 刪除所有 Transit 金鑰
echo "[INFO] 刪除 Transit 金鑰..."
# 列出所有 transit 金鑰
transit_keys=$($VAULT_CMD list -mount=transit keys 2>/dev/null | tail -n +2 || echo "")

if [ -n "$transit_keys" ]; then
    echo "$transit_keys" | while IFS= read -r key; do
        # 移除前導空格和特殊字符
        key=$(echo "$key" | tr -d '[:space:]' | tr -d '/')
        if [ -n "$key" ]; then
            echo "   刪除 Transit 金鑰: $key"
            $VAULT_CMD write -mount=transit keys/$key/rotate 2>/dev/null || true
            $VAULT_CMD delete -mount=transit keys/$key 2>/dev/null || true
        fi
    done
else
    echo "   未發現 Transit 金鑰"
fi

echo ""
echo "[DONE] Vault 資料清空完成！"
echo ""
echo "⚠️  注意："
echo "   - SQLite wallet 中的引用可能已失效（若存在）"
echo "   - 建議同時清理 SQLite 資料庫（若需要）"
echo "   - 鏈上的加密報告資料不受影響（僅 KSM 金鑰被清除）"


