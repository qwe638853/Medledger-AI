#!/usr/bin/env bash
set -euo pipefail
# 7) 可選：對 CA 進行使用者 revoke + identity remove（需管理員身份）
# 環境變數：
#   DO_CA_REMOVE=1                啟用此步驟（預設 0）
#   CA_URL=http://localhost:7054  CA URL
#   CA_ADMIN_HOME=/path/to/msp    CA 管理員 MSP 目錄（需含 keystore/signcerts）
DO_CA_REMOVE=${DO_CA_REMOVE:-1}
CA_URL=${CA_URL:-http://localhost:7054}
CA_ADMIN_HOME=${CA_ADMIN_HOME:-/home/Medledger-AI/hyperledger/root-admin}
ROOT_DIR=${ROOT_DIR:-/home/Medledger-AI/hyperledger/go_server}
# 新增：支援從 SQLite 錢包表讀取身份（避免只依賴 msp-data */）
DB_PATH=${DB_PATH:-$ROOT_DIR/database/user_data.sqlite}

if [ "$DO_CA_REMOVE" = "1" ]; then
  echo "[INFO] 啟用 CA 憑證撤銷與身份註銷 (revoke + identity remove)"
  if ! command -v fabric-ca-client >/dev/null 2>&1; then
    echo "[WARN] 找不到 fabric-ca-client，略過 CA 註銷步驟"
  elif [ -z "$CA_ADMIN_HOME" ] || [ ! -d "$CA_ADMIN_HOME" ]; then
    echo "[WARN] 未設定 CA_ADMIN_HOME 或路徑不存在，略過 CA 註銷步驟"
  else
    export FABRIC_CA_CLIENT_HOME="$CA_ADMIN_HOME"
    echo "[INFO] 使用 CA 管理員 MSP：$FABRIC_CA_CLIENT_HOME"
    echo "[INFO] CA URL: $CA_URL"

    # 收集要註銷的身份
    IDS=()
    # 1) 從 SQLite 的 wallet 表讀取（若存在）
    if [ -f "$DB_PATH" ]; then
      echo "[INFO] 從 SQLite 讀取 wallet 標籤: $DB_PATH"
      while IFS= read -r id; do
        [ -n "$id" ] && IDS+=("$id")
      done < <(sqlite3 "$DB_PATH" "SELECT label FROM wallet;" 2>/dev/null || true)
    fi

    # 2) 從 msp-data 目錄推斷（相容舊流程）
    if [ -d "$ROOT_DIR/msp-data/users" ]; then
      while IFS= read -r -d '' d; do
        IDS+=("$(basename "$d")")
      done < <(find "$ROOT_DIR/msp-data/users" -mindepth 1 -maxdepth 1 -type d -print0)
    fi
    if [ -d "$ROOT_DIR/msp-data/insurers" ]; then
      while IFS= read -r -d '' d; do
        IDS+=("$(basename "$d")")
      done < <(find "$ROOT_DIR/msp-data/insurers" -mindepth 1 -maxdepth 1 -type d -print0)
    fi
    # 平台身份：無論是否存在 msp-data，都一併嘗試（由 CA 決定是否存在）
    IDS+=("platform")
    # 新增：診所身分（若存在）
    if [ -d "$ROOT_DIR/msp-data/clinic" ]; then
      while IFS= read -r -d '' d; do
        IDS+=("$(basename "$d")")
      done < <(find "$ROOT_DIR/msp-data/clinic" -mindepth 1 -maxdepth 1 -type d -print0)
    fi

    # 去重
    if [ ${#IDS[@]} -gt 0 ]; then
      mapfile -t IDS < <(printf "%s\n" "${IDS[@]}" | awk 'NF' | sort -u)
    fi

    if [ ${#IDS[@]} -eq 0 ]; then
      echo "[INFO] 未發現可註銷的身份，略過"
    else
      echo "[INFO] 將對下列身份執行 revoke + identity remove：${IDS[*]}"
      for id in "${IDS[@]}"; do
        echo "[INFO] revoking $id ..."
        fabric-ca-client revoke -u "$CA_URL" -e "$id" || true
        echo "[INFO] identity remove $id ..."
        fabric-ca-client identity remove "$id" -u "$CA_URL" || true
      done
    fi
  fi
fi

echo "[DONE] 清理完成（未動到任何憑證與 MSP 目錄）。"


