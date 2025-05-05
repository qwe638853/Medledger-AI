package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("無法開啟 SQLite 資料庫: %v", err)
	}

	createStmt := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password TEXT,
		name TEXT,
		date TEXT,
		email TEXT,
		phone TEXT
	);`
	_, err = DB.Exec(createStmt)
	if err != nil {
		return fmt.Errorf("建立資料表失敗: %v", err)
	}

	log.Println("✅ SQLite 初始化成功")
	return nil
}

// 查詢帳號是否存在
func IsUserExists(username string) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	return count > 0, err
}

// 新增使用者
func InsertUser(username, password, name, date, email, phone string) error {
	_, err := DB.Exec("INSERT INTO users(username, password, name, date, email, phone) VALUES (?, ?, ?, ?, ?, ?)",
		username, password, name, date, email, phone)
	return err
}
