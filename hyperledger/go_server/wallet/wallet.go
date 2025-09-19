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

/**
 * @notice 錢包介面，定義身份寫入與讀取能力
 * @dev 以 userID 作為標籤管理身份，提供儲存、查詢與讀取功能
 */
type WalletInterface interface {
	PutFile(userID, certPath, keyPath, mspID string) error
	Exists(label string) bool
	Get(userID string) (*Entry, bool)
}

/**
 * @notice 錢包條目，包含 X.509 身份與交易簽章器
 * @dev 供 Fabric Gateway 建立合約並簽署交易使用
 */
type Entry struct {
	ID     *identity.X509Identity
	Signer identity.Sign
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

/**
 * @notice 讀取 PEM 檔並寫入錢包
 * @dev 讀取 cert/key 檔後委派至 PutRaw 儲存
 * @param userID 使用者識別
 * @param certPath 憑證檔路徑
 * @param keyPath 私鑰檔路徑
 * @param mspID MSP 名稱（例如 Org1MSP）
 * @return error 讀檔或存檔失敗
 */
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

/**
 * @notice 以位元組形式將憑證與私鑰寫入錢包
 * @dev 驗證 PEM 格式並以 JSON 存入 SQLite，鍵為 userID
 * @param userID 使用者識別
 * @param certPEM 憑證 PEM 內容
 * @param keyPEM 私鑰 PEM 內容
 * @param mspID MSP 名稱（例如 Org1MSP）
 * @return error 解析或存取失敗
 */
func (w *Wallet) PutRaw(userID string, certPEM, keyPEM []byte, mspID string) error {
	// 讀取憑證與私鑰
	cert, err := identity.CertificateFromPEM(certPEM)
	if err != nil {
		return err
	}
	// TODO: 可以針對私鑰及憑證做更多檢查

	// 建立 JSON 格式
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

	return nil
}

/**
 * @notice 由錢包取回身份與簽章器
 * @dev 讀取資料庫 JSON，還原 X.509 身份與 private key signer
 * @param userID 使用者識別
 * @return *Entry 還原的錢包條目, bool 是否存在
 */
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

/**
 * @notice 檢查標籤是否存在於錢包
 * @param label userID 或其他標籤
 * @return bool 是否存在
 */
func (w *Wallet) Exists(label string) bool {
	row := database.DB.QueryRow(`SELECT 1 FROM wallet WHERE label=?`, label)
	var dummy int
	return row.Scan(&dummy) == nil
}

/**
 * @notice 自錢包刪除指定標籤
 * @param label userID 或其他標籤
 * @return error 刪除失敗
 */
func (w *Wallet) Remove(label string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	_, err := database.DB.Exec(`DELETE FROM wallet WHERE label=?`, label)
	return err
}

/**
 * @notice 列出目前錢包中的所有標籤
 * @return []string 標籤清單, error 查詢失敗
 */
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
