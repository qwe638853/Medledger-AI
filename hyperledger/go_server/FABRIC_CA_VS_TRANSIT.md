# 🔗 Fabric CA 與 Transit 的關係說明

## 🤔 您的問題

**Fabric CA 與 Transit 有關係嗎？**

**簡單答案**：✅ **有關係，但它們是不同的系統，互相協作！**

---

## 📊 它們分別是什麼？

### 🏛️ Fabric CA（Certificate Authority）

**定位**：**發放憑證的機構**（就像真實世界的戶政事務所）

**功能**：
- ✅ **註冊身份**：註冊新用戶到系統中
- ✅ **簽發憑證**：發放 X.509 數位憑證
- ✅ **身份認證**：驗證用戶身份

**特點**：
- 🔐 只負責「發放憑證」，不管私鑰
- 📜 只關心「身份」和「憑證」
- ⚙️ 運行在：`http://localhost:7054`

---

### 🔑 Transit（Vault Transit Engine）

**定位**：**金鑰管理服務**（就像銀行的保險庫）

**功能**：
- ✅ **產生和儲存金鑰**：簽章金鑰、加密金鑰
- ✅ **簽章**：代替傳統私鑰檔案簽名
- ✅ **加密/解密**：包裝 dataKey

**特點**：
- 🔐 管理所有「私鑰和金鑰」
- 🚫 **永遠不讓金鑰離開** Vault
- ⚙️ 運行在：Vault 伺服器中的 Transit Engine

---

## 🔗 它們如何協作？

### 📝 完整流程：用戶註冊

讓我們看實際代碼：

```91:130:service/user_service.go
    // ✅ 以 TransitSigner 產生 CSR（私鑰不出庫）
    log.Printf("[DEBUG] 開始 Transit 產生 CSR，用戶ID: %s", req.UserId)
    store, err := vs.NewFromEnv()
    if err != nil { return &pb.RegisterResponse{Success:false, Message:"Vault 初始化失敗"}, nil }
    // 確保 Transit key 存在
    if err := store.EnsureTransitKey(ctx, "user-"+req.UserId); err != nil {
        return &pb.RegisterResponse{Success:false, Message:"建立使用者 Transit 金鑰失敗"}, nil
    }
    pub, err := store.TransitGetPublicKey(ctx, "user-"+req.UserId)
    if err != nil { return &pb.RegisterResponse{Success:false, Message:"取得公鑰失敗"}, nil }
    signerObj, err := sg.NewTransitSignerWithPublicKey(store, "user-"+req.UserId, pub)
    if err != nil { return &pb.RegisterResponse{Success:false, Message:"建立 Transit 簽章器失敗"}, nil }
    tmpl := x509.CertificateRequest{ Subject: pkix.Name{ CommonName: req.UserId } }
    csrDER, err := x509.CreateCertificateRequest(rand.Reader, &tmpl, signerObj)
    if err != nil { return &pb.RegisterResponse{Success:false, Message:"CSR 產生失敗"}, nil }
    csrPEM := pem.EncodeToMemory(&pem.Block{Type:"CERTIFICATE REQUEST", Bytes: csrDER})
    keyPEM := []byte("")

	// ✅ Enroll 產生證書
	enrollReq := fc.EnrollRequest{
		Certificate_request: string(csrPEM),
	}

    certPem, err := fc.EnrollUser("http://localhost:7054", req.UserId, req.Password, enrollReq)
	if err != nil {
		log.Fatalf("Enroll 失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "Enroll 憑證註冊失敗"}, nil
	}
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
    // signerUri 使用使用者專屬 Transit key：user-<id>
    err = wallet.PutReference(req.UserId, "Org1MSP", "transit://user-"+req.UserId, "kv://users/"+req.UserId)
```

**詳細步驟拆解**：

---

## 📋 步驟 1️⃣：Transit 產生簽章金鑰

```go
// 確保 Transit 簽章金鑰存在
store.EnsureTransitKey(ctx, "user-"+req.UserId)
// → 在 Vault 建立：user-patient_001 (ecdsa-p256)
```

**比喻**：在銀行開一個保險箱，專屬於這個用戶

**結果**：
- ✅ Vault 中有一個簽章金鑰：`user-patient_001`
- ✅ 類型：`ecdsa-p256`（非對稱加密，有公鑰和私鑰）
- ✅ 私鑰**永遠不離開** Vault

---

## 📋 步驟 2️⃣：從 Transit 取得公鑰

```go
pub, err := store.TransitGetPublicKey(ctx, "user-"+req.UserId)
// → 從 Vault 取得 public key
```

**比喻**：從保險箱取得「公鑰」，這個可以公開給別人

**結果**：
- ✅ 取得 Transit 簽章金鑰的**公鑰**
- ✅ 公鑰是公開的，不敏感

---

## 📋 步驟 3️⃣：用 Transit 公鑰產生 CSR

```go
signerObj, err := sg.NewTransitSignerWithPublicKey(store, "user-"+req.UserId, pub)
csrDER, err := x509.CreateCertificateRequest(rand.Reader, &tmpl, signerObj)
// → 用 TransitSigner 簽署 CSR
```

**關鍵**：CSR 是用 Transit 簽章金鑰簽名的，不是傳統私鑰檔案！

**比喻**：
- 用保險箱的鑰匙（Transit）簽署一份「申請書」（CSR）
- 申請書上寫著：「我是 patient_001，請給我憑證」
- 這份申請書必須用 Transit 簽章，證明身份

**結果**：
- ✅ 產生 CSR（Certificate Signing Request）
- ✅ CSR 包含 Transit 的公鑰
- ✅ CSR 用 Transit 私鑰簽名（但私鑰不離開 Vault）

---

## 📋 步驟 4️⃣：先註冊到 Fabric CA

```go
fc.RegisterUser(
    "http://localhost:7054",
    ...,
    api.RegistrationRequest{
        Name: req.UserId,
        Secret: req.Password,
        ...
    },
)
```

**比喻**：向戶政事務所（Fabric CA）登記名字，拿到一個註冊序號

**結果**：
- ✅ 在 Fabric CA 註冊身份
- ✅ 獲得 `userId` 和 `password`
- ⚠️ 但**還沒有憑證**

---

## 📋 步驟 5️⃣：提交 CSR 給 Fabric CA，取得憑證

```go
certPem, err := fc.EnrollUser("http://localhost:7054", req.UserId, req.Password, enrollReq)
// → 把 CSR 送給 Fabric CA，CA 簽發憑證
```

**關鍵流程**：
```
CSR（包含 Transit 公鑰）
  ↓
送給 Fabric CA
  ↓
Fabric CA 驗證 CSR 簽章（用 Transit 公鑰驗證）
  ↓
Fabric CA 簽發憑證（包含 Transit 公鑰）
  ↓
返回憑證
```

**比喻**：
- 把申請書（CSR）交給戶政事務所（Fabric CA）
- 戶政事務所驗證：申請書確實是用保險箱鑰匙簽的
- 戶政事務所發給您身份證（憑證），上面記錄了您的公鑰

**結果**：
- ✅ Fabric CA 簽發 X.509 憑證
- ✅ 憑證包含 Transit 的公鑰
- ✅ 憑證證明：`user-patient_001` 這個身份是合法的

---

## 📋 步驟 6️⃣：儲存憑證到 Vault

```go
store.WriteUserMaterial(ctx, req.UserId, csrPEM, keyPEM, certPem)
// → 寫入 Vault KV
```

**儲存位置**：
```
/v1/kv/data/users/patient_001
{
  "data": {
    "csr": "CSR 內容",
    "key": "",  // 空字串（不使用傳統私鑰）
    "cert": "Fabric CA 簽發的憑證"  ← 重要！
  }
}
```

---

## 🎯 關係圖解

```
┌──────────────────────────────────────────────────────────────┐
│                    用戶註冊流程                               │
└──────────────────────────────────────────────────────────────┘

步驟 1: Transit 產生簽章金鑰
┌──────────────┐
│ Vault        │
│ Transit      │ → 建立 user-patient_001 (ecdsa-p256)
│              │ → 有公鑰和私鑰（私鑰不離開）
└──────┬───────┘
       │ 取得公鑰
       ▼
步驟 2-3: 用 Transit 公鑰產生 CSR
┌──────────────┐
│ Transit      │ → 用 Transit 私鑰簽署 CSR
│ Signer       │ → CSR 包含 Transit 公鑰
└──────┬───────┘
       │ 提交 CSR
       ▼
步驟 4-5: Fabric CA 簽發憑證
┌──────────────┐
│ Fabric CA    │ → 驗證 CSR 簽章
│              │ → 簽發憑證（包含 Transit 公鑰）
└──────┬───────┘
       │ 返回憑證
       ▼
步驟 6: 儲存憑證
┌──────────────┐
│ Vault KV     │ → 儲存 CA 簽發的憑證
│              │ → cert: Fabric CA 憑證
└──────────────┘

┌──────────────────────────────────────────────────────────────┐
│                    後續使用                                  │
└──────────────────────────────────────────────────────────────┘

Fabric 交易簽章：
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│ 取得憑證     │ →   │ 取得 Transit │ →   │ 簽章交易      │
│ (從 Vault KV)│     │ 簽章金鑰引用 │     │ (TransitSign) │
└──────────────┘     └──────────────┘     └──────────────┘
  ↑                    ↑
  │                    │
  └─ Fabric CA 簽發    └─ Transit 提供簽章能力
```

---

## 🔍 關鍵理解點

### ✅ Fabric CA 的角色

1. **註冊身份**：
   - 登記用戶到系統
   - 給用戶 ID 和密碼

2. **簽發憑證**：
   - 收到 CSR（包含 Transit 公鑰）
   - 驗證 CSR 是合法的（用 Transit 公鑰驗證）
   - 簽發憑證證明這個身份

3. **不知道私鑰**：
   - ❌ Fabric CA **不知道** Transit 私鑰
   - ✅ 只知道 Transit 公鑰（在 CSR 和憑證中）

---

### ✅ Transit 的角色

1. **管理私鑰**：
   - 產生簽章金鑰（公鑰+私鑰對）
   - **私鑰永遠不離開** Vault

2. **提供簽章能力**：
   - 簽署 CSR
   - 簽署 Fabric 交易
   - 所有簽章操作都在 Vault 內完成

3. **不直接與 Fabric CA 溝通**：
   - Transit 和 Fabric CA **不會直接對話**
   - 它們通過「Go Server」協作

---

## 🎓 類比說明

### 🏛️ 真實世界的類比

**Fabric CA = 戶政事務所**
- 負責發身份證
- 需要看您的「申請書」（CSR）
- 驗證申請書是真的（用公鑰驗證簽章）
- 發給您身份證（憑證）

**Transit = 銀行的保險箱**
- 存放您的「印章」（私鑰）
- 您用印章簽署申請書
- 印章永遠不離開銀行
- 銀行幫您保管印章

**協作流程**：
1. 在銀行開保險箱，拿到印章（Transit 產生金鑰）
2. 用印章簽署申請書（CSR）
3. 把申請書送到戶政事務所（Fabric CA）
4. 戶政事務所驗證申請書（用公鑰驗證簽章）
5. 戶政事務所發身份證（憑證）
6. 身份證存回家裡（Vault KV）

**後續使用**：
- 需要簽署文件時，去銀行用印章（Transit 簽章）
- 向別人證明身份時，出示身份證（Fabric CA 憑證）

---

## 📊 它們的關係總結

| 特性 | Fabric CA | Transit |
|------|----------|---------|
| **職責** | 發放憑證 | 管理金鑰和簽章 |
| **知道什麼** | Transit 公鑰（在 CSR 和憑證中） | 私鑰（永遠不離開） |
| **不知道什麼** | Transit 私鑰 | Fabric CA 不知道私鑰 |
| **直接溝通** | ❌ 不直接溝通 | ❌ 不直接溝通 |
| **協作方式** | 通過 Go Server 協作 | 通過 Go Server 協作 |
| **產生的東西** | X.509 憑證 | 簽章金鑰、加密金鑰 |
| **儲存位置** | 憑證存 Vault KV | 金鑰存 Vault Transit |

---

## 🔗 協作流程示意圖

```
┌─────────────────────────────────────────────────────────────┐
│                    Go Server (協調者)                        │
│  ┌────────────────────────────────────────────────────┐     │
│  │  1. 產生 Transit 簽章金鑰                          │     │
│  │     → 呼叫 Vault Transit                          │     │
│  │                                                     │     │
│  │  2. 取得 Transit 公鑰                              │     │
│  │     → 呼叫 Vault Transit                          │     │
│  │                                                     │     │
│  │  3. 產生 CSR（用 Transit 簽章）                      │     │
│  │     → 呼叫 TransitSigner                          │     │
│  │                                                     │     │
│  │  4. 註冊到 Fabric CA                               │     │
│  │     → 呼叫 Fabric CA API                          │     │
│  │                                                     │     │
│  │  5. Enroll 取得憑證（提交 CSR）                      │     │
│  │     → 呼叫 Fabric CA API                          │     │
│  │                                                     │     │
│  │  6. 儲存憑證到 Vault KV                            │     │
│  │     → 呼叫 Vault KV API                           │     │
│  └────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
       │                                    │
       ▼                                    ▼
┌─────────────┐                    ┌──────────────┐
│ Fabric CA   │                    │ Vault        │
│             │                    │              │
│ - 註冊身份  │                    │ - Transit    │
│ - 簽發憑證  │                    │   (簽章金鑰) │
│             │                    │ - KV         │
│ 不直接溝通！ │                    │   (存憑證)   │
└─────────────┘                    └──────────────┘
```

---

## ❓ 常見問題

### Q1: Fabric CA 和 Transit 會互相知道對方嗎？

**A**: 
- ❌ **不會直接溝通**
- ✅ 它們都**不知道對方的存在**
- ✅ 所有協作都通過 **Go Server** 完成

### Q2: 為什麼要用 Transit，不用傳統私鑰檔案？

**A**: 
- ✅ **更安全**：私鑰永遠不離開 Vault
- ✅ **更易管理**：不需要處理 PEM 檔案
- ✅ **支援審計**：Vault 記錄所有簽章操作
- ✅ **支援輪換**：可以自動輪換金鑰

### Q3: Fabric CA 會知道 Transit 私鑰嗎？

**A**: 
- ❌ **絕對不會！**
- ✅ Fabric CA **只看到** Transit 公鑰（在 CSR 和憑證中）
- ✅ 私鑰永遠在 Vault Transit 中，**從不離開**

### Q4: 可以不用 Transit 嗎？

**A**: 
- ⚠️ **理論上可以**，但需要：
  - 在本地產生 ECDSA 私鑰檔案
  - 用私鑰檔案簽署 CSR
  - 管理私鑰檔案的安全性
- ✅ **使用 Transit 更好**：
  - 統一金鑰管理
  - 更高的安全性
  - 更容易操作

### Q5: Transit 簽章金鑰和 wrap 金鑰是同一個嗎？

**A**: 
- ❌ **不是！** 它們是完全不同的金鑰：
  - **簽章金鑰**：`user-{id}`（ecdsa-p256，用於簽章）
  - **Wrap 金鑰**：`user-{id}-wrap`（aes256-gcm96，用於加密 dataKey）

---

## 📝 總結

### ✅ Fabric CA 和 Transit 的關係

1. **它們是協作關係**：
   - ✅ Transit 提供私鑰功能（簽章 CSR）
   - ✅ Fabric CA 簽發憑證（證明身份）
   - ✅ 它們通過 Go Server 協作

2. **它們不知道對方**：
   - ❌ Fabric CA 不知道 Transit 的存在
   - ❌ Transit 不知道 Fabric CA 的存在
   - ✅ 它們只關心自己的工作

3. **分工明確**：
   - **Fabric CA**：負責「身份認證和憑證管理」
   - **Transit**：負責「金鑰管理和簽章」

4. **共同目標**：
   - ✅ 讓用戶能夠安全地在 Fabric 區塊鏈上操作
   - ✅ 確保私鑰安全，身份可驗證

### 🎯 一句話總結

**Fabric CA 發身份證，Transit 是印章保險庫；**
**身份證證明你是誰，印章用來簽署文件；**
**它們互相配合，但不直接打交道。**

---

這樣的設計讓系統更加**安全**（私鑰不離開 Vault）和**靈活**（可以獨立管理身份和金鑰）！



