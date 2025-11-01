package wrap

import (
    "context"
    "crypto/aes"
    "crypto/cipher"
    "crypto/sha256"
    "crypto/rand"
    b64 "encoding/base64"
    "encoding/json"
    "fmt"
    "log"
    "sort"
    "strings"

    vs "go_server/vaultstore"
)

// TransitWrapper 封裝以 Vault Transit 進行資料金鑰包裝/解包與 AES-GCM 主密文加解
type TransitWrapper struct {
    store         *vs.Store
}

// NewTransitWrapperFromEnv 以環境變數建立 Wrapper
// 需要：VAULT_ADDR, VAULT_TOKEN；可選：VAULT_MOUNT（kv）、VAULT_TRANSIT_MOUNT（預設 transit）
func NewTransitWrapperFromEnv() (*TransitWrapper, error) {
    st, err := vs.NewFromEnv()
    if err != nil { return nil, err }
    log.Printf("[Transit] NewWrapper (clinic/patient only)")
    return &TransitWrapper{store: st}, nil
}

// EncryptReportTransit：
// 1) 產生隨機 32B dataKey 與 12B nonce
// 2) AES-GCM(dataKey) 加密 plaintext
// 3) 以 Transit 為 clinicKey 與 platformKey 分別包 dataKey
// 4) 回傳 Envelope JSON：{"ciphertext": base64(ct), "nonce": base64(nonce), "wrappedKeys": {label:{"type":"transit","ct":...}}, ...}
// EncryptReportTransit 已移除：請使用 EncryptReportTransitClinicOnly

// EncryptReportTransitClinicOnly：只以診所 wrap dataKey，不包含平台 wrapped key
// 使用情境：避免平台具備解包能力，後續追加收件者改由診所/病患的 wrap key 進行 rewrap
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

    // 以診所簽章金鑰簽署 Envelope（簽署固定序列化的摘要，不依賴 map 順序）
    signerKey := clinicKey
    signerKey = strings.TrimSuffix(signerKey, "-wrap")
    derB64, sigErr := w.signEnvelope(ctx, env, signerKey)
    if sigErr != nil { return nil, sigErr }
    env["sig"] = map[string]string{
        "alg": "ecdsa-p256-sha256-asn1",
        "by":  fmt.Sprintf("transit://%s", signerKey),
        "der": derB64,
    }
    envCipherB64 := env["ciphertext"].(string)
    log.Printf("[Transit] Envelope (clinic-only) ready: cipherB64Len=%d", len(envCipherB64))
    return json.Marshal(env)
}



// AddRecipientTransitFrom：指定 unwrap 來源（label/baseKey），於 Vault 內部完成「解包 -> 以新 key 重包」
// 參數：
//   - unwrapLabel：既有 wrapped key 的標籤（例如診所 label 或 patient）
//   - unwrapBaseKey：對應 unwrapLabel 的 Transit base key 名稱（例如 "clinic-<hash>" 或 "user-<hash>"），函式內會自動加上 -wrap
//   - newLabel/newBaseKey：要新增的收件者標籤與 Transit base key 名稱（同樣自動加 -wrap）
// 注意：所有 baseKey 統一使用 hash 值（與註冊邏輯一致）
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

    // 重新簽章（使用 unwrapBaseKey 的簽章金鑰）
    signerKey := strings.TrimSuffix(unwrapBaseKey, "-wrap")
    // 將可序列化物件重建以包含 sig
    out := map[string]any{
        "ciphertext": env.Ciphertext,
        "nonce":      env.Nonce,
        "wrappedKeys": env.WrappedKeys,
        "enc":        env.Enc,
        "kdf":        env.KDF,
        "curve":      env.Curve,
    }
    derB64, sigErr := w.signEnvelope(ctx, out, signerKey)
    if sigErr != nil { return nil, sigErr }
    out["sig"] = map[string]string{
        "alg": "ecdsa-p256-sha256-asn1",
        "by":  fmt.Sprintf("transit://%s", signerKey),
        "der": derB64,
    }
    return json.Marshal(out)
}

// DecryptReportTransit：使用指定角色/標籤與對應 wrap key 解出明文
// label 例："patient"、"insurer"、"platform"；baseKey 例：user-<hash>、insurer-<hash>、clinic-<hash>
// 注意：baseKey 統一使用 hash 值（與註冊邏輯一致），函式內會自動加上 -wrap 後綴
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

type wrapError struct{ msg string }
func (e *wrapError) Error() string { return e.msg }


// signEnvelope 以排序後的 wrappedKeys 建立穩定摘要並使用 Transit 簽章
func (w *TransitWrapper) signEnvelope(ctx context.Context, env map[string]any, signerKey string) (string, error) {
    if signerKey == "" { return "", &wrapError{"empty signerKey"} }
    // 準備待簽內容
    cipherB64, _ := env["ciphertext"].(string)
    nonceB64, _ := env["nonce"].(string)
    enc, _ := env["enc"].(string)
    kdf, _ := env["kdf"].(string)
    curve, _ := env["curve"].(string)

    // 取出 wrappedKeys 並以標籤排序，僅納入 type+ct（避免 map 順序問題）
    wks, _ := env["wrappedKeys"].(map[string]any)
    labels := make([]string, 0, len(wks))
    for k := range wks { labels = append(labels, k) }
    sort.Strings(labels)
    // 建立 canonical 字串
    b := strings.Builder{}
    b.WriteString("ciphertext="); b.WriteString(cipherB64); b.WriteString("\n")
    b.WriteString("nonce="); b.WriteString(nonceB64); b.WriteString("\n")
    b.WriteString("enc="); b.WriteString(enc); b.WriteString("\n")
    b.WriteString("kdf="); b.WriteString(kdf); b.WriteString("\n")
    b.WriteString("curve="); b.WriteString(curve); b.WriteString("\n")
    for _, label := range labels {
        recAny := wks[label]
        rec, _ := recAny.(map[string]string)
        if rec == nil { // 兼容 map[string]any
            if tmp, ok := recAny.(map[string]any); ok {
                rec = map[string]string{
                    "type": fmt.Sprint(tmp["type"]),
                    "ct":   fmt.Sprint(tmp["ct"]),
                }
            }
        }
        b.WriteString("wrapped:"); b.WriteString(label); b.WriteString("=")
        b.WriteString(rec["type"]); b.WriteString(":"); b.WriteString(rec["ct"]) ; b.WriteString("\n")
    }
    sum := sha256.Sum256([]byte(b.String()))

    // 使用 Transit 簽章（要求 SHA-256 digest）
    if err := w.store.EnsureTransitKey(ctx, signerKey); err != nil { return "", err }
    der, err := w.store.TransitSign(ctx, signerKey, sum[:])
    if err != nil { return "", err }
    return b64.StdEncoding.EncodeToString(der), nil
}

