package signer

import (
    "context"
    "crypto"
    "crypto/x509"
    "encoding/pem"
    "errors"
    "io"

    vs "go_server/vaultstore"
)

// TransitSigner 實作 crypto.Signer，將 Sign 委派給 Vault Transit
// 注意：Public 需回傳對應公鑰；此實作提供兩種來源：
// 1) 從 PEM 載入公鑰（建議把 certPEM 存在 Vault KV）
// 2) 或後續擴充：從 Transit 讀取公鑰（未在此實作）
type TransitSigner struct {
    store   *vs.Store
    keyName string
    pub     crypto.PublicKey
}

// NewTransitSignerFromCert 建立以憑證公鑰為 Public 的 Transit 簽章器
func NewTransitSignerFromCert(store *vs.Store, keyName string, certPEM []byte) (*TransitSigner, error) {
    cert, err := x509.ParseCertificate(decodePEMBytes(certPEM))
    if err != nil { return nil, err }
    return &TransitSigner{store: store, keyName: keyName, pub: cert.PublicKey}, nil
}

// NewTransitSignerWithPublicKey 以 Transit 公鑰建立簽章器（用於產 CSR）
func NewTransitSignerWithPublicKey(store *vs.Store, keyName string, pub crypto.PublicKey) (*TransitSigner, error) {
    if store == nil || keyName == "" || pub == nil { return nil, errors.New("bad args") }
    return &TransitSigner{store: store, keyName: keyName, pub: pub}, nil
}

func (t *TransitSigner) Public() crypto.PublicKey {
    return t.pub
}

// Sign 期望 digest 已為 SHA-256；alg 需為 &crypto.SHA256
func (t *TransitSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
    if opts == nil || opts.HashFunc() != crypto.SHA256 {
        return nil, errors.New("transit signer requires SHA-256 digest")
    }
    der, err := t.store.TransitSign(context.Background(), t.keyName, digest)
    if err != nil { return nil, err }
    return der, nil
}

// --- helpers ---
func decodePEMBytes(pemBytes []byte) []byte {
    block, _ := pem.Decode(pemBytes)
    if block != nil { return block.Bytes }
    return pemBytes
}


