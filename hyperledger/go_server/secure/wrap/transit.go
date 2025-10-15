package wrap

import (
    "context"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    b64 "encoding/base64"
    "encoding/json"
    "log"
    "os"
    "strings"

    vs "go_server/vaultstore"
)

// TransitWrapper 封裝以 Vault Transit 進行資料金鑰包裝/解包與 AES-GCM 主密文加解
type TransitWrapper struct {
    store         *vs.Store
    platformKey   string // 例如 "platform-wrap"（AES 用於 encrypt/decrypt）
}

// NewTransitWrapperFromEnv 以環境變數建立 Wrapper
// 需要：VAULT_ADDR, VAULT_TOKEN；可選：VAULT_MOUNT（kv）、VAULT_TRANSIT_MOUNT（預設 transit）、VAULT_TRANSIT_KEY_PLATFORM（預設 platform）
func NewTransitWrapperFromEnv() (*TransitWrapper, error) {
    st, err := vs.NewFromEnv()
    if err != nil { return nil, err }
    plat := os.Getenv("VAULT_TRANSIT_KEY_PLATFORM")
    if plat == "" { plat = "platform" }
    // 為避免使用 ECDSA key 進行 encrypt，改用獨立 wrap key 後綴
    if !strings.HasSuffix(plat, "-wrap") { plat = plat + "-wrap" }
    log.Printf("[Transit] NewWrapper: platformWrapKey=%s", plat)
    return &TransitWrapper{store: st, platformKey: plat}, nil
}

// EncryptReportTransit：
// 1) 產生隨機 32B dataKey 與 12B nonce
// 2) AES-GCM(dataKey) 加密 plaintext
// 3) 以 Transit 為 clinicKey 與 platformKey 分別包 dataKey
// 4) 回傳 Envelope JSON：{"ciphertext": base64(ct), "nonce": base64(nonce), "wrappedKeys": {label:{"type":"transit","ct":...}}, ...}
func (w *TransitWrapper) EncryptReportTransit(ctx context.Context, plaintext []byte, clinicLabel string, clinicKey string) ([]byte, error) {
    log.Printf("[Transit] EncryptReportTransit: clinicLabel=%s clinicKey=%s platformWrapKey=%s plainLen=%d", clinicLabel, clinicKey, w.platformKey, len(plaintext))
    dk := make([]byte, 32)
    if _, err := rand.Read(dk); err != nil { return nil, err }
    nonce := make([]byte, 12)
    if _, err := rand.Read(nonce); err != nil { return nil, err }
    block, err := aes.NewCipher(dk)
    if err != nil { return nil, err }
    gcm, err := cipher.NewGCM(block)
    if err != nil { return nil, err }
    ct := gcm.Seal(nil, nonce, plaintext, nil)
    log.Printf("[Transit] AES-GCM done: ctLen=%d nonceLen=%d", len(ct), len(nonce))

    // 確保 wrap 用 AES key 存在
    if err := w.store.EnsureTransitKeyOfType(ctx, w.platformKey, "aes256-gcm96"); err != nil { return nil, err }
    // 診所端也用獨立 wrap key（若傳入的是簽章 key 名稱，則補 -wrap）
    wrapClinicKey := clinicKey
    if !strings.HasSuffix(wrapClinicKey, "-wrap") { wrapClinicKey = wrapClinicKey + "-wrap" }
    if err := w.store.EnsureTransitKeyOfType(ctx, wrapClinicKey, "aes256-gcm96"); err != nil { return nil, err }

    // Transit 包 dataKey（AES wrap）
    ctPlatform, err := w.store.TransitEncrypt(ctx, w.platformKey, dk)
    if err != nil { log.Printf("[Transit] TransitEncrypt(platform) error: %v", err); return nil, err }
    ctClinic, err := w.store.TransitEncrypt(ctx, wrapClinicKey, dk)
    if err != nil { log.Printf("[Transit] TransitEncrypt(clinic %s) error: %v", clinicKey, err); return nil, err }
    log.Printf("[Transit] Wrapped dataKey: platformCtLen=%d clinicCtLen=%d", len(ctPlatform), len(ctClinic))

    env := map[string]any{
        "ciphertext": b64.StdEncoding.EncodeToString(ct),
        "nonce":      b64.StdEncoding.EncodeToString(nonce),
        "wrappedKeys": map[string]any{
            clinicLabel: map[string]string{"type":"transit", "ct": ctClinic},
            "platform":  map[string]string{"type":"transit", "ct": ctPlatform},
        },
        "enc":  "AES-256-GCM",
        "kdf":  "n/a",
        "curve":"n/a",
    }
    // 預估回傳大小（不含 JSON overhead）
    envCipherB64 := env["ciphertext"].(string)
    log.Printf("[Transit] Envelope ready: cipherB64Len=%d", len(envCipherB64))
    return json.Marshal(env)
}

// AddRecipientTransit：
// 以平台 Transit key 解出 dataKey，再以 newKeyName 為新收件者包一份，更新 wrappedKeys
func (w *TransitWrapper) AddRecipientTransit(ctx context.Context, envJSON []byte, newLabel string, newKeyName string) ([]byte, error) {
    log.Printf("[Transit] AddRecipientTransit: newLabel=%s newKey=%s platformKey=%s envSize=%d", newLabel, newKeyName, w.platformKey, len(envJSON))
    var env struct {
        Ciphertext string                                 `json:"ciphertext"`
        Nonce      string                                 `json:"nonce"`
        WrappedKeys map[string]map[string]string          `json:"wrappedKeys"`
        Enc        string                                  `json:"enc"`
        KDF        string                                  `json:"kdf"`
        Curve      string                                  `json:"curve"`
    }
    if err := json.Unmarshal(envJSON, &env); err != nil { return nil, err }
    plat, ok := env.WrappedKeys["platform"]
    if !ok { return nil, ErrNoPlatformKey }
    ct := plat["ct"]
    dk, err := w.store.TransitDecrypt(ctx, w.platformKey, ct)
    if err != nil { return nil, err }
    log.Printf("[Transit] Decrypt platform wrapped key: dkLen=%d", len(dk))
    // 以 AES wrap key 重新包裝新收件者的 dataKey
    wrapKey := newKeyName
    if !strings.HasSuffix(wrapKey, "-wrap") { wrapKey = wrapKey + "-wrap" }
    if err := w.store.EnsureTransitKeyOfType(ctx, wrapKey, "aes256-gcm96"); err != nil { return nil, err }
    ctNew, err := w.store.TransitEncrypt(ctx, wrapKey, dk)
    if err != nil { return nil, err }
    if env.WrappedKeys == nil { env.WrappedKeys = map[string]map[string]string{} }
    env.WrappedKeys[newLabel] = map[string]string{"type":"transit", "ct": ctNew}
    log.Printf("[Transit] Added recipient: label=%s wrapKey=%s ctLen=%d", newLabel, wrapKey, len(ctNew))
    return json.Marshal(env)
}

// DecryptReportTransit：使用指定角色/標籤與對應 wrap key 解出明文
// label 例："patient"、"insurer"、"platform"；baseKey 例：user-<id>、insurer-<id>、platform
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

// --- errors ---
var ErrNoPlatformKey = &wrapError{"no platform wrapped key"}
type wrapError struct{ msg string }
func (e *wrapError) Error() string { return e.msg }


