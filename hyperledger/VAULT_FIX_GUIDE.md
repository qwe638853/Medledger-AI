# Vault 權限修復指南

## 問題描述

保險業者註冊時，Go server 嘗試將憑證寫入 Vault 時收到 403 錯誤：
```
[Vault] 寫入保險業者材料失敗：Error making API request.
URL: PUT http://127.0.0.1:8200/v1/kv-v2/data/insurers/FFFFF1
Code: 403. Errors: permission denied
```

這表示 Vault token 沒有寫入 `kv-v2/data/insurers/*` 的權限。

## 解決方案

### 方法 1: 使用 Root Token（最簡單）

如果 Vault 使用 root token，root token 擁有所有權限，不需要額外設定 policy。

1. **找到 root token**：
   ```bash
   # 檢查是否有 root token 檔案
   cat hyperledger/vault/init/root_token.txt
   
   # 或從 Vault 初始化輸出中取得
   ```

2. **設定環境變數**：
   ```bash
   export VAULT_ADDR=http://127.0.0.1:8200
   export VAULT_TOKEN=<你的_root_token>
   export VAULT_MOUNT=kv-v2  # 如果 Vault 使用 kv-v2 mount
   ```

3. **在 Go server 啟動時使用這些環境變數**：
   ```bash
   cd hyperledger/go_server
   export VAULT_ADDR=http://127.0.0.1:8200
   export VAULT_TOKEN=<你的_root_token>
   export VAULT_MOUNT=kv-v2
   go run main.go
   ```

### 方法 2: 設定 Policy（推薦用於生產環境）

1. **確保 Vault 已初始化並解封**：
   ```bash
   cd hyperledger
   bash scripts/start_vault.sh
   ```

2. **設定 policy**：
   ```bash
   # 使用 root token
   export VAULT_ADDR=http://127.0.0.1:8200
   export VAULT_TOKEN=<你的_root_token>
   
   # 寫入 policy
   vault policy write medledger-policy hyperledger/vault/config/policy.hcl
   ```

3. **檢查 KV mount 名稱**：
   ```bash
   vault secrets list
   # 應該會看到 kv 或 kv-v2
   ```

4. **確保 mount 存在**（如果不存在）：
   ```bash
   # 如果看到 kv-v2，確保它存在
   vault secrets enable -version=2 -path=kv-v2 kv
   
   # 或如果使用 kv
   vault secrets enable -version=2 -path=kv kv
   ```

5. **建立 token 並附加 policy**（可選，如果不想用 root token）：
   ```bash
   vault token create -policy=medledger-policy
   # 使用輸出的 token 作為 VAULT_TOKEN
   ```

### 方法 3: 檢查並修復現有設定

1. **檢查當前使用的 mount**：
   從 Go server 的 log 可以看到：
   ```
   [Vault] init client addr=... kvMount=kv-v2 transitMount=transit
   ```
   這表示 Go server 使用的是 `kv-v2` mount。

2. **確保 policy 包含 kv-v2 路徑**：
   `hyperledger/vault/config/policy.hcl` 已經更新，包含 `kv-v2` 的路徑。

3. **設定 policy**：
   ```bash
   export VAULT_ADDR=http://127.0.0.1:8200
   export VAULT_TOKEN=<你的_root_token>
   vault policy write medledger-policy hyperledger/vault/config/policy.hcl
   ```

4. **驗證 policy**：
   ```bash
   vault policy read medledger-policy
   # 應該會看到 kv-v2/data/insurers/* 的路徑
   ```

## 快速修復腳本

如果 Vault 已經初始化，可以直接執行：

```bash
cd hyperledger

# 1. 讀取 root token（如果存在）
export VAULT_TOKEN=$(cat vault/init/root_token.txt 2>/dev/null || echo "")

# 2. 如果沒有 root token，需要手動設定
if [ -z "$VAULT_TOKEN" ]; then
    echo "請手動設定 VAULT_TOKEN 環境變數"
    exit 1
fi

export VAULT_ADDR=http://127.0.0.1:8200

# 3. 確保 kv-v2 mount 存在
docker exec vault vault secrets enable -version=2 -path=kv-v2 kv 2>/dev/null || true

# 4. 設定 policy
docker exec vault vault policy write medledger-policy /vault/config/policy.hcl

# 5. 驗證
docker exec vault vault policy read medledger-policy
```

## 驗證修復

1. **重新啟動 Go server**（確保使用正確的環境變數）：
   ```bash
   cd hyperledger/go_server
   export VAULT_ADDR=http://127.0.0.1:8200
   export VAULT_TOKEN=<你的_token>
   export VAULT_MOUNT=kv-v2  # 根據實際 mount 名稱設定
   go run main.go
   ```

2. **測試保險業者註冊**：
   嘗試註冊一個新的保險業者，應該不會再出現 403 錯誤。

3. **檢查 Vault 中是否有資料**：
   ```bash
   export VAULT_ADDR=http://127.0.0.1:8200
   export VAULT_TOKEN=<你的_token>
   vault kv get kv-v2/insurers/FFFFF1
   ```

## 常見問題

### Q: 如何知道 Vault 使用哪個 mount 名稱？
A: 從 Go server 的 log 可以看到 `[Vault] init client ... kvMount=...`，或執行 `vault secrets list`。

### Q: 為什麼會出現 403 錯誤？
A: 因為 Vault token 沒有寫入 `kv-v2/data/insurers/*` 的權限。需要：
   - 使用 root token（擁有所有權限）
   - 或設定 policy 並附加到 token

### Q: Policy 已經更新了，為什麼還是不行？
A: Policy 檔案更新後，需要重新寫入 Vault：
   ```bash
   vault policy write medledger-policy hyperledger/vault/config/policy.hcl
   ```

### Q: 如何確認 token 有正確的權限？
A: 使用 `vault token capabilities <path>` 檢查：
   ```bash
   vault token capabilities kv-v2/data/insurers/test
   # 應該顯示: create, update, read, delete, list
   ```

