package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// HashString 將字串轉換為 SHA256 雜湊值
func HashString(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func InitDB(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("無法開啟 SQLite 資料庫: %v", err)
	}

	// 用戶表
	createUserStmt := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password TEXT,
		name TEXT,
		date TEXT,
		email TEXT,
		phone TEXT
	);`

	_, err = DB.Exec(createUserStmt)
	if err != nil {
		return fmt.Errorf("建立用戶資料表失敗: %v", err)
	}

	// 保險業者表
	createInsurerStmt := `
	CREATE TABLE IF NOT EXISTS insurers (
		insurer_id TEXT PRIMARY KEY,
		password TEXT,
		company_name TEXT,
		contact_person TEXT,
		email TEXT,
		phone TEXT
	);`

	_, err = DB.Exec(createInsurerStmt)
	if err != nil {
		return fmt.Errorf("建立保險業者資料表失敗: %v", err)
	}

	log.Println("✅ SQLite 初始化成功")
	return nil
}

// 查詢用戶帳號是否存在
func IsUserExists(username string) (bool, error) {
	hashedUsername := HashString(username)
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", hashedUsername).Scan(&count)
	return count > 0, err
}

// 查詢保險業者帳號是否存在
func IsInsurerExists(insurerId string) (bool, error) {
	hashedInsurerId := HashString(insurerId)
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM insurers WHERE insurer_id = ?", hashedInsurerId).Scan(&count)
	return count > 0, err
}

// 新增用戶
func InsertUser(username, password, name, date, email, phone string) error {
	log.Printf("[Debug] 新增用戶: %s", username)
	hashedUsername := HashString(username)
	hashedPassword := HashString(password)
	_, err := DB.Exec("INSERT INTO users(username, password, name, date, email, phone) VALUES (?, ?, ?, ?, ?, ?)",
		hashedUsername, hashedPassword, name, date, email, phone)
	return err
}

// 新增保險業者
func InsertInsurer(insurerId, password, companyName, contactPerson, email, phone string) error {
	log.Printf("[Debug] 新增保險業者: %s", insurerId)
	hashedInsurerId := HashString(insurerId)
	hashedPassword := HashString(password)
	_, err := DB.Exec("INSERT INTO insurers(insurer_id, password, company_name, contact_person, email, phone) VALUES (?, ?, ?, ?, ?, ?)",
		hashedInsurerId, hashedPassword, companyName, contactPerson, email, phone)
	return err
}

// 取得用戶密碼
func GetUserPassword(username string) (string, error) {
	hashedUsername := HashString(username)
	var password string
	err := DB.QueryRow("SELECT password FROM users WHERE username = ?", hashedUsername).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

// 取得保險業者密碼
func GetInsurerPassword(insurerId string) (string, error) {
	hashedInsurerId := HashString(insurerId)
	var password string
	err := DB.QueryRow("SELECT password FROM insurers WHERE insurer_id = ?", hashedInsurerId).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

// InsurerInfo 存儲保險業者的基本資訊
type InsurerInfo struct {
	InsurerID    string
	CompanyName  string
	Name         string
	Email        string
	Phone        string
}

// UserInfo 存儲用戶的基本資訊
type UserInfo struct {
	Username string
	Name     string
	Date     string
	Email    string
	Phone    string
}

// GetInsurerByHash 根據雜湊值獲取保險業者資訊
func GetInsurerByHash(insurerHash string) (*InsurerInfo, error) {
	var info InsurerInfo
	err := DB.QueryRow(`
		SELECT insurer_id, company_name, contact_person, email, phone 
		FROM insurers 
		WHERE insurer_id = ?`, insurerHash).Scan(
		&info.InsurerID,
		&info.CompanyName,
		&info.Name,
		&info.Email,
		&info.Phone,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("找不到對應的保險業者資訊")
		}
		return nil, fmt.Errorf("查詢保險業者資訊失敗: %v", err)
	}
	return &info, nil
}

// GetUserByHash 根據雜湊值獲取用戶資訊
func GetUserByHash(userHash string) (*UserInfo, error) {
	var info UserInfo
	err := DB.QueryRow(`
		SELECT username, name, date, email, phone 
		FROM users 
		WHERE username = ?`, userHash).Scan(
		&info.Username,
		&info.Name,
		&info.Date,
		&info.Email,
		&info.Phone,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("找不到對應的用戶資訊")
		}
		return nil, fmt.Errorf("查詢用戶資訊失敗: %v", err)
	}
	return &info, nil
}
