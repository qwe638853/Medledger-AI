#!/usr/bin/env bash
set -euo pipefail

# 修復 Vault 權限問題
# - 檢查並初始化 Vault（如需要）
# - 設定 policy
# - 確保 KV mount 存在且名稱正確
# - 提供修復建議

ROOT_DIR="/home/Medledger-AI/hyperledger"
INIT_DIR="$ROOT_DIR/vault/init"
POLICY_FILE="$ROOT_DIR/vault/config/policy.hcl"
POLICY_NAME="medledger-policy"
VAULT_HTTP="http://127.0.0.1:8200"
VAULT_ADDR="${VAULT_ADDR:-$VAULT_HTTP}"

echo "=========================================="
echo "  Vault 權限修復腳本"
echo "=========================================="
echo ""

# 1. 檢查 Vault 是否運行
echo "[1/5] 檢查 Vault 狀態..."
if ! curl -sS "$VAULT_HTTP/v1/sys/health" -o /dev/null 2>/dev/null; then
    echo "❌ Vault 未運行或無法連接"
    echo "   請先啟動 Vault："
    echo "   cd $ROOT_DIR && bash scripts/start_vault.sh"
    exit 1
fi

HEALTH_CODE=$(curl -s -o /dev/null -w '%{http_code}' "$VAULT_HTTP/v1/sys/health" 2>/dev/null || echo "000")
echo "   Vault 健康狀態碼: $HEALTH_CODE (200=正常, 503=封印, 501=未初始化)"

# 2. 初始化 Vault（如需要）
if [ "$HEALTH_CODE" = "501" ]; then
    echo ""
    echo "[2/5] Vault 未初始化，正在初始化..."
    bash "$ROOT_DIR/scripts/start_vault.sh"
    HEALTH_CODE=$(curl -s -o /dev/null -w '%{http_code}' "$VAULT_HTTP/v1/sys/health" 2>/dev/null || echo "000")
fi

# 3. 解封 Vault（如需要）
if [ "$HEALTH_CODE" = "503" ]; then
    echo ""
    echo "[2/5] Vault 已封印，正在解封..."
    if [ -f "$INIT_DIR/unseal_key.txt" ]; then
        UNSEAL_KEY=$(cat "$INIT_DIR/unseal_key.txt")
        docker exec vault vault operator unseal "$UNSEAL_KEY" >/dev/null 2>&1 || true
        HEALTH_CODE=$(curl -s -o /dev/null -w '%{http_code}' "$VAULT_HTTP/v1/sys/health" 2>/dev/null || echo "000")
        if [ "$HEALTH_CODE" = "200" ]; then
            echo "   ✅ Vault 已解封"
        else
            echo "   ⚠️  解封失敗，請手動解封"
        fi
    else
        echo "   ⚠️  找不到 unseal key，請手動解封"
    fi
fi

if [ "$HEALTH_CODE" != "200" ]; then
    echo "❌ Vault 無法正常運行（狀態碼: $HEALTH_CODE）"
    exit 1
fi

# 4. 讀取 root token
echo ""
echo "[3/5] 讀取 Vault token..."
if [ -z "${VAULT_TOKEN:-}" ]; then
    if [ ! -f "$INIT_DIR/root_token.txt" ]; then
        echo "❌ 找不到 root token 檔案：$INIT_DIR/root_token.txt"
        echo "   請先初始化 Vault："
        echo "   cd $ROOT_DIR && bash scripts/start_vault.sh"
        exit 1
    fi
    export VAULT_TOKEN=$(cat "$INIT_DIR/root_token.txt")
    echo "   ✅ 已從檔案讀取 root token"
else
    echo "   ✅ 使用環境變數中的 VAULT_TOKEN"
fi

export VAULT_ADDR

# 使用容器內的 vault CLI
VAULT_CMD="docker exec vault vault"
if command -v vault >/dev/null 2>&1; then
    VAULT_CMD="vault"
fi

# 5. 檢查並建立 KV mount
echo ""
echo "[4/5] 檢查 KV mount..."
MOUNTS=$($VAULT_CMD secrets list -format=json 2>/dev/null || echo "{}")
KV_MOUNT=""

if echo "$MOUNTS" | grep -q '"kv"'; then
    KV_MOUNT="kv"
    echo "   ✅ 找到 KV mount: kv"
elif echo "$MOUNTS" | grep -q '"kv-v2"'; then
    KV_MOUNT="kv-v2"
    echo "   ✅ 找到 KV mount: kv-v2"
else
    echo "   ⚠️  未找到 KV mount，正在建立 kv-v2..."
    $VAULT_CMD secrets enable -version=2 -path=kv-v2 kv 2>/dev/null || {
        # 如果建立失敗，嘗試建立 kv
        $VAULT_CMD secrets enable -version=2 -path=kv kv 2>/dev/null || {
            echo "   ❌ 建立 KV mount 失敗"
            exit 1
        }
        KV_MOUNT="kv"
    }
    if [ -z "$KV_MOUNT" ]; then
        KV_MOUNT="kv-v2"
    fi
    echo "   ✅ 已建立 KV mount: $KV_MOUNT"
fi

# 6. 設定 policy
echo ""
echo "[5/5] 設定 Vault policy..."
if [ ! -f "$POLICY_FILE" ]; then
    echo "❌ Policy 檔案不存在：$POLICY_FILE"
    exit 1
fi

$VAULT_CMD policy write "$POLICY_NAME" "$POLICY_FILE" 2>/dev/null || {
    echo "❌ 寫入 policy 失敗"
    exit 1
}
echo "   ✅ Policy '$POLICY_NAME' 已寫入"

# 7. 確保 Transit mount 存在
echo ""
echo "[6/6] 檢查 Transit mount..."
if ! echo "$MOUNTS" | grep -q '"transit"'; then
    echo "   ⚠️  未找到 Transit mount，正在建立..."
    $VAULT_CMD secrets enable transit 2>/dev/null || {
        echo "   ❌ 建立 Transit mount 失敗"
        exit 1
    }
    echo "   ✅ 已建立 Transit mount"
else
    echo "   ✅ Transit mount 已存在"
fi

echo ""
echo "=========================================="
echo "  ✅ 修復完成！"
echo "=========================================="
echo ""
echo "📋 摘要："
echo "   - KV Mount: $KV_MOUNT"
echo "   - Policy: $POLICY_NAME"
echo "   - Root Token: 已使用（擁有所有權限）"
echo ""
echo "💡 下一步："
echo "   1. 確保 Go server 使用正確的 mount 名稱："
echo "      export VAULT_MOUNT=$KV_MOUNT"
echo ""
echo "   2. 或在 Go server 的 .env 檔案中設定："
echo "      VAULT_MOUNT=$KV_MOUNT"
echo "      VAULT_ADDR=$VAULT_ADDR"
echo "      VAULT_TOKEN=\$(cat $INIT_DIR/root_token.txt)"
echo ""
echo "   3. 重新啟動 Go server 並測試保險業者註冊"

