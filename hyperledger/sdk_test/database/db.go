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
		return fmt.Errorf("開啟資料庫失敗: %v", err)
	}

	// 建立表格
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT,
		email TEXT,
		password TEXT
	);
	`
	_, err = DB.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("建立表格失敗: %v", err)
	}

	log.Println("✅ SQLite 初始化成功")
	return nil
}

func InsertUser(id, name, email, password string) error {
	_, err := DB.Exec("INSERT INTO users(id, name, email, password) VALUES (?, ?, ?, ?)", id, name, email, password)
	return err
}

func GetUserByID(id string) (string, string, string, error) {
	row := DB.QueryRow("SELECT name, email, password FROM users WHERE id = ?", id)
	
	var name, email, password string
	err := row.Scan(&name, &email, &password)
	return name, email, password, err
}
