package wallet

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/sha256"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "net/http/httptest"
    "math/big"
    "net/http"
    "testing"
    "time"

    db "go_server/database"
    vs "go_server/vaultstore"
)

// genKeyAndSelfCert 產生一組 ECDSA P-256 私鑰與自簽憑證（測試用）
func genKeyAndSelfCert(t *testing.T, cn string) (certPEM, keyPEM []byte) {
    t.Helper()
    // 產生 ECDSA P-256 私鑰
    priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil { t.Fatalf("gen key: %v", err) }

    // 建立自簽 X.509 憑證
    tmpl := &x509.Certificate{
        SerialNumber: big.NewInt(time.Now().UnixNano()),
        Subject: pkix.Name{CommonName: cn},
        NotBefore: time.Now().Add(-time.Hour),
        NotAfter:  time.Now().Add(24 * time.Hour),
        KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
        IsCA:      true,
        BasicConstraintsValid: true,
    }
    derCert, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &priv.PublicKey, priv)
    if err != nil { t.Fatalf("create cert: %v", err) }

    certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derCert})
    derKey, err := x509.MarshalPKCS8PrivateKey(priv)
    if err != nil { t.Fatalf("marshal key: %v", err) }
    keyPEM = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: derKey})
    return
}

// TestWalletPutRawGet
// 目的：驗證 PutRaw 寫入 PEM 後，Get 能還原出 X.509 身份與可用的簽章器
func TestWalletPutRawGet(t *testing.T) {
    // 初始化測試用 SQLite（記憶體）
    if err := db.InitDB(":memory:"); err != nil {
        t.Fatalf("init db: %v", err)
    }
    w := New()

    // 準備測試用憑證與私鑰
    certPEM, keyPEM := genKeyAndSelfCert(t, "alice")
    if err := w.PutRaw("alice", certPEM, keyPEM, "Org1MSP"); err != nil {
        t.Fatalf("PutRaw: %v", err)
    }

    // 讀取並檢查
    e, ok := w.Get("alice")
    if !ok || e == nil { t.Fatalf("Get: not found") }
    if e.ID == nil || e.Signer == nil || e.Cert == nil {
        t.Fatalf("entry not fully restored: id=%v signer=%v cert=%v", e.ID, e.Signer, e.Cert)
    }

    // 進一步驗證簽章器可用（對 SHA-256 雜湊簽章）
    digest := sha256.Sum256([]byte("hello"))
    sig, err := e.Signer(digest[:])
    if err != nil || len(sig) == 0 {
        t.Fatalf("sign failed: %v len=%d", err, len(sig))
    }
}

// TestWalletPutReferenceGet
// 目的：驗證引用模式寫入（不存 PEM），Get 僅回傳 signerUri 與 mspId，其他留空
func TestWalletPutReferenceGet(t *testing.T) {
    if err := db.InitDB(":memory:"); err != nil {
        t.Fatalf("init db: %v", err)
    }
    w := New()

    if err := w.PutReference("bob", "Org1MSP", "transit://users/bob", "kv://users/bob"); err != nil {
        t.Fatalf("PutReference: %v", err)
    }
    e, ok := w.Get("bob")
    if !ok || e == nil { t.Fatalf("Get: not found") }
    if e.SignerURI != "users/bob" { t.Fatalf("signerUri mismatch: %s", e.SignerURI) }
    if e.MspID != "Org1MSP" { t.Fatalf("mspId mismatch: %s", e.MspID) }
    if e.ID != nil || e.Signer != nil || e.Cert != nil {
        t.Fatalf("expected unresolved entry, got id=%v signer=%v cert=%v", e.ID, e.Signer, e.Cert)
    }
}

// TestGetResolvedWithVault
// 目的：驗證引用模式下，透過 GetResolved 能從 Vault 拉回 cert/key 並填入 Entry
func TestGetResolvedWithVault(t *testing.T) {
    // 建立 fake Vault KV v2 伺服器
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 只處理 /v1/kv/data/users/alice
        if r.Method == http.MethodPost && r.URL.Path == "/v1/kv/data/users/alice" {
            w.Write([]byte(`{"data":{"metadata":{}}}`))
            return
        }
        if r.Method == http.MethodGet && r.URL.Path == "/v1/kv/data/users/alice" {
            // 回傳先前寫入的資料（固定值）
            // 測試使用先前 genKeyAndSelfCert 生成的資料不易；此處使用另一組自簽快速生成
            certPEM, keyPEM := genKeyAndSelfCert(t, "alice")
            _ = keyPEM
            resp := `{"data":{"data":{"csr":"CSR","key":"` + string(keyPEM) + `","cert":"` + string(certPEM) + `"},"metadata":{}}}`
            w.Write([]byte(resp))
            return
        }
        http.NotFound(w, r)
    }))
    defer srv.Close()

    // 設定 Vault env
    t.Setenv("VAULT_ADDR", srv.URL)
    t.Setenv("VAULT_TOKEN", "dev-token")
    t.Setenv("VAULT_MOUNT", "kv")

    if err := db.InitDB(":memory:"); err != nil { t.Fatalf("init db: %v", err) }
    w := New()

    // 只存引用
    if err := w.PutReference("alice", "Org1MSP", "transit://users/alice", "kv://users/alice"); err != nil {
        t.Fatalf("PutReference: %v", err)
    }
    // 解析
    e, ok := w.GetResolved("alice")
    if !ok || e == nil || e.ID == nil || e.Signer == nil || e.Cert == nil || e.Priv == nil {
        t.Fatalf("GetResolved failed: ok=%v id=%v signer=%v cert=%v priv=%v", ok, e.ID, e.Signer, e.Cert, e.Priv)
    }
}



