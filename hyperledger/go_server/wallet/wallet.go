package wallet

// Integrated SQLite-backed Wallet — no separate store layer
// Requires go_server/database.DB already opened. Provides PutFile, PutRaw, Get, etc.

import (
    "context"
    "crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"log"
	"strings"
	"sync"
    "crypto/sha256"
    "encoding/asn1"
    "math/big"

	"go_server/database"
    vs "go_server/vaultstore"

	"github.com/hyperledger/fabric-gateway/pkg/identity"
)

/**
 * @notice 錢包介面，定義身份寫入與讀取能力
 * @dev 僅支援 signerUri 引用模式（DB 不存 PEM）
 */
type WalletInterface interface {
    PutReference(userID string, mspID, signerURI string, certURI string) error
	Exists(label string) bool
    GetResolved(userID string) (*Entry, bool)
}

/**
 * @notice 錢包條目，包含 X.509 身份與交易簽章器
 * @dev 供 Fabric Gateway 建立合約並簽署交易使用
 */
type Entry struct {
    ID        *identity.X509Identity
    Signer    identity.Sign
    Cert      *x509.Certificate
    Priv      *ecdsa.PrivateKey
    SignerURI string
    MspID     string
}

/**
 * @notice SQLite 支援的錢包實作
 * @dev 使用互斥鎖確保執行緒安全，並以 once 確保資料表只建立一次
 */
type Wallet struct {
	once sync.Once
	mu   sync.RWMutex
}

/**
 * @notice 確保建立錢包資料表
 * @dev 建立 `wallet(label TEXT PRIMARY KEY, content BLOB)`，若已存在則忽略
 */
func ensureTable() {
	const ddl = `CREATE TABLE IF NOT EXISTS wallet (
        label   TEXT PRIMARY KEY,
        content BLOB NOT NULL
    );`
	database.DB.Exec(ddl)
}

/**
 * @notice 建立錢包執行個體
 * @dev 於第一次呼叫時建立資料表
 * @return *Wallet 錢包執行個體
 */
func New() *Wallet {
	w := &Wallet{}
	w.once.Do(ensureTable)
	return w
}

// PutReference 僅存外部引用（不存憑證與私鑰），透過 signerUri 指向 Vault 等外部儲存
func (w *Wallet) PutReference(userID string, mspID, signerURI string, certURI string) error {
    payload := map[string]any{
        "mspId":     mspID,
        "signerUri": signerURI,
        "certUri":   certURI,
    }
    content, _ := json.Marshal(payload)
    w.mu.Lock()
    defer w.mu.Unlock()
    _, err := database.DB.Exec(`INSERT INTO wallet(label,content) VALUES(?,?)
        ON CONFLICT(label) DO UPDATE SET content=excluded.content`, userID, content)
    return err
}

// GetResolved：僅支援引用模式（signerUri），自 Vault 還原 cert/key 於記憶體中建立 ID 與 Signer
func (w *Wallet) GetResolved(userID string) (*Entry, bool) {
    w.mu.RLock()
    row := database.DB.QueryRow(`SELECT content FROM wallet WHERE label=?`, userID)
    var blob []byte
    err := row.Scan(&blob)
    w.mu.RUnlock()
    if err != nil {
        log.Printf("[Wallet] GetResolved: wallet row not found user=%s err=%v", userID, err)
        return nil, false
    }

    var payload struct {
        MspID     string `json:"mspId"`
        SignerURI string `json:"signerUri"`
        CertURI   string `json:"certUri"`
    }
    if err := json.Unmarshal(blob, &payload); err != nil {
        log.Printf("[Wallet] GetResolved: bad payload user=%s err=%v", userID, err)
        return nil, false
    }
    if payload.SignerURI == "" {
        log.Printf("[Wallet] GetResolved: empty signerUri user=%s", userID)
        return nil, false
    }

    store, err := vs.NewFromEnv(); if err != nil { log.Printf("[Wallet] GetResolved: vault init failed user=%s err=%v", userID, err); return nil, false }
    // 憑證路徑僅使用 certUri（必須為 kv:// 前綴）
    if payload.CertURI == "" { log.Printf("[Wallet] GetResolved: empty certUri user=%s", userID); return nil, false }
    if !strings.HasPrefix(payload.CertURI, "kv://") { log.Printf("[Wallet] GetResolved: certUri must start with kv:// user=%s certUri=%s", userID, payload.CertURI); return nil, false }
    kvPath := strings.TrimPrefix(payload.CertURI, "kv://")
    _, _, certPEM, err := store.ReadPath(context.Background(), kvPath)
    if err != nil { log.Printf("[Wallet] GetResolved: vault read cert failed user=%s path=%s err=%v", userID, kvPath, err); return nil, false }
    cert, err := identity.CertificateFromPEM(certPEM); if err != nil { log.Printf("[Wallet] GetResolved: parse cert failed user=%s err=%v", userID, err); return nil, false }
    id, _ := identity.NewX509Identity(payload.MspID, cert)

    // 構建 Fabric 所需的 identity.Sign（Transit 簽章 + low-S 正規化）
    keyName := strings.TrimPrefix(payload.SignerURI, "transit://")
    pubKey, _ := cert.PublicKey.(*ecdsa.PublicKey)
    fabricSigner := func(digest []byte) ([]byte, error) {
        der, err := store.TransitSign(context.Background(), keyName, digest)
        if err != nil { return nil, err }
        // low-S normalization for ECDSA signatures (Fabric requires canonical signatures)
        if pubKey != nil && pubKey.Params() != nil {
            var sig struct{ R, S *big.Int }
            if _, err := asn1.Unmarshal(der, &sig); err == nil && sig.R != nil && sig.S != nil {
                n := pubKey.Params().N
                half := new(big.Int).Rsh(new(big.Int).Set(n), 1)
                if sig.S.Cmp(half) == 1 { // S > N/2
                    sig.S = new(big.Int).Sub(n, sig.S)
                    if fixed, mErr := asn1.Marshal(sig); mErr == nil {
                        der = fixed
                    }
                }
            }
        }
        return der, nil
    }
    // 自檢：用 cert 公鑰驗證 Transit 簽章是否可被驗證（固定探針）
    if pub, ok := cert.PublicKey.(*ecdsa.PublicKey); ok {
        probe := sha256.Sum256([]byte("wallet-probe"))
        sig, err := fabricSigner(probe[:])
        if err != nil {
            log.Printf("[Wallet] GetResolved: probe sign failed user=%s err=%v", userID, err)
            return nil, false
        }
        if !ecdsa.VerifyASN1(pub, probe[:], sig) {
            log.Printf("[Wallet] GetResolved: probe verify failed user=%s", userID)
            return nil, false
        }
    }
    log.Printf("[Wallet] GetResolved: success user=%s msp=%s key=%s certPath=%s", userID, payload.MspID, keyName, kvPath)
    return &Entry{ID: id, Signer: fabricSigner, Cert: cert, SignerURI: payload.SignerURI, MspID: payload.MspID}, true
}

// Exists 檢查標籤是否存在
func (w *Wallet) Exists(label string) bool {
	row := database.DB.QueryRow(`SELECT 1 FROM wallet WHERE label=?`, label)
	var dummy int
	return row.Scan(&dummy) == nil
}

// Remove 自錢包刪除指定標籤
func (w *Wallet) Remove(label string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	_, err := database.DB.Exec(`DELETE FROM wallet WHERE label=?`, label)
	return err
}

// List 列出目前錢包中的所有標籤
func (w *Wallet) List() ([]string, error) {
	rows, err := database.DB.Query(`SELECT label FROM wallet`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ls []string
	for rows.Next() {
		var l string
		rows.Scan(&l)
		ls = append(ls, l)
	}
	return ls, nil
}
