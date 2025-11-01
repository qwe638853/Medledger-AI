# 📊 鏈上與鏈下資料儲存說明

## 🎯 快速總覽

| 儲存位置 | 資料類型 | 說明 |
|---------|---------|------|
| **🟢 區塊鏈（Hyperledger Fabric）** | 健檢報告、授權記錄 | 不可篡改、可追蹤 |
| **🟡 SQLite 資料庫** | 用戶帳號、密碼、個人資訊 | 快速查詢、帳號管理 |
| **🔵 Vault** | 加密金鑰、憑證、私鑰 | 專業金鑰管理 |
| **🔴 記憶體/處理中** | 明文報告（加密前） | 僅在處理過程中存在 |

---

## 🔷 鏈上資料（Hyperledger Fabric 區塊鏈）

### 📄 1. 健檢報告（HealthReport）

**儲存位置**：`chaincode-go/health_contract.go`

**資料結構**：
```go
type HealthReport struct {
    DocType     string `json:"docType"`      // "HealthRecord"
    ReportID    string `json:"reportId"`      // 報告唯一ID，例如 "report_12345"
    PatientHash string `json:"patientHash"`  // 病患ID的SHA256雜湊值
    ClinicID    string `json:"clinicId"`      // 診所ID（從憑證中取得）
    ResultJSON  string `json:"resultJson"`   // ⚠️ 加密後的 Envelope JSON
    CreatedAt   int64  `json:"createdAt"`     // Unix 時間戳
}
```

**實際內容範例**：
```json
{
  "docType": "HealthRecord",
  "reportId": "report_12345",
  "patientHash": "a1b2c3d4...",  // SHA256("patient_001")
  "clinicId": "clinic-001",
  "resultJson": "{\"ciphertext\":\"base64...\",\"wrappedKeys\":{...}}",
  "createdAt": 1703123456
}
```

**⚠️ 重要說明**：
- ✅ **`resultJson` 是加密後的 Envelope**，包含：
  - `ciphertext`：AES-256-GCM 加密的報告內容
  - `wrappedKeys`：不同角色加密的 dataKey
  - `nonce`：加密用的隨機數
- ✅ **永遠不會在鏈上存明文報告**
- ✅ **任何人讀取區塊鏈都只能看到加密內容**

---

### 🎫 2. 授權票據（AuthTicket）

**用途**：記錄病患授權給誰查看報告

**資料結構**：
```go
type AuthTicket struct {
    DocType     string `json:"docType"`      // "AuthTicket"
    PatientHash string `json:"patientHash"`  // 病患雜湊
    TargetHash  string `json:"targetHash"`   // 被授權者雜湊（如保險公司）
    ReportID    string `json:"reportId"`      // 報告ID
    GrantedAt   int64  `json:"grantedAt"`    // 授權時間
    Expiry      int64  `json:"expiry"`       // 過期時間
}
```

**範例**：
```json
{
  "docType": "AuthTicket",
  "patientHash": "a1b2c3d4...",
  "targetHash": "e5f6g7h8...",  // 保險公司雜湊
  "reportId": "report_12345",
  "grantedAt": 1703123456,
  "expiry": 1703818256  // 7天後過期
}
```

---

### 📋 3. 授權請求（AccessRequest）

**用途**：記錄保險業者向病患申請授權的記錄

**資料結構**：
```go
type AccessRequest struct {
    DocType       string `json:"docType"`         // "AccessRequest"
    RequestID     string `json:"requestId"`      // 請求唯一ID
    ReportID      string `json:"reportId"`        // 報告ID
    PatientHash   string `json:"patientHash"`     // 病患雜湊
    RequesterHash string `json:"requesterHash"`  // 申請者雜湊（保險公司）
    Reason        string `json:"reason"`          // 申請理由
    RequestedAt   int64  `json:"requestedAt"`     // 申請時間
    Expiry        int64  `json:"expiry"`          // 請求過期時間
    Status        string `json:"status"`          // "pending"/"approved"/"rejected"
}
```

**範例**：
```json
{
  "docType": "AccessRequest",
  "requestId": "req_001",
  "reportId": "report_12345",
  "patientHash": "a1b2c3d4...",
  "requesterHash": "e5f6g7h8...",
  "reason": "理賠審核需要",
  "requestedAt": 1703123456,
  "expiry": 1703818256,
  "status": "pending"
}
```

---

### 🔍 鏈上查詢方法

**從鏈碼可以看到**：
```go
// 讀取單一報告
func ReadMyReport(ctx, reportID) → 返回 HealthReport

// 列出報告 metadata
func ListMyReportMeta(ctx) → 返回 ReportMeta[] (只有 ID、診所、時間)

// 查詢授權記錄
func GetAuthTicket(ctx, patientHash, targetHash, reportID) → 返回 AuthTicket
```

---

## 🔶 鏈下資料（SQLite 資料庫）

### 👤 1. 用戶表（users）

**儲存位置**：`database/user_data.sqlite`

**資料表結構**：
```sql
CREATE TABLE users (
    username TEXT PRIMARY KEY,  -- ⚠️ 儲存雜湊後的用戶ID（SHA256）
    password TEXT,               -- ⚠️ 儲存雜湊後的密碼（SHA256）
    name TEXT,                   -- 真實姓名
    date TEXT,                   -- 生日
    email TEXT,                   -- 電子郵件
    phone TEXT                    -- 電話號碼
);
```

**實際資料範例**：
```sql
-- 用戶ID: "patient_001"
-- 儲存時會先 SHA256 雜湊
username: "3f7a2b8c..."  (SHA256("patient_001"))
password: "9d4e1f2a..."  (SHA256("原始密碼"))
name: "張三"
date: "1990-01-01"
email: "zhang@example.com"
phone: "0912345678"
```

**⚠️ 隱私保護**：
- ✅ **username 和 password 都是雜湊值**，不是原始ID
- ✅ 無法從資料庫直接反推出真實身份證字號
- ✅ 需要用相同的雜湊演算法比對才能驗證身份

---

### 🏢 2. 保險業者表（insurers）

**資料表結構**：
```sql
CREATE TABLE insurers (
    insurer_id TEXT PRIMARY KEY,  -- ⚠️ 儲存雜湊後的保險業者ID
    password TEXT,                 -- ⚠️ 儲存雜湊後的密碼
    company_name TEXT,             -- 公司名稱
    contact_person TEXT,           -- 聯絡人
    email TEXT,                    -- 電子郵件
    phone TEXT                     -- 電話號碼
);
```

**範例**：
```sql
insurer_id: "5a6b7c8d..."  (SHA256("insurer_001"))
password: "1e2f3a4b..."
company_name: "XX人壽保險"
contact_person: "李四"
email: "insurer@example.com"
phone: "02-12345678"
```

---

### 🔑 3. 錢包表（wallet）

**用途**：儲存用戶在 Fabric 的身份引用

**資料表結構**：
```sql
CREATE TABLE wallet (
    label TEXT PRIMARY KEY,   -- 用戶ID（如 "patient_001"）
    content BLOB NOT NULL     -- JSON格式的引用資訊
);
```

**儲存內容範例**：
```json
{
  "mspID": "Org1MSP",
  "signerUri": "transit://user-patient_001",
  "certUri": "kv://users/patient_001"
}
```

**⚠️ 重要**：
- ✅ **不存實際憑證和私鑰**
- ✅ 只存**引用路徑**指向 Vault
- ✅ 實際憑證在 Vault 的 KV 儲存中

---

## 🔵 Vault 儲存（HashiCorp Vault）

### 🔐 1. 用戶憑證和私鑰（KV 儲存）

**路徑結構**：
```
/users/{userID}/
  - csr: 證書簽名請求
  - key: 私鑰（實際上空，因為用 Transit）
  - cert: X.509 憑證

/insurers/{insurerID}/
  - csr, key, cert

/clinics/{clinicID}/
  - csr, key, cert

/platform/
  - csr, key, cert
```

**範例路徑**：
```
/v1/kv/data/users/patient_001
/v1/kv/data/clinics/clinic-001
```

---

### 🔑 2. Transit 加密金鑰

**用途**：用於包裝 dataKey 的加密金鑰

**金鑰命名規則**：
```
{role}-{id}-wrap

範例：
- clinic-001-wrap     → 診所ID 001的包裝金鑰
- user-patient_001-wrap → 病患ID的包裝金鑰
- insurer-001-wrap    → 保險公司ID的包裝金鑰
```

**Vault API 路徑**：
```
/v1/transit/keys/{keyName}
/v1/transit/encrypt/{keyName}
/v1/transit/decrypt/{keyName}
```

**特點**：
- ✅ 金鑰**永遠不離開 Vault**
- ✅ 只有 Vault 知道實際金鑰值
- ✅ 支援金鑰輪換（不影響已加密資料）

---

## 🔴 記憶體/處理中資料（暫時性）

### ⚠️ 明文報告

**存在時機**：
1. 前端送達 Go Server（`req.TestResultsJson`）
2. 加密處理過程中
3. **加密完成後立即清除**

**流程**：
```
前端明文
  ↓
Go Server 接收
  ↓
TransitWrapper.EncryptReportTransitClinicOnly()
  ↓
加密成 Envelope JSON
  ↓
明文從記憶體清除 ❌
  ↓
只有加密內容上鏈 ✅
```

**安全措施**：
- ✅ 明文**永遠不上鏈**
- ✅ 明文**不存入資料庫**
- ✅ 明文**只在處理時短暫存在記憶體**

---

## 📊 完整資料流向圖

```
┌─────────────────────────────────────────────────────────────┐
│                      前端應用程式                             │
│  ┌──────────────────┐  ┌──────────────────┐              │
│  │ 明文健檢報告       │  │ 用戶ID、密碼       │              │
│  │ (僅處理時存在)    │  │ (雜湊後存SQLite)  │              │
│  └────────┬─────────┘  └────────┬─────────┘              │
└───────────┼─────────────────────┼────────────────────────────┘
            │                     │
            ▼                     ▼
┌─────────────────────────────────────────────────────────────┐
│                    Go Server (處理層)                         │
│  ┌──────────────────┐  ┌──────────────────┐              │
│  │ 1. 接收明文報告    │  │ 2. 驗證帳號密碼    │              │
│  │ 2. 加密成 Envelope│  │ 3. 查詢/寫入SQLite│              │
│  │ 3. 明文清除 ✅    │  │                  │              │
│  └────────┬─────────┘  └──────────────────┘              │
└───────────┼─────────────────────────────────────────────────┘
            │
            ├─────────────────┬─────────────────┐
            ▼                 ▼                 ▼
┌───────────────────┐  ┌────────────┐  ┌────────────┐
│ Hyperledger       │  │ SQLite DB  │  │ Vault       │
│ Fabric 區塊鏈      │  │            │  │             │
│                   │  │ users      │  │ Transit Keys│
│ HealthReport      │  │ insurers   │  │ Certificates│
│ (加密Envelope)    │  │ wallet     │  │ Private Keys│
│ AuthTicket        │  │            │  │             │
│ AccessRequest     │  │            │  │             │
└───────────────────┘  └────────────┘  └────────────┘
```

---

## 🎯 資料分類總結表

| 資料類型 | 儲存位置 | 是否加密 | 可訪問性 |
|---------|---------|---------|---------|
| **健檢報告內容** | 區塊鏈 | ✅ AES-256-GCM + Transit | 只有有權限的角色 |
| **報告 metadata** | 區塊鏈 | ❌ 明文（ID、時間等） | 有權限的角色 |
| **授權記錄** | 區塊鏈 | ❌ 明文（雜湊值） | 相關角色 |
| **用戶帳號密碼** | SQLite | ✅ SHA256 雜湊 | 系統內部 |
| **用戶個人資訊** | SQLite | ❌ 明文 | 系統內部 |
| **Fabric 憑證** | Vault KV | ❌ PEM 格式 | Vault 授權用戶 |
| **Transit Wrap Keys** | Vault Transit | ✅ Vault 內部加密 | Vault 管理 |
| **報告明文** | 記憶體 | ❌ 僅處理時 | 處理過程 |

---

## 🔒 隱私保護策略

### ✅ 實施的保護措施

1. **身份隱私**
   - ✅ 區塊鏈只存 SHA256 雜湊值，非真實身份證
   - ✅ SQLite 中的 username 也是雜湊

2. **資料隱私**
   - ✅ 報告內容加密才上鏈
   - ✅ 使用混合加密（AES + Transit）
   - ✅ 多重金鑰包裝（不同角色不同金鑰）

3. **金鑰管理**
   - ✅ 私鑰和加密金鑰存在 Vault（專業管理）
   - ✅ 錢包只存引用，不存實際金鑰
   - ✅ 金鑰永不離開 Vault

4. **可追蹤性**
   - ✅ 所有操作記錄在區塊鏈（不可篡改）
   - ✅ 授權記錄完整保存
   - ✅ 符合醫療資料法規要求

---

## ❓ 常見問題

### Q: 為什麼不把所有資料都存在區塊鏈？

**A**: 
- ❌ **成本考量**：區塊鏈儲存成本高
- ❌ **效能考量**：每次查詢都需要共識，速度慢
- ✅ **分層設計**：重要資料上鏈，查詢資料存 SQLite

### Q: 如果 SQLite 被駭了會怎樣？

**A**:
- ⚠️ 可以拿到用戶個人資訊（姓名、電話、email）
- ✅ **無法拿到健檢報告**（報告在區塊鏈，且加密）
- ✅ **無法登入系統**（密碼是雜湊，無法反推）
- ✅ **無法偽造身份**（需要 Vault 中的憑證才能操作鏈碼）

### Q: 為什麼要同時用 SQLite 和區塊鏈？

**A**:
- ✅ **SQLite**：快速查詢、帳號管理、會話管理
- ✅ **區塊鏈**：不可篡改、可追蹤、多方信任
- ✅ **互補**：各取所長

### Q: Vault 壞了怎麼辦？

**A**:
- ⚠️ **無法加密/解密新的報告**
- ⚠️ **無法簽署新的交易**
- ✅ **已加密的報告仍在區塊鏈**（資料不會丟失）
- ✅ 需要用 Vault 備份和金鑰恢復機制

---

## 📝 總結

### 🟢 區塊鏈（公開、不可篡改）
- ✅ 加密的健檢報告
- ✅ 授權記錄和請求
- ✅ 所有操作的歷史記錄

### 🟡 SQLite（快速、可查詢）
- ✅ 用戶帳號和密碼（雜湊）
- ✅ 用戶個人資訊
- ✅ 錢包引用資訊

### 🔵 Vault（安全、專業）
- ✅ 加密金鑰（Transit）
- ✅ 憑證和私鑰
- ✅ 金鑰管理服務

### 🔴 記憶體（暫時性）
- ✅ 明文報告（僅處理時）
- ✅ 處理後立即清除

這樣的設計確保了：
1. 🔒 **資料安全**：敏感資料加密
2. 🔍 **可追蹤性**：所有操作有記錄
3. ⚡ **效能**：快速查詢和處理
4. 🛡️ **隱私保護**：身份資訊雜湊處理



