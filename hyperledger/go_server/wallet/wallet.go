package wallet

// Integrated SQLite-backed Wallet — no separate store layer
// Requires go_server/database.DB already opened. Provides PutFile, PutRaw, Get, etc.

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"go_server/database"

	"github.com/hyperledger/fabric-gateway/pkg/identity"
)

type WalletInterface interface {
	PutFile(userID, certPath, keyPath, mspID string) error
	Exists(label string) bool
	Get(userID string) (*Entry, bool)
}

type Entry struct {
	ID     *identity.X509Identity
	Signer identity.Sign
}

type Wallet struct {
	once sync.Once
	mu   sync.RWMutex
}

func ensureTable() {
	const ddl = `CREATE TABLE IF NOT EXISTS wallet (
        label   TEXT PRIMARY KEY,
        content BLOB NOT NULL
    );`
	database.DB.Exec(ddl)
}

func New() *Wallet {
	w := &Wallet{}
	w.once.Do(ensureTable)
	return w
}

// PutFile reads PEM files under MSP folder then stores into wallet table.
func (w *Wallet) PutFile(userID, certPath, keyPath, mspID string) error {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("read cert: %w", err)
	}
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("read key: %w", err)
	}
	return w.PutRaw(userID, certPEM, keyPEM, mspID)
}

// PutRaw stores cert & key bytes directly in wallet table.
func (w *Wallet) PutRaw(userID string, certPEM, keyPEM []byte, mspID string) error {
	// 讀取憑證與私鑰
	cert, err := identity.CertificateFromPEM(certPEM)
	if err != nil {
		return err
	}
	privKey, err := identity.PrivateKeyFromPEM(keyPEM)
	if err != nil {
		return err
	}

	id, _ := identity.NewX509Identity(mspID, cert)

	// build JSON manually to keep both cert and key for later retrieval
	payload := map[string]any{
		"mspId":       mspID,
		"certificate": string(certPEM),
		"privateKey":  string(keyPEM),
	}
	content, _ := json.Marshal(payload)

	// 確保這個時間只有一個 goroutine 在執行這段程式碼
	w.mu.Lock()
	defer w.mu.Unlock()
	_, err = database.DB.Exec(`INSERT INTO wallet(label,content) VALUES(?,?)
        ON CONFLICT(label) DO UPDATE SET content=excluded.content`, userID, content)
	if err != nil {
		return err
	}

	// cache signer optional
	_ = privKey
	_ = id // silence unused
	return nil
}

// Get reconstructs Entry from DB JSON; ok=false if not exist or malformed.
func (w *Wallet) Get(userID string) (*Entry, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	row := database.DB.QueryRow(`SELECT content FROM wallet WHERE label=?`, userID)
	var blob []byte
	if err := row.Scan(&blob); err != nil {
		return nil, false
	}
	var payload struct {
		MspID       string `json:"mspId"`
		Certificate string `json:"certificate"`
		PrivateKey  string `json:"privateKey"`
	}
	if err := json.Unmarshal(blob, &payload); err != nil {

		return nil, false
	}
	cert, err := identity.CertificateFromPEM([]byte(payload.Certificate))
	if err != nil {
		return nil, false
	}
	id, _ := identity.NewX509Identity(payload.MspID, cert)
	privKey, err := identity.PrivateKeyFromPEM([]byte(payload.PrivateKey))
	if err != nil {
		return nil, false
	}
	signer, _ := identity.NewPrivateKeySign(privKey)
	return &Entry{ID: id, Signer: signer}, true
}

func (w *Wallet) Exists(label string) bool {
	row := database.DB.QueryRow(`SELECT 1 FROM wallet WHERE label=?`, label)
	var dummy int
	return row.Scan(&dummy) == nil
}

func (w *Wallet) Remove(label string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	_, err := database.DB.Exec(`DELETE FROM wallet WHERE label=?`, label)
	return err
}

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
