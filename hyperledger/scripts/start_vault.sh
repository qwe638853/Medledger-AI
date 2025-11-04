#!/usr/bin/env bash
set -euo pipefail

# 啟動 Vault（正式模式）並自動初始化/解封（若需要）
# - 依賴 docker compose 與 curl
# - 初次執行會將 unseal key 與 root token 寫入 vault/init/*（權限 600）

ROOT_DIR="/home/Medledger-AI/hyperledger"
COMPOSE_FILE="$ROOT_DIR/docker-compose.yaml"
INIT_DIR="$ROOT_DIR/vault/init"
VAULT_HTTP="http://127.0.0.1:8200"


echo "[INFO] 等待 Vault API 埠就緒: $VAULT_HTTP"
for i in {1..60}; do
  if curl -sS "$VAULT_HTTP/v1/sys/health" -o /dev/null; then
    break
  fi
  sleep 1
  if [ $i -eq 60 ]; then
    echo "[ERROR] Vault API 未就緒" >&2
    exit 1
  fi
done

code=$(curl -s -o /dev/null -w '%{http_code}' "$VAULT_HTTP/v1/sys/health") || code=000
echo "[INFO] 健康檢查狀態碼: $code (200=運行中, 503=已初始化但封印, 501=未初始化)"

mkdir -p "$INIT_DIR"
chmod 700 "$INIT_DIR"

if [ "$code" = "501" ]; then
  echo "[INFO] Vault 未初始化，開始初始化 ..."
  # 使用容器內 vault CLI 初始化並輸出 JSON
  out=$(docker exec vault sh -lc 'vault operator init -key-shares=1 -key-threshold=1 -format=json')
  echo "$out" > "$INIT_DIR/init.json"
  chmod 600 "$INIT_DIR/init.json"
  unseal_key=$(echo "$out" | jq -r '.unseal_keys_b64[0]')
  root_token=$(echo "$out" | jq -r '.root_token')
  printf "%s\n" "$unseal_key" > "$INIT_DIR/unseal_key.txt"
  printf "%s\n" "$root_token" > "$INIT_DIR/root_token.txt"
  chmod 600 "$INIT_DIR/unseal_key.txt" "$INIT_DIR/root_token.txt"
  echo "[INFO] 完成初始化，已將金鑰寫入 $INIT_DIR"

  echo "[INFO] 解封 Vault ..."
  docker exec vault sh -lc "vault operator unseal \"$unseal_key\"" >/dev/null
  code=$(curl -s -o /dev/null -w '%{http_code}' "$VAULT_HTTP/v1/sys/health") || code=000
fi

if [ "$code" = "503" ]; then
  echo "[INFO] Vault 已初始化但封印，嘗試使用 $INIT_DIR/unseal_key.txt 解封 ..."
  if [ -f "$INIT_DIR/unseal_key.txt" ]; then
    unseal_key=$(cat "$INIT_DIR/unseal_key.txt")
    docker exec vault sh -lc "vault operator unseal \"$unseal_key\"" >/dev/null || true
    code=$(curl -s -o /dev/null -w '%{http_code}' "$VAULT_HTTP/v1/sys/health") || code=000
  else
    echo "[WARN] 找不到 $INIT_DIR/unseal_key.txt，請手動解封：docker exec -it vault vault operator unseal" >&2
  fi
fi

if [ "$code" != "200" ]; then
  echo "[WARN] Vault 仍未處於運行狀態（HTTP $code）。請檢查容器日誌：docker logs vault" >&2
else
  echo "[INFO] Vault 已啟動且可用：$VAULT_HTTP"
  if [ -f "$INIT_DIR/root_token.txt" ]; then
    echo "[INFO] Root Token 存於：$INIT_DIR/root_token.txt"
    echo "[HINT] 使用方式：export VAULT_ADDR=$VAULT_HTTP && export VAULT_TOKEN=\"$(cat $INIT_DIR/root_token.txt)\""
  fi
fi


