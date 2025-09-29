#!/usr/bin/env bash
set -euo pipefail

# 清理 Hyperledger Fabric 節點資料（不動憑證/MSP）
# - 停止並移除 docker-compose 服務與資料卷（peer/orderer/couchdb/chaincode 容器資料）
# - 刪除通道產物（channel-artifacts/、system-genesis-block/）
# - 刪除 Explorer 的本地資料（不影響憑證）
# - 不刪除 organizations/*、msp-data/*、各類 MSP/證書目錄
# - 可選：以 CA 管理員身分對 CA 進行 revoke + identity remove（需設定環境變數）

ROOT_DIR="/home/Medledger-AI/hyperledger"
COMPOSE_MAIN="$ROOT_DIR/docker-compose.yaml"
COMPOSE_EXPLORER="$ROOT_DIR/explorer/docker-compose.yaml"

echo "[INFO] 清理環境：$ROOT_DIR"

cd "$ROOT_DIR"

# 1) 關閉並移除主要網路容器與資料卷
if [ -f "$COMPOSE_MAIN" ]; then
  echo "[INFO] docker compose down (main) ..."
  if command -v docker compose >/dev/null 2>&1; then
    docker compose -f "$COMPOSE_MAIN" down -v --remove-orphans || true
  else
    docker-compose -f "$COMPOSE_MAIN" down -v --remove-orphans || true
  fi
else
  echo "[WARN] 找不到 $COMPOSE_MAIN，略過主要網路關閉"
fi

# 2) 關閉並移除 Explorer 容器與資料卷（若存在）
if [ -f "$COMPOSE_EXPLORER" ]; then
  echo "[INFO] docker compose down (explorer) ..."
  if command -v docker compose >/dev/null 2>&1; then
    docker compose -f "$COMPOSE_EXPLORER" down -v --remove-orphans || true
  else
    docker-compose -f "$COMPOSE_EXPLORER" down -v --remove-orphans || true
  fi
fi

# 3) 刪除 dev-peer* 連鎖碼容器與映像（安全：不影響基底映像與憑證）
echo "[INFO] 刪除 dev-peer* 連鎖碼容器/映像 ..."
DEV_PEER_CONTAINERS=$(docker ps -a --format '{{.ID}} {{.Names}}' | awk '/dev-peer/ {print $1}') || true
if [ -n "${DEV_PEER_CONTAINERS:-}" ]; then
  echo "$DEV_PEER_CONTAINERS" | xargs -r docker rm -f || true
fi
DEV_PEER_IMAGES=$(docker images --format '{{.Repository}}:{{.Tag}} {{.ID}}' | awk '/dev-peer/ {print $2}') || true
if [ -n "${DEV_PEER_IMAGES:-}" ]; then
  echo "$DEV_PEER_IMAGES" | xargs -r docker rmi -f || true
fi

# 4) 刪除通道產物與生成塊（不含 MSP 與憑證）
echo "[INFO] 刪除 channel-artifacts/ 與 system-genesis-block/ ..."
rm -rf "$ROOT_DIR/channel-artifacts" || true
rm -rf "$ROOT_DIR/system-genesis-block" || true

# 5) 刪除 Explorer 本地資料（不含憑證）
if [ -d "$ROOT_DIR/explorer" ]; then
  echo "[INFO] 清理 Explorer 本地資料 ..."
  rm -rf "$ROOT_DIR/explorer/.data" || true
  rm -rf "$ROOT_DIR/explorer/pgdata" || true
fi

# 6) 額外安全提示：不刪除下列敏感目錄
echo "[INFO] 保留目錄（未刪除）："
echo "       $ROOT_DIR/organizations (MSP/憑證)"
echo "       $ROOT_DIR/msp-data (本地身份憑證)"
echo "       $ROOT_DIR/root-ca (CA 憑證與資料庫)"
