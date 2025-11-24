#!/usr/bin/env bash
set -euo pipefail

# 快速修復 Vault 權限問題
# 使用環境變數中的 VAULT_TOKEN（如果有的話）

ROOT_DIR="/home/Medledger-AI/hyperledger"
POLICY_FILE="$ROOT_DIR/vault/config/policy.hcl"
POLICY_NAME="medledger-policy"
VAULT_ADDR="${VAULT_ADDR:-http://127.0.0.1:8200}"

echo "🔧 快速修復 Vault 權限..."
echo ""

# 檢查 Vault 是否運行
if ! curl -sS "$VAULT_ADDR/v1/sys/health" -o /dev/null 2>/dev/null; then
    echo "❌ Vault 未運行"
    exit 1
fi

# 檢查 token
if [ -z "${VAULT_TOKEN:-}" ]; then
    INIT_DIR="$ROOT_DIR/vault/init"
    if [ -f "$INIT_DIR/root_token.txt" ]; then
        export VAULT_TOKEN=$(cat "$INIT_DIR/root_token.txt")
        echo "✅ 已從檔案讀取 root token"
    else
        echo "❌ 找不到 VAULT_TOKEN 環境變數或 root token 檔案"
        echo "   請設定：export VAULT_TOKEN=<你的_token>"
        exit 1
    fi
else
    echo "✅ 使用環境變數中的 VAULT_TOKEN"
fi

export VAULT_ADDR

# 使用容器內的 vault CLI
VAULT_CMD="docker exec -e VAULT_ADDR=$VAULT_ADDR -e VAULT_TOKEN=$VAULT_TOKEN vault vault"
if command -v vault >/dev/null 2>&1; then
    VAULT_CMD="vault"
fi

# 檢查 KV mount
echo ""
echo "📋 檢查 KV mount..."
MOUNTS=$($VAULT_CMD secrets list -format=json 2>/dev/null || echo "{}")
KV_MOUNT=""

if echo "$MOUNTS" | grep -q '"kv-v2"'; then
    KV_MOUNT="kv-v2"
    echo "   ✅ 找到 kv-v2 mount"
elif echo "$MOUNTS" | grep -q '"kv"'; then
    KV_MOUNT="kv"
    echo "   ✅ 找到 kv mount"
else
    echo "   ⚠️  未找到 KV mount，建立 kv-v2..."
    $VAULT_CMD secrets enable -version=2 -path=kv-v2 kv 2>/dev/null || true
    KV_MOUNT="kv-v2"
    echo "   ✅ 已建立 kv-v2 mount"
fi

# 設定 policy
echo ""
echo "📝 設定 policy..."
if [ ! -f "$POLICY_FILE" ]; then
    echo "❌ Policy 檔案不存在：$POLICY_FILE"
    exit 1
fi

# 複製 policy 到容器內（如果使用 docker exec）
if [ "$VAULT_CMD" != "vault" ]; then
    docker cp "$POLICY_FILE" vault:/tmp/policy.hcl >/dev/null 2>&1 || true
    POLICY_PATH="/tmp/policy.hcl"
else
    POLICY_PATH="$POLICY_FILE"
fi

$VAULT_CMD policy write "$POLICY_NAME" "$POLICY_PATH" 2>/dev/null || {
    echo "❌ 寫入 policy 失敗"
    echo "   請確認 VAULT_TOKEN 有足夠權限"
    exit 1
}
echo "   ✅ Policy '$POLICY_NAME' 已寫入"

# 驗證 policy
echo ""
echo "🔍 驗證 policy..."
$VAULT_CMD policy read "$POLICY_NAME" | grep -q "kv-v2/data/insurers" && {
    echo "   ✅ Policy 包含 kv-v2/data/insurers/* 路徑"
} || {
    echo "   ⚠️  Policy 可能不包含 kv-v2 路徑，但已寫入"
}

echo ""
echo "=========================================="
echo "  ✅ 修復完成！"
echo "=========================================="
echo ""
echo "📋 摘要："
echo "   - KV Mount: $KV_MOUNT"
echo "   - Policy: $POLICY_NAME"
echo ""
echo "💡 下一步："
echo "   1. 確保 Go server 使用正確的 mount："
echo "      export VAULT_MOUNT=$KV_MOUNT"
echo ""
echo "   2. 重新啟動 Go server 並測試"

