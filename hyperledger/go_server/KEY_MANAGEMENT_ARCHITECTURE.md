# 🔑 金鑰管理架構說明

## 🎯 您的問題確認

**問題**：現在是否將 CA 存進 KSM，並用 Transit 進行加密 dataKey？

**答案**：✅ **部分正確！** 讓我詳細說明：

---

## 📊 實際架構總覽

### ✅ 1. CA 憑證儲存（在 Vault KV）

**是的，CA 簽發的憑證存進 Vault（KSM）**

從代碼可以看到：

```119:128:service/user_service.go
    // ✅ 寫入 Vault + 錢包改為僅存引用（DB 不存證書與私鑰）
    if store, verr := vs.NewFromEnv(); verr != nil {
        log.Printf("[Vault] 初始化失敗（略過）：%v", verr)
    } else {
        if werr := store.WriteUserMaterial(ctx, req.UserId, csrPEM, keyPEM, certPem); werr != nil {
            log.Printf("[Vault] 寫入使用者材料失敗：%v", werr)
        } else {
            log.Printf("[Vault] ✅ 已寫入使用者材料至 Vault：%s", req.UserId)
        }
    }
```

**儲存內容**：
```go
WriteUserMaterial(userID, csrPEM, keyPEM, certPEM)
```

**實際儲存到 Vault KV**：
```
/v1/kv/data/users/{userID}
{
  "data": {
    "csr": "CSR 內容（PEM格式）",
    "key": "",  // ⚠️ 實際上為空字串！
    "cert": "Fabric CA 簽發的憑證（PEM格式）"
  }
}
```

**路徑結構**：
- 用戶：`/v1/kv/data/users/{userID}`
- 診所：`/v1/kv/data/clinics/{clinicID}`
- 保險業者：`/v1/kv/data/insurers/{insurerID}`
- 平台：`/v1/kv/data/platform`

---

### ✅ 2. DataKey 用 Transit 加密

**是的，dataKey 用 Vault Transit 加密**

從加密流程可以看到：

```56:62:secure/wrap/transit.go
    // 診所端使用獨立 wrap key（若傳入的是簽章 key 名稱，則補 -wrap）
    wrapClinicKey := clinicKey
    if !strings.HasSuffix(wrapClinicKey, "-wrap") { wrapClinicKey = wrapClinicKey + "-wrap" }
    if err := w.store.EnsureTransitKeyOfType(ctx, wrapClinicKey, "aes256-gcm96"); err != nil { return nil, err }

    // 以 Transit 包 dataKey（僅診所）
    ctClinic, err := w.store.TransitEncrypt(ctx, wrapClinicKey, dk)
```

**實際調用**：
```go
store.TransitEncrypt(ctx, "clinic-001-wrap", dataKey)
// → 調用 Vault API: POST /v1/transit/encrypt/clinic-001-wrap
// → 返回: "vault:v1:加密後的dataKey"
```

---

### ⚠️ 3. 重要差異：私鑰處理方式

**關鍵發現**：系統**不使用傳統私鑰檔案**！

從代碼可以看到：

```106:107:service/user_service.go
    csrPEM := pem.EncodeToMemory(&pem.Block{Type:"CERTIFICATE REQUEST", Bytes: csrDER})
    keyPEM := []byte("")  // ⚠️ 空字串！
```

**實際架構**：
- ❌ **不使用傳統 ECDSA 私鑰檔案**
- ✅ **使用 Vault Transit 簽章金鑰** 代替私鑰

**簽章流程**：

```129:130:service/user_service.go
    // signerUri 使用使用者專屬 Transit key：user-<id>
    err = wallet.PutReference(req.UserId, "Org1MSP", "transit://user-"+req.UserId, "kv://users/"+req.UserId)
```

**簽章時**：

```131:134:wallet/wallet.go
    keyName := strings.TrimPrefix(payload.SignerURI, "transit://")
    pubKey, _ := cert.PublicKey.(*ecdsa.PublicKey)
    fabricSigner := func(digest []byte) ([]byte, error) {
        der, err := store.TransitSign(context.Background(), keyName, digest)
```

---

## 🔐 完整架構圖

```
┌─────────────────────────────────────────────────────────────┐
│                    Fabric CA                                │
│  ┌────────────────────────────────────────────────────┐     │
│  │ 1. 用戶產生 CSR（用 Transit 公鑰）                  │     │
│  │ 2. CA 簽發憑證                                     │     │
│  │ 3. 憑證存回 Vault                                   │     │
│  └────────────┬───────────────────────────────────────┘     │
└───────────────┼──────────────────────────────────────────────┘
                │ 憑證 (certPEM)
                ▼
┌─────────────────────────────────────────────────────────────┐
│              Vault (KSM)                                    │
│  ┌──────────────────┐  ┌──────────────────┐              │
│  │ KV Storage       │  │ Transit Engine   │              │
│  │                  │  │                  │              │
│  │ /users/{id}      │  │ Sign Keys:       │              │
│  │  - csr           │  │  - user-{id}     │              │
│  │  - key: ""       │  │  - clinic-{id}  │              │
│  │  - cert ✅       │  │  - insurer-{id}  │              │
│  │                  │  │                  │              │
│  │ /clinics/{id}    │  │ Wrap Keys:       │              │
│  │  - csr           │  │  - user-{id}-wrap│              │
│  │  - key: ""       │  │  - clinic-{id}-wrap ✅          │
│  │  - cert ✅       │  │  - insurer-{id}-wrap            │
│  └──────────────────┘  └──────────────────┘              │
└─────────────────────────────────────────────────────────────┘
                              │
                ┌─────────────┴─────────────┐
                │                           │
                ▼                           ▼
        ┌───────────────┐          ┌──────────────┐
        │ 簽章交易      │          │ 加密 DataKey  │
        │ (TransitSign) │          │ (TransitEncrypt)│
        └───────────────┘          └──────────────┘
```

---

## 📋 詳細說明各組件

### 1. CA 憑證管理流程

```
用戶註冊
  ↓
產生 CSR（用 Transit 公鑰）
  ↓
提交給 Fabric CA
  ↓
CA 簽發憑證
  ↓
寫入 Vault KV ✅
  ├─ /users/{id}/cert    ← Fabric CA 簽發的憑證
  ├─ /users/{id}/csr     ← 證書簽名請求
  └─ /users/{id}/key     ← 空字串（不使用）
```

**關鍵點**：
- ✅ CA 憑證存在 Vault KV 中
- ✅ 不需要管理私鑰檔案
- ✅ 錢包只存引用路徑（`kv://users/{id}`）

---

### 2. Transit 雙重用途

#### 🔐 用途 1：簽章金鑰（代替私鑰）

**金鑰名稱**：`user-{id}`, `clinic-{id}`, `insurer-{id}`

**用於**：
- Fabric 交易簽章
- CSR 生成（簽署 CSR）

**API 調用**：
```go
// 簽章
store.TransitSign(ctx, "user-patient_001", digest)
// → POST /v1/transit/sign/user-patient_001
```

**特點**：
- ✅ 金鑰類型：`ecdsa-p256`
- ✅ 用於數位簽章，不是加密
- ✅ 代替傳統私鑰檔案

---

#### 🔑 用途 2：Wrap 金鑰（加密 dataKey）

**金鑰名稱**：`user-{id}-wrap`, `clinic-{id}-wrap`, `insurer-{id}-wrap`

**用於**：
- 加密報告的 dataKey
- 包裝/解包 dataKey

**API 調用**：
```go
// 加密 dataKey
store.TransitEncrypt(ctx, "clinic-001-wrap", dataKey)
// → POST /v1/transit/encrypt/clinic-001-wrap

// 解密 dataKey
store.TransitDecrypt(ctx, "user-patient_001-wrap", encryptedDataKey)
// → POST /v1/transit/decrypt/user-patient_001-wrap
```

**特點**：
- ✅ 金鑰類型：`aes256-gcm96`
- ✅ 用於對稱加密（包裝 dataKey）
- ✅ 金鑰命名規則：`{role}-{id}-wrap`

---

## 🎯 完整資料流向

### 場景：診所上傳報告

```
1. 建立 Fabric 交易
   │
   ├─ 取得憑證（從 Vault KV）
   │  └─ /v1/kv/data/users/clinic-001
   │
   └─ 簽章交易（用 Transit 簽章金鑰）
      └─ /v1/transit/sign/clinic-001  ✅

2. 加密報告
   │
   ├─ 產生 dataKey（32 bytes）
   │
   ├─ AES-256-GCM 加密報告內容
   │
   └─ 用 Transit wrap key 加密 dataKey
      └─ /v1/transit/encrypt/clinic-001-wrap  ✅

3. 打包 Envelope
   └─ {
        "ciphertext": "加密的報告",
        "wrappedKeys": {
          "clinic": "vault:v1:加密的dataKey"
        }
      }

4. 提交到區塊鏈
   └─ Hyperledger Fabric
```

---

## 📊 金鑰對照表

| 用途 | Transit 金鑰名稱 | 金鑰類型 | API 端點 | 儲存位置 |
|------|----------------|---------|---------|---------|
| **簽章交易** | `user-{id}` | `ecdsa-p256` | `/transit/sign/` | Vault Transit |
| **簽章交易** | `clinic-{id}` | `ecdsa-p256` | `/transit/sign/` | Vault Transit |
| **簽章交易** | `insurer-{id}` | `ecdsa-p256` | `/transit/sign/` | Vault Transit |
| **加密 dataKey** | `user-{id}-wrap` | `aes256-gcm96` | `/transit/encrypt/` | Vault Transit |
| **加密 dataKey** | `clinic-{id}-wrap` | `aes256-gcm96` | `/transit/encrypt/` | Vault Transit |
| **加密 dataKey** | `insurer-{id}-wrap` | `aes256-gcm96` | `/transit/encrypt/` | Vault Transit |
| **Fabric 憑證** | N/A | X.509 憑證 | N/A | Vault KV |

---

## ✅ 總結：您的架構確認

### ✅ 正確的部分

1. **CA 憑證存進 Vault（KSM）** ✅
   - 儲存在 Vault KV：`/v1/kv/data/users/{id}/cert`
   - 包含 CSR 和 Fabric CA 簽發的憑證

2. **DataKey 用 Transit 加密** ✅
   - 使用 Transit wrap keys：`{role}-{id}-wrap`
   - 類型：`aes256-gcm96`
   - 用於包裝報告的 dataKey

### 🔍 額外說明

3. **私鑰處理** ⚠️
   - **不使用傳統私鑰檔案**
   - 使用 **Transit 簽章金鑰** 代替私鑰
   - 金鑰名稱：`user-{id}`, `clinic-{id}`（沒有 `-wrap` 後綴）
   - 類型：`ecdsa-p256`
   - 用於 Fabric 交易簽章

---

## 🔑 金鑰命名規則總覽

```
簽章金鑰（用於交易簽章）：
  - user-{userId}           → 病患的簽章金鑰
  - clinic-{clinicId}      → 診所的簽章金鑰
  - insurer-{insurerId}    → 保險業者的簽章金鑰

Wrap 金鑰（用於加密 dataKey）：
  - user-{userId}-wrap           → 病患的加密金鑰
  - clinic-{clinicId}-wrap       → 診所的加密金鑰
  - insurer-{insurerId}-wrap     → 保險業者的加密金鑰
```

---

## 🎓 架構優勢

1. **統一金鑰管理**
   - ✅ 所有金鑰都在 Vault
   - ✅ 不需要管理私鑰檔案
   - ✅ 支援金鑰輪換

2. **安全性**
   - ✅ 私鑰永遠不離開 Vault
   - ✅ 金鑰有完整審計日誌
   - ✅ 支援細粒度權限控制

3. **靈活性**
   - ✅ 簽章和加密分離（不同金鑰）
   - ✅ 可以獨立管理不同角色的金鑰
   - ✅ 支援動態添加收件者

---

## ❓ 常見問題

### Q: 為什麼不用傳統私鑰檔案？

**A**: 
- ✅ **更安全**：私鑰永遠不離開 Vault
- ✅ **更易管理**：不需要處理 PEM 檔案
- ✅ **支援輪換**：可以自動輪換金鑰
- ✅ **審計完整**：Vault 記錄所有簽章操作

### Q: 簽章金鑰和 wrap 金鑰有什麼區別？

**A**:
| 特性 | 簽章金鑰 | Wrap 金鑰 |
|------|---------|----------|
| **用途** | Fabric 交易簽章 | 加密 dataKey |
| **類型** | `ecdsa-p256`（非對稱） | `aes256-gcm96`（對稱） |
| **命名** | `user-{id}` | `user-{id}-wrap` |
| **API** | `/transit/sign/` | `/transit/encrypt/` |

### Q: Vault KV 中的 `key` 欄位為什麼是空的？

**A**: 
- 因為不使用傳統私鑰檔案
- 簽章直接用 Transit 簽章金鑰
- 所以 `keyPEM := []byte("")` 是空字串
- 真正的簽章能力在 Transit 中（`user-{id}` 金鑰）

---

## 📝 完整流程示例

### 示例：用戶註冊

```go
1. 產生 Transit 簽章金鑰（在 Vault）
   → 確保存在：user-patient_001 (ecdsa-p256)

2. 取得 Transit 公鑰
   → 用於生成 CSR

3. 產生 CSR
   → 用 TransitSigner 簽署 CSR
   → 公鑰來自 Transit

4. 提交給 Fabric CA
   → CA 簽發憑證

5. 寫入 Vault KV
   → /v1/kv/data/users/patient_001
   → {
        "csr": "CSR PEM",
        "key": "",  // 空字串
        "cert": "CA 簽發的憑證 PEM"
      }

6. 寫入錢包引用
   → SQLite wallet 表
   → {
        "signerUri": "transit://user-patient_001",
        "certUri": "kv://users/patient_001"
      }

7. 產生 wrap 金鑰（第一次加密時自動產生）
   → 確保存在：user-patient_001-wrap (aes256-gcm96)
```

---

## 🎯 結論

您的理解**基本正確**，但需要注意：

✅ **CA 憑證** → Vault KV 儲存  
✅ **DataKey 加密** → Vault Transit wrap keys  
⚠️ **私鑰** → 不是傳統檔案，而是 **Transit 簽章金鑰**

這是一個**混合架構**：
- **身份管理**：Fabric CA + Vault KV（憑證儲存）
- **簽章**：Vault Transit 簽章金鑰（代替私鑰）
- **資料加密**：Vault Transit wrap keys（加密 dataKey）

這樣的設計提供了**統一的金鑰管理**和**更高的安全性**！



