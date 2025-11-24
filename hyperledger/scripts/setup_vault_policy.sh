#!/usr/bin/env bash
set -euo pipefail

# 設定 Vault Policy 並確保 token 有正確權限
# - 讀取 policy.hcl 並寫入 Vault
# - 將 policy 附加到 root token（或指定的 token）

ROOT_DIR="/home/Medledger-AI/hyperledger"
POLICY_FILE="$ROOT_DIR/vault/config/policy.hcl"
POLICY_NAME="medledger-policy"
INIT_DIR="$ROOT_DIR/vault/init"
VAULT_HTTP="http://127.0.0.1:8200"
VAULT_ADDR="${VAULT_ADDR:-$VAULT_HTTP}"

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

echo "[INFO] 開始設定 Vault Policy..."

# 1. 檢查 policy 檔案是否存在
if [ ! -f "$POLICY_FILE" ]; then
    echo "[ERROR] Policy 檔案不存在：$POLICY_FILE" >&2
    exit 1
fi

# 2. 寫入 policy
echo "[INFO] 寫入 policy: $POLICY_NAME"
$VAULT_CMD policy write "$POLICY_NAME" "$POLICY_FILE" || {
    echo "[ERROR] 寫入 policy 失敗" >&2
    exit 1
}

# 3. 將 policy 附加到 root token（root token 本身不需要 policy，但可以設定）
# 實際上，root token 擁有所有權限，不需要 policy
# 但我們可以建立一個 token 並附加 policy，供 Go server 使用
echo "[INFO] Policy '$POLICY_NAME' 已寫入 Vault"

# 4. 檢查並確保 KV mount 存在
echo "[INFO] 檢查 KV mount..."
MOUNTS=$($VAULT_CMD auth list -format=json 2>/dev/null || echo "{}")
KV_MOUNT=""
if $VAULT_CMD secrets list -format=json 2>/dev/null | grep -q "kv"; then
    KV_MOUNT="kv"
elif $VAULT_CMD secrets list -format=json 2>/dev/null | grep -q "kv-v2"; then
    KV_MOUNT="kv-v2"
else
    echo "[INFO] 未找到 KV mount，嘗試建立 kv-v2 mount..."
    $VAULT_CMD secrets enable -version=2 -path=kv-v2 kv 2>/dev/null || {
        echo "[WARN] 建立 kv-v2 mount 失敗，可能已存在"
    }
    KV_MOUNT="kv-v2"
fi

echo "[INFO] 使用 KV mount: $KV_MOUNT"

# 5. 確保 Transit mount 存在
echo "[INFO] 檢查 Transit mount..."
if ! $VAULT_CMD secrets list -format=json 2>/dev/null | grep -q "transit"; then
    echo "[INFO] 建立 Transit mount..."
    $VAULT_CMD secrets enable transit 2>/dev/null || {
        echo "[WARN] 建立 Transit mount 失敗，可能已存在"
    }
fi

echo ""
echo "[DONE] Vault Policy 設定完成！"
echo ""
echo "📋 摘要："
echo "   - Policy 名稱: $POLICY_NAME"
echo "   - KV Mount: $KV_MOUNT"
echo "   - Root Token: 已使用（擁有所有權限）"
echo ""
echo "💡 提示："
echo "   - 如果 Go server 使用 root token，則不需要額外設定"
echo "   - 如果使用其他 token，請確保該 token 附加了 '$POLICY_NAME' policy"
echo "   - 環境變數 VAULT_MOUNT 應設為: $KV_MOUNT"

