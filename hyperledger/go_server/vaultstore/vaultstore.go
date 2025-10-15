package vaultstore

import (
    "context"
    "fmt"
    "os"
    "strings"
    "encoding/pem"
    "time"
    "log"

    vault "github.com/hashicorp/vault/api"
    base64Std "encoding/base64"
    "crypto/x509"
    "crypto/ecdsa"
)

// Store 封裝 Vault 用戶端與預設路徑
type Store struct {
    client *vault.Client
    mount  string // e.g. "kv" (對應於 /v1/kv)
    transitMount string // e.g. "transit"
}

// NewFromEnv 透過環境變數建立 Store（統一使用 KV v2）
// 需要：VAULT_ADDR、VAULT_TOKEN；可選：VAULT_MOUNT（預設 kv）
func NewFromEnv() (*Store, error) {
    cfg := vault.DefaultConfig()
    // 讀取所有 Vault 相關環境變數（如 VAULT_ADDR/VAULT_CACERT/VAULT_CAPATH/VAULT_CLIENT_CERT...）
    _ = cfg.ReadEnvironment()
    if addr := os.Getenv("VAULT_ADDR"); addr != "" {
        cfg.Address = addr
    }
    cli, err := vault.NewClient(cfg)
    if err != nil {
        return nil, err
    }
    if tok := os.Getenv("VAULT_TOKEN"); tok != "" {
        cli.SetToken(tok)
    }
    mount := os.Getenv("VAULT_MOUNT")
    if mount == "" { mount = "kv" }
    tMount := os.Getenv("VAULT_TRANSIT_MOUNT")
    if tMount == "" { tMount = "transit" }
    log.Printf("[Vault] init client addr=%s kvMount=%s transitMount=%s", cfg.Address, mount, tMount)
    return &Store{client: cli, mount: mount, transitMount: tMount}, nil
}

// WriteUserMaterial 將使用者的 CSR/私鑰/憑證寫入 Vault（KV v2）
// path: "{mount}/data/users/{userID}"，資料放在 data 欄位
func (s *Store) WriteUserMaterial(ctx context.Context, userID string, csrPEM, keyPEM, certPEM []byte) error {
    _, err := s.client.Logical().WriteWithContext(ctx, fmt.Sprintf("%s/data/users/%s", s.mount, userID), map[string]interface{}{
        "data": map[string]interface{}{
            "csr":  string(csrPEM),
            "key":  string(keyPEM),
            "cert": string(certPEM),
        },
    })
    return err
}

// WriteInsurerMaterial 將保險業者 CSR/私鑰/憑證寫入 Vault（KV v2）
func (s *Store) WriteInsurerMaterial(ctx context.Context, insurerID string, csrPEM, keyPEM, certPEM []byte) error {
    _, err := s.client.Logical().WriteWithContext(ctx, fmt.Sprintf("%s/data/insurers/%s", s.mount, insurerID), map[string]interface{}{
        "data": map[string]interface{}{
            "csr":  string(csrPEM),
            "key":  string(keyPEM),
            "cert": string(certPEM),
        },
    })
    return err
}

// WriteClinicMaterial 將健檢中心 CSR/私鑰/憑證寫入 Vault（KV v2）
func (s *Store) WriteClinicMaterial(ctx context.Context, clinicID string, csrPEM, keyPEM, certPEM []byte) error {
    _, err := s.client.Logical().WriteWithContext(ctx, fmt.Sprintf("%s/data/clinics/%s", s.mount, clinicID), map[string]interface{}{
        "data": map[string]interface{}{
            "csr":  string(csrPEM),
            "key":  string(keyPEM),
            "cert": string(certPEM),
        },
    })
    return err
}

// WritePlatformMaterial 將平台 CSR/私鑰/憑證寫入 Vault（KV v2）
func (s *Store) WritePlatformMaterial(ctx context.Context, csrPEM, keyPEM, certPEM []byte) error {
    _, err := s.client.Logical().WriteWithContext(ctx, fmt.Sprintf("%s/data/platform", s.mount), map[string]interface{}{
        "data": map[string]interface{}{
            "csr":  string(csrPEM),
            "key":  string(keyPEM),
            "cert": string(certPEM),
        },
    })
    return err
}

// 讀取類別材料（KV v2）
func (s *Store) ReadUserMaterial(ctx context.Context, userID string) (csrPEM, keyPEM, certPEM []byte, err error) {
    sec, err := s.client.Logical().ReadWithContext(ctx, fmt.Sprintf("%s/data/users/%s", s.mount, userID))
    if err != nil { return nil, nil, nil, err }
    if sec == nil || sec.Data == nil { return nil, nil, nil, fmt.Errorf("vault: empty secret") }
    data, _ := sec.Data["data"].(map[string]interface{})
    if data == nil { return nil, nil, nil, fmt.Errorf("vault: missing data field") }
    return []byte(fmt.Sprint(data["csr"])), []byte(fmt.Sprint(data["key"])), []byte(fmt.Sprint(data["cert"])), nil
}

func (s *Store) ReadInsurerMaterial(ctx context.Context, insurerID string) (csrPEM, keyPEM, certPEM []byte, err error) {
    sec, err := s.client.Logical().ReadWithContext(ctx, fmt.Sprintf("%s/data/insurers/%s", s.mount, insurerID))
    if err != nil { return nil, nil, nil, err }
    if sec == nil || sec.Data == nil { return nil, nil, nil, fmt.Errorf("vault: empty secret") }
    data, _ := sec.Data["data"].(map[string]interface{})
    if data == nil { return nil, nil, nil, fmt.Errorf("vault: missing data field") }
    return []byte(fmt.Sprint(data["csr"])), []byte(fmt.Sprint(data["key"])), []byte(fmt.Sprint(data["cert"])), nil
}

func (s *Store) ReadClinicMaterial(ctx context.Context, clinicID string) (csrPEM, keyPEM, certPEM []byte, err error) {
    sec, err := s.client.Logical().ReadWithContext(ctx, fmt.Sprintf("%s/data/clinics/%s", s.mount, clinicID))
    if err != nil { return nil, nil, nil, err }
    if sec == nil || sec.Data == nil { return nil, nil, nil, fmt.Errorf("vault: empty secret") }
    data, _ := sec.Data["data"].(map[string]interface{})
    if data == nil { return nil, nil, nil, fmt.Errorf("vault: missing data field") }
    return []byte(fmt.Sprint(data["csr"])), []byte(fmt.Sprint(data["key"])), []byte(fmt.Sprint(data["cert"])), nil
}

func (s *Store) ReadPlatformMaterial(ctx context.Context) (csrPEM, keyPEM, certPEM []byte, err error) {
    sec, err := s.client.Logical().ReadWithContext(ctx, fmt.Sprintf("%s/data/platform", s.mount))
    if err != nil { return nil, nil, nil, err }
    if sec == nil || sec.Data == nil { return nil, nil, nil, fmt.Errorf("vault: empty secret") }
    data, _ := sec.Data["data"].(map[string]interface{})
    if data == nil { return nil, nil, nil, fmt.Errorf("vault: missing data field") }
    return []byte(fmt.Sprint(data["csr"])), []byte(fmt.Sprint(data["key"])), []byte(fmt.Sprint(data["cert"])), nil
}

// ReadPath 以相對路徑讀取（KV v2）：會從 {mount}/data/{relPath} 取出 data.csr/key/cert
func (s *Store) ReadPath(ctx context.Context, relPath string) (csrPEM, keyPEM, certPEM []byte, err error) {
    sec, err := s.client.Logical().ReadWithContext(ctx, fmt.Sprintf("%s/data/%s", s.mount, relPath))
    if err != nil { return nil, nil, nil, err }
    if sec == nil || sec.Data == nil { return nil, nil, nil, fmt.Errorf("vault: empty secret") }
    data, _ := sec.Data["data"].(map[string]interface{})
    if data == nil { return nil, nil, nil, fmt.Errorf("vault: missing data field") }
    return []byte(fmt.Sprint(data["csr"])), []byte(fmt.Sprint(data["key"])), []byte(fmt.Sprint(data["cert"])), nil
}

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

// TransitDecrypt 將 Transit ciphertext 解密為原始 dataKey
func (s *Store) TransitDecrypt(ctx context.Context, keyName string, ciphertext string) ([]byte, error) {
    if keyName == "" || ciphertext == "" { return nil, fmt.Errorf("empty args") }
    path := fmt.Sprintf("%s/decrypt/%s", s.transitMount, keyName)
    if _, ok := ctx.Deadline(); !ok {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
        defer cancel()
    }
    log.Printf("[Vault] TransitDecrypt path=/v1/%s", path)
    sec, err := s.client.Logical().WriteWithContext(ctx, path, map[string]interface{}{
        "ciphertext": ciphertext,
    })
    if err != nil { return nil, err }
    if sec == nil || sec.Data == nil { return nil, fmt.Errorf("empty transit response") }
    data, _ := sec.Data["plaintext"].(string)
    if data == "" { return nil, fmt.Errorf("missing plaintext") }
    raw, err := decodeB64(data)
    if err != nil { return nil, err }
    return raw, nil
}

// --- helpers ---
func encodeB64(b []byte) string { return base64Std.StdEncoding.EncodeToString(b) }
func decodeB64(s string) ([]byte, error) { return base64Std.StdEncoding.DecodeString(s) }

// TransitSign 使用 Transit 對 digest 簽章（ECDSA-P256；預設 asn1 DER）
// digest 必須為 SHA-256 結果；此函數會以 prehashed=true 呼叫 Vault。
func (s *Store) TransitSign(ctx context.Context, keyName string, digest []byte) ([]byte, error) {
    if keyName == "" || len(digest) == 0 { return nil, fmt.Errorf("empty args") }
    path := fmt.Sprintf("%s/sign/%s", s.transitMount, keyName)
    body := map[string]interface{}{
        "input":                 encodeB64(digest),
        "prehashed":             true,
        "hash_algorithm":        "sha2-256",
        "marshaling_algorithm":  "asn1",
    }
    if _, ok := ctx.Deadline(); !ok {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
        defer cancel()
    }
    log.Printf("[Vault] TransitSign path=/v1/%s", path)
    sec, err := s.client.Logical().WriteWithContext(ctx, path, body)
    if err != nil { return nil, err }
    if sec == nil || sec.Data == nil { return nil, fmt.Errorf("empty transit sign response") }
    sig, _ := sec.Data["signature"].(string)
    if sig == "" { return nil, fmt.Errorf("missing signature") }
    // 解析 vault:vX:<b64>
    parts := strings.Split(sig, ":")
    b64part := parts[len(parts)-1]
    der, err := decodeB64(b64part)
    if err != nil { return nil, err }
    return der, nil
}

// TransitGetPublicKey 讀取 Transit 簽章金鑰的公鑰（ECDSA-P256）
func (s *Store) TransitGetPublicKey(ctx context.Context, keyName string) (*ecdsa.PublicKey, error) {
    if keyName == "" { return nil, fmt.Errorf("empty keyName") }
    path := fmt.Sprintf("%s/keys/%s", s.transitMount, keyName)
    if _, ok := ctx.Deadline(); !ok {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
        defer cancel()
    }
    log.Printf("[Vault] TransitGetPublicKey path=/v1/%s", path)
    sec, err := s.client.Logical().ReadWithContext(ctx, path)
    if err != nil { return nil, err }
    if sec == nil || sec.Data == nil { return nil, fmt.Errorf("empty transit key response") }
    // 嘗試從 keys.latest 或 public_key 取公鑰 PEM
    // 先取 top-level public_key
    if pk, ok := sec.Data["public_key"].(string); ok && pk != "" {
        return parseECDSAPublicKey([]byte(pk))
    }
    // 再嘗試 keys -> map -> 任一版本的 public_key
    if km, ok := sec.Data["keys"].(map[string]interface{}); ok {
        for _, v := range km {
            if rec, ok := v.(map[string]interface{}); ok {
                if pk, ok := rec["public_key"].(string); ok && pk != "" {
                    return parseECDSAPublicKey([]byte(pk))
                }
            }
        }
    }
    return nil, fmt.Errorf("no public key in transit response")
}

// EnsureTransitKey 確保指定 Transit 簽章金鑰存在；若不存在則以 ecdsa-p256 建立
func (s *Store) EnsureTransitKey(ctx context.Context, keyName string) error {
    if keyName == "" { return fmt.Errorf("empty keyName") }
    path := fmt.Sprintf("%s/keys/%s", s.transitMount, keyName)
    // 試讀
    sec, err := s.client.Logical().ReadWithContext(ctx, path)
    if err == nil && sec != nil && sec.Data != nil { return nil }
    // 嘗試建立
    _, werr := s.client.Logical().WriteWithContext(ctx, path, map[string]interface{}{"type": "ecdsa-p256"})
    return werr
}

// EnsureTransitKeyOfType 確保指定類型的 Transit 金鑰存在（例如 aes256-gcm96 用於加解密）
func (s *Store) EnsureTransitKeyOfType(ctx context.Context, keyName string, keyType string) error {
    if keyName == "" || keyType == "" { return fmt.Errorf("empty args") }
    path := fmt.Sprintf("%s/keys/%s", s.transitMount, keyName)
    sec, err := s.client.Logical().ReadWithContext(ctx, path)
    if err == nil && sec != nil && sec.Data != nil { return nil }
    _, werr := s.client.Logical().WriteWithContext(ctx, path, map[string]interface{}{"type": keyType})
    return werr
}

func parseECDSAPublicKey(pemBytes []byte) (*ecdsa.PublicKey, error) {
    block, _ := pem.Decode(pemBytes)
    if block == nil { return nil, fmt.Errorf("bad pem") }
    pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil { return nil, err }
    pub, ok := pubAny.(*ecdsa.PublicKey)
    if !ok { return nil, fmt.Errorf("not ecdsa pub") }
    return pub, nil
}



