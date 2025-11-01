# 🔐 Medledger 加密架構完整說明

## 📖 目錄
1. [核心概念：保險箱比喻](#核心概念保險箱比喻)
2. [系統架構總覽](#系統架構總覽)
3. [加密流程詳解](#加密流程詳解)
4. [解密流程詳解](#解密流程詳解)
5. [動態授權機制](#動態授權機制)
6. [技術細節](#技術細節)

---

## 🏦 核心概念：保險箱比喻

想像您有一份**健檢報告**需要安全儲存：

### 傳統方式（不安全）
```
報告 → 放在一個大保險箱 → 很多人有這把鑰匙 ❌
問題：只要有鑰匙的人都能看
```

### 我們的混合加密方式（安全）
```
報告 → 用一次性密碼鎖上（AES-256-GCM）
       ↓
密碼 → 分別鎖在多個獨立保險箱中（Vault Transit）
       ↓
       ┌─ 診所保險箱（clinic-key）
       ├─ 病患保險箱（patient-key）  
       └─ 保險公司保險箱（insurer-key）
       
✅ 優點：每個人只能用自己的鑰匙打開自己的保險箱，拿到密碼後才能解開報告
```

---

## 🏗️ 系統架構總覽

```
┌─────────────────────────────────────────────────────────────┐
│                    前端應用程式                              │
│              (上傳報告 / 查看報告)                           │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                   Go Server (gRPC/HTTP)                     │
│  ┌────────────────────────────────────────────────────┐     │
│  │         service/report_service.go                 │     │
│  │  - HandleUploadReport()                            │     │
│  │  - HandleReadMyReport()                            │     │
│  └────────────────┬───────────────────────────────────┘     │
│                   │ 調用                                        │
│                   ▼                                            │
│  ┌────────────────────────────────────────────────────┐     │
│  │      secure/wrap/transit.go (TransitWrapper)       │     │
│  │  - EncryptReportTransitClinicOnly()                 │     │
│  │  - DecryptReportTransit()                           │     │
│  │  - AddRecipientTransitFrom()                        │     │
│  └────────────────┬───────────────────────────────────┘     │
│                   │ 使用                                        │
│                   ▼                                            │
│  ┌────────────────────────────────────────────────────┐     │
│  │    vaultstore/vaultstore.go (VaultStore)           │     │
│  │  - TransitEncrypt()                                 │     │
│  │  - TransitDecrypt()                                 │     │
│  └────────────────┬───────────────────────────────────┘     │
└───────────────────┼───────────────────────────────────────────┘
                    │ HTTP API
                    ▼
┌─────────────────────────────────────────────────────────────┐
│              HashiCorp Vault Server                         │
│  ┌────────────────────────────────────────────────────┐     │
│  │         Transit Engine (金鑰管理服務)               │     │
│  │  - 儲存和管理 wrap keys (aes256-gcm96)             │     │
│  │  - 提供加密/解密 API                                │     │
│  │  - 金鑰：clinic-<id>-wrap, user-<id>-wrap 等      │     │
│  └────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
                    │
                    ▼ 儲存加密後的報告
┌─────────────────────────────────────────────────────────────┐
│           Hyperledger Fabric 區塊鏈                          │
│  ┌────────────────────────────────────────────────────┐     │
│  │            Chaincode (Smart Contract)              │     │
│  │  - 儲存加密後的 Envelope JSON                      │     │
│  │  - 管理授權邏輯                                    │     │
│  └────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔒 加密流程詳解

### 場景：診所上傳健檢報告

讓我們逐步追蹤一份報告從明文到加密儲存的過程：

#### 步驟 1️⃣：前端發送上傳請求
```javascript
// 前端送出
{
  report_id: "report_12345",
  user_id: "patient_001", 
  test_results_json: '{"血糖": 95, "膽固醇": 180, ...}'
}
```

#### 步驟 2️⃣：進入 EncryptReportTransitClinicOnly()

讓我們看看實際代碼做了什麼：

```40:78:secure/wrap/transit.go
func (w *TransitWrapper) EncryptReportTransitClinicOnly(ctx context.Context, plaintext []byte, clinicLabel string, clinicKey string) ([]byte, error) {
    log.Printf("[Transit] EncryptReportTransitClinicOnly: clinicLabel=%s clinicKey=%s plainLen=%d", clinicLabel, clinicKey, len(plaintext))
    // 產生隨機 32 位元組的 dataKey
    dk := make([]byte, 32)              
    if _, err := rand.Read(dk); err != nil { return nil, err }
    // 產生隨機 12 位元組的 nonce（GCM 標準長度）
    nonce := make([]byte, 12)
    if _, err := rand.Read(nonce); err != nil { return nil, err }
    // 建立 AES cipher 並使用 GCM 模式加密
    block, err := aes.NewCipher(dk)
    if err != nil { return nil, err }
    gcm, err := cipher.NewGCM(block)
    if err != nil { return nil, err }
    ct := gcm.Seal(nil, nonce, plaintext, nil)
    log.Printf("[Transit] AES-GCM done (clinic-only): ctLen=%d nonceLen=%d", len(ct), len(nonce))

    // 診所端使用獨立 wrap key（若傳入的是簽章 key 名稱，則補 -wrap）
    wrapClinicKey := clinicKey
    if !strings.HasSuffix(wrapClinicKey, "-wrap") { wrapClinicKey = wrapClinicKey + "-wrap" }
    if err := w.store.EnsureTransitKeyOfType(ctx, wrapClinicKey, "aes256-gcm96"); err != nil { return nil, err }

    // 以 Transit 包 dataKey（僅診所）
    ctClinic, err := w.store.TransitEncrypt(ctx, wrapClinicKey, dk)
    if err != nil { log.Printf("[Transit] TransitEncrypt(clinic %s) error: %v", clinicKey, err); return nil, err }

    env := map[string]any{
        "ciphertext": b64.StdEncoding.EncodeToString(ct),
        "nonce":      b64.StdEncoding.EncodeToString(nonce),
        "wrappedKeys": map[string]any{
            clinicLabel: map[string]string{"type":"transit", "ct": ctClinic},
        },
        "enc":  "AES-256-GCM",
        "kdf":  "n/a",
        "curve":"n/a",
    }
    envCipherB64 := env["ciphertext"].(string)
    log.Printf("[Transit] Envelope (clinic-only) ready: cipherB64Len=%d", len(envCipherB64))
    return json.Marshal(env)
}
```

**詳細拆解每個步驟：**

#### 📝 步驟 2.1：產生一次性金鑰（Data Key）
```go
dk := make([]byte, 32)  // 256 位元的隨機金鑰
rand.Read(dk)           // 從系統安全隨機數產生器讀取
```
**比喻**：就像為這份報告產生一個**唯一的密碼**，每份報告的密碼都不同。

#### 📝 步驟 2.2：產生 Nonce（一次性數字）
```go
nonce := make([]byte, 12)  // 12 字節的隨機數
rand.Read(nonce)
```
**比喻**：這是防止重放攻擊的**鹽值**，確保即使相同內容加密結果也不同。

#### 📝 步驟 2.3：用 AES-256-GCM 加密報告內容
```go
block, _ := aes.NewCipher(dk)      // 用 dataKey 建立 AES cipher
gcm, _ := cipher.NewGCM(block)      // 轉成 GCM 模式（可認證加密）
ct := gcm.Seal(nil, nonce, plaintext, nil)  // 加密！
```
**過程**：
```
明文報告 → [AES-256-GCM + dataKey + nonce] → 加密報告 (ciphertext)
```
**比喻**：用一次性密碼鎖把報告鎖進保險箱。

#### 📝 步驟 2.4：把密碼（dataKey）鎖進診所的保險箱
```go
wrapClinicKey := "clinic-" + clinicID + "-wrap"  // 例如：clinic-001-wrap
w.store.TransitEncrypt(ctx, wrapClinicKey, dk)   // 請 Vault 用診所鑰匙加密 dataKey
```
**過程**：
```
dataKey → [Vault Transit + clinic-001-wrap] → 加密的 dataKey (vault:v1:...)
```
**比喻**：把保險箱密碼放進診所專屬的保險箱，只有診所有鑰匙。

#### 📝 步驟 2.5：打包成 Envelope（信封）
最終生成一個 JSON 結構：
```json
{
  "ciphertext": "base64(加密後的報告)",
  "nonce": "base64(12字節隨機數)",
  "wrappedKeys": {
    "clinic": {
      "type": "transit",
      "ct": "vault:v1:加密後的dataKey"
    }
  },
  "enc": "AES-256-GCM",
  "kdf": "n/a",
  "curve": "n/a"
}
```

**比喻**：就像一個安全信封：
- 📄 **ciphertext**：鎖住的報告
- 🔑 **wrappedKeys**：可以取得密碼的鑰匙盒（目前只有診所的）
- 📋 **nonce**：解鎖需要的參數

---

#### 步驟 3️⃣：為病患也準備一份鑰匙

從 `service/report_service.go` 可以看到：

```65:73:service/report_service.go
    // Transit：改為 clinic-only 加密，平台不再持有解包能力
    tw, err := sw.NewTransitWrapperFromEnv(); if err != nil { return nil, status.Error(codes.Internal, "Vault 初始化失敗") }
    // label 固定為 "clinic"，實際身分由 baseKey（clinic-<clinicId>-wrap）決定
    envJSON, err := tw.EncryptReportTransitClinicOnly(ctx, []byte(req.TestResultsJson), "clinic", "clinic-"+userID)
    if err != nil { return nil, status.Error(codes.Internal, "資料加密失敗") }
    // 上傳階段預先授權病患：用診所 unwrap 來源追加 patient wrapped key（病患可離線）
    // unwrapLabel 固定 "clinic"；newLabel 固定 "patient"
    envJSON, err = tw.AddRecipientTransitFrom(ctx, envJSON, "clinic", "clinic-"+userID, "patient", "user-"+req.UserId)
```

**說明**：
1. 先用診所 key 加密（步驟 2）
2. 再用 `AddRecipientTransitFrom` 為病患添加一份 wrapped key

讓我們看 `AddRecipientTransitFrom` 做了什麼：

```89:122:secure/wrap/transit.go
func (w *TransitWrapper) AddRecipientTransitFrom(ctx context.Context, envJSON []byte, unwrapLabel string, unwrapBaseKey string, newLabel string, newBaseKey string) ([]byte, error) {
    log.Printf("[Transit] AddRecipientTransitFrom: unwrapLabel=%s unwrapBaseKey=%s newLabel=%s newBaseKey=%s envSize=%d", unwrapLabel, unwrapBaseKey, newLabel, newBaseKey, len(envJSON))
    var env struct {
        Ciphertext  string                        `json:"ciphertext"`
        Nonce       string                        `json:"nonce"`
        WrappedKeys map[string]map[string]string `json:"wrappedKeys"`
        Enc         string                        `json:"enc"`
        KDF         string                        `json:"kdf"`
        Curve       string                        `json:"curve"`
    }
    if err := json.Unmarshal(envJSON, &env); err != nil { return nil, err }
    src, ok := env.WrappedKeys[unwrapLabel]
    if !ok { return nil, &wrapError{"missing wrapped key label: " + unwrapLabel} }
    ctSrc := src["ct"]

    // 使用來源金鑰解包 dataKey（自動補 -wrap 後綴）
    unwrapKey := unwrapBaseKey
    if !strings.HasSuffix(unwrapKey, "-wrap") { unwrapKey = unwrapKey + "-wrap" }
    if err := w.store.EnsureTransitKeyOfType(ctx, unwrapKey, "aes256-gcm96"); err != nil { return nil, err }
    dk, err := w.store.TransitDecrypt(ctx, unwrapKey, ctSrc)
    if err != nil { return nil, err }
    log.Printf("[Transit] AddRecipientTransitFrom: decrypted dataKey len=%d via %s", len(dk), unwrapKey)

    // 以 newBaseKey 重包
    wrapKey := newBaseKey
    if !strings.HasSuffix(wrapKey, "-wrap") { wrapKey = wrapKey + "-wrap" }
    if err := w.store.EnsureTransitKeyOfType(ctx, wrapKey, "aes256-gcm96"); err != nil { return nil, err }
    ctNew, err := w.store.TransitEncrypt(ctx, wrapKey, dk)
    if err != nil { return nil, err }
    if env.WrappedKeys == nil { env.WrappedKeys = map[string]map[string]string{} }
    env.WrappedKeys[newLabel] = map[string]string{"type":"transit", "ct": ctNew}
    log.Printf("[Transit] AddRecipientTransitFrom: added label=%s via wrapKey=%s ctLen=%d", newLabel, wrapKey, len(ctNew))
    return json.Marshal(env)
}
```

**過程**：
1. **解包**：用診所的 key (`clinic-001-wrap`) 解出 dataKey
2. **重包**：用病患的 key (`user-patient_001-wrap`) 重新加密 dataKey
3. **更新**：在 `wrappedKeys` 中添加病患的條目

**最終 Envelope 變成**：
```json
{
  "ciphertext": "base64(加密報告)",
  "nonce": "base64(nonce)",
  "wrappedKeys": {
    "clinic": {
      "type": "transit",
      "ct": "vault:v1:診所加密的dataKey"
    },
    "patient": {                    // ← 新增！
      "type": "transit", 
      "ct": "vault:v1:病患加密的dataKey"
    }
  },
  "enc": "AES-256-GCM",
  ...
}
```

#### 步驟 4️⃣：儲存到區塊鏈
```go
contract.SubmitTransaction("UploadReport", reportID, patientHash, string(envJSON))
```

---

## 🔓 解密流程詳解

### 場景：病患查看自己的報告

#### 步驟 1️⃣：從區塊鏈讀取 Envelope
```go
result, err := contract.EvaluateTransaction("ReadMyReport", reportID)
// result 就是之前儲存的 Envelope JSON
```

#### 步驟 2️⃣：解密過程

看 `DecryptReportTransit` 函數：

```124:158:secure/wrap/transit.go
func (w *TransitWrapper) DecryptReportTransit(ctx context.Context, envJSON []byte, label string, baseKey string) ([]byte, error) {
    log.Printf("[Transit] DecryptReportTransit: label=%s baseKey=%s envSize=%d", label, baseKey, len(envJSON))
    var env struct {
        Ciphertext string                        `json:"ciphertext"`
        Nonce      string                        `json:"nonce"`
        WrappedKeys map[string]map[string]string `json:"wrappedKeys"`
    }
    if err := json.Unmarshal(envJSON, &env); err != nil { log.Printf("[Transit] DecryptReportTransit: bad envelope json: %v", err); return nil, err }
    rec, ok := env.WrappedKeys[label]
    if !ok { log.Printf("[Transit] DecryptReportTransit: missing wrapped label=%s", label); return nil, &wrapError{"missing wrapped key label: " + label} }
    ctWrap := rec["ct"]
    // 以 AES wrap key 解出 dataKey
    wrapKey := baseKey
    if !strings.HasSuffix(wrapKey, "-wrap") { wrapKey = wrapKey + "-wrap" }
    log.Printf("[Transit] DecryptReportTransit: using wrapKey=%s", wrapKey)
    if err := w.store.EnsureTransitKeyOfType(ctx, wrapKey, "aes256-gcm96"); err != nil { log.Printf("[Transit] Ensure wrap key failed: %v", err); return nil, err }
    dataKey, err := w.store.TransitDecrypt(ctx, wrapKey, ctWrap)
    if err != nil { log.Printf("[Transit] TransitDecrypt(wrap) error: %v", err); return nil, err }
    log.Printf("[Transit] DecryptReportTransit: dataKeyLen=%d", len(dataKey))
    // 以 dataKey 做 AES-GCM 解密主密文
    ctBytes, err := b64.StdEncoding.DecodeString(env.Ciphertext)
    if err != nil { log.Printf("[Transit] Bad ciphertext b64: %v", err); return nil, err }
    nonce, err := b64.StdEncoding.DecodeString(env.Nonce)
    if err != nil { log.Printf("[Transit] Bad nonce b64: %v", err); return nil, err }
    block, err := aes.NewCipher(dataKey)
    if err != nil { log.Printf("[Transit] AES cipher init error: %v", err); return nil, err }
    gcm, err := cipher.NewGCM(block)
    if err != nil { log.Printf("[Transit] GCM init error: %v", err); return nil, err }
    pt, err := gcm.Open(nil, nonce, ctBytes, nil)
    if err != nil { log.Printf("[Transit] GCM open error: %v", err); return nil, err }
    log.Printf("[Transit] DecryptReportTransit: plainLen=%d", len(pt))
    return pt, nil
}
```

**詳細步驟：**

##### 🔑 步驟 2.1：找到病患的 wrapped key
```go
rec := env.WrappedKeys["patient"]  // 取得病患的加密 dataKey
ctWrap := rec["ct"]                 // "vault:v1:..."
```

##### 🔓 步驟 2.2：用病患的鑰匙解出 dataKey
```go
wrapKey := "user-patient_001-wrap"
dataKey := w.store.TransitDecrypt(ctx, wrapKey, ctWrap)
```
**過程**：
```
vault:v1:加密的dataKey → [Vault Transit + user-xxx-wrap] → dataKey (32 bytes)
```
**比喻**：用病患的鑰匙打開病患保險箱，拿到密碼。

##### 📖 步驟 2.3：用 dataKey 解開報告
```go
gcm.Open(nil, nonce, ciphertext, nil)  // 解密！
```
**過程**：
```
加密報告 + dataKey + nonce → [AES-256-GCM 解密] → 明文報告
```
**比喻**：用密碼打開主保險箱，取出報告。

---

## 🔄 動態授權機制

### 場景：保險公司申請查看報告

#### 情況 1：病患批准後

系統會再次呼叫 `AddRecipientTransitFrom`：

```go
// 用病患的 key 解出 dataKey，再用保險公司的 key 重包
envJSON = tw.AddRecipientTransitFrom(
    ctx, envJSON,
    "patient",           // 從病患的 wrapped key
    "user-patient_001",  // 病患的 baseKey
    "insurer",           // 新增保險公司標籤
    "insurer-001"        // 保險公司的 baseKey
)
```

**最終 Envelope**：
```json
{
  "ciphertext": "base64(加密報告)",
  "wrappedKeys": {
    "clinic": {"type": "transit", "ct": "vault:v1:..."},
    "patient": {"type": "transit", "ct": "vault:v1:..."},
    "insurer": {"type": "transit", "ct": "vault:v1:..."}  // ← 新增！
  },
  ...
}
```

**重要特性**：
- ✅ **不需要重新加密報告本身**：只有 dataKey 被重包裝
- ✅ **零知識架構**：平台伺服器永遠看不到明文
- ✅ **細粒度控制**：每個角色只能用自己的 key 解包

---

## 🔧 技術細節

### 1. 為什麼用混合加密？

**問題**：為什麼不直接用 Vault Transit 加密整個報告？

**答案**：
- ❌ **性能問題**：大型報告加密/解密慢
- ❌ **靈活性差**：無法動態添加收件者（需重新加密整個文件）
- ✅ **混合方案**：小 dataKey 用 Transit，大內容用 AES（快速）

### 2. 為什麼是 Clinic-Only 初始加密？

**設計決策**：`EncryptReportTransitClinicOnly` 只為診所包裝 dataKey

**原因**：
- ✅ **零信任**：平台無法解包查看內容
- ✅ **患者隱私**：只有診所和病患能解密
- ✅ **動態授權**：後續可添加其他收件者

### 3. 金鑰命名規則

```go
// 診所
"clinic-001-wrap"  // 診所 ID 001 的包裝金鑰

// 病患  
"user-patient_001-wrap"  // 病患 ID 的包裝金鑰

// 保險公司
"insurer-001-wrap"  // 保險公司 ID 的包裝金鑰
```

規則：
- `{role}-{id}-wrap` 格式
- 自動補 `-wrap` 後綴（如果沒有）

### 4. Vault Transit 如何工作？

**Vault Store 封裝**：

```145:165:vaultstore/vaultstore.go
// TransitEncrypt 將 dataKey 以 Transit 指定 keyName 加密，回傳 ciphertext（vault:vX:...）
func (s *Store) TransitEncrypt(ctx context.Context, keyName string, data []byte) (string, error) {
    if keyName == "" { return "", fmt.Errorf("empty transit keyName") }
    // plaintext 需為 base64；Vault SDK 會轉送 body，這裡直接傳字串
    b64 := map[string]interface{}{
        "plaintext":  encodeB64(data),
    }
    path := fmt.Sprintf("%s/encrypt/%s", s.transitMount, keyName)
    if _, ok := ctx.Deadline(); !ok {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
        defer cancel()
    }
    log.Printf("[Vault] TransitEncrypt path=/v1/%s dataLen=%d", path, len(data))
    sec, err := s.client.Logical().WriteWithContext(ctx, path, b64)
    if err != nil { return "", err }
    if sec == nil || sec.Data == nil { return "", fmt.Errorf("empty transit response") }
    ct, _ := sec.Data["ciphertext"].(string)
    if ct == "" { return "", fmt.Errorf("missing ciphertext") }
    return ct, nil
}
```

**流程**：
1. 調用 Vault API：`POST /v1/transit/encrypt/clinic-001-wrap`
2. 傳入 base64 編碼的 dataKey
3. Vault 用該 key 加密並返回：`vault:v1:加密結果`

**優點**：
- 🔐 **金鑰永遠不離開 Vault**：只有 Vault 知道 wrap key
- 🔒 **審計日誌**：Vault 記錄每次加密/解密
- 🔄 **金鑰輪換**：Vault 可自動輪換，不影響已加密數據

### 5. AES-256-GCM 為什麼好？

**AES-256**：256 位金鑰強度（32 字節）

**GCM 模式**：
- ✅ **認證加密（AEAD）**：同時提供加密和完整性驗證
- ✅ **隨機 IV/nonce**：每次加密都不同，防止重放攻擊
- ✅ **高效能**：硬體加速支援，適合大量數據

---

## 📊 完整流程圖

```
┌─────────────┐
│  明文報告   │
└──────┬──────┘
       │
       ▼
┌─────────────────────────┐
│ 1. 產生隨機 dataKey (32B)│
│ 2. 產生隨機 nonce (12B) │
└──────┬──────────────────┘
       │
       ▼
┌─────────────────────────┐      ┌──────────────────┐
│ 3. AES-256-GCM 加密     │      │ 4. Transit 包裝   │
│    plaintext + dataKey  │      │    dataKey        │
│    → ciphertext         │      │    → wrapped key  │
└──────┬──────────────────┘      └──────┬───────────┘
       │                                │
       └──────────┬─────────────────────┘
                  ▼
         ┌─────────────────┐
         │ 5. 打包 Envelope │
         │   (JSON)         │
         └────────┬─────────┘
                  │
                  ▼
    ┌─────────────────────────┐
    │ 6. 添加病患 wrapped key  │
    │    (AddRecipient)        │
    └────────┬────────────────┘
             │
             ▼
    ┌─────────────────────────┐
    │ 7. 提交到區塊鏈         │
    └─────────────────────────┘
```

**解密流程（反向）**：
```
┌─────────────────────────┐
│ 1. 從區塊鏈讀取 Envelope│
└──────┬──────────────────┘
       │
       ▼
┌─────────────────────────┐
│ 2. Transit 解包          │
│    wrapped key → dataKey│
└──────┬──────────────────┘
       │
       ▼
┌─────────────────────────┐
│ 3. AES-256-GCM 解密     │
│    ciphertext + dataKey │
│    → plaintext          │
└──────┬──────────────────┘
       │
       ▼
┌─────────────┐
│  明文報告   │
└─────────────┘
```

---

## ✅ 安全特性總結

1. **零知識架構**：平台無法查看內容
2. **角色分離**：每個角色有獨立金鑰
3. **動態授權**：可後續添加收件者
4. **標準加密**：AES-256-GCM + Vault Transit
5. **不可抵賴**：區塊鏈記錄所有操作
6. **金鑰管理**：Vault 專業金鑰管理

---

## 🎓 總結

這套系統的核心思想是：
- **混合加密**：小 dataKey 用 Transit，大內容用 AES
- **分層授權**：每個人用自己的鑰匙開自己的保險箱
- **動態擴展**：授權時再添加新的 wrapped key，不需重新加密
- **零信任**：即使是平台也看不到內容

就像銀行保險庫系統：
- 📦 保險箱本身（報告）用一次性密碼鎖（AES）
- 🔑 密碼分別鎖在不同房間（Transit）
- 👤 每個人只能進自己的房間拿密碼
- ✅ 即使銀行員工也看不到保險箱內容


