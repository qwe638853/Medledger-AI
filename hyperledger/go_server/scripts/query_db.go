package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go run query_db.go <命令>")
		fmt.Println("命令:")
		fmt.Println("  users     - 查詢所有用戶")
		fmt.Println("  insurers  - 查詢所有保險業者")
		fmt.Println("  count     - 統計數量")
		os.Exit(1)
	}

	dbPath := "database/user_data.sqlite"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Fatalf("❌ 資料庫文件不存在: %s", dbPath)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("❌ 開啟資料庫失敗: %v", err)
	}
	defer db.Close()

	command := os.Args[1]

	switch command {
	case "users":
		queryUsers(db)
	case "insurers":
		queryInsurers(db)
	case "count":
		countRecords(db)
	default:
		fmt.Printf("❌ 未知命令: %s\n", command)
		os.Exit(1)
	}
}

func queryUsers(db *sql.DB) {
	fmt.Println("==========================================")
	fmt.Println("📋 所有用戶資料")
	fmt.Println("==========================================")
	fmt.Printf("%-64s %-20s %-15s %-30s %-15s\n", "Username (雜湊)", "Name", "Date", "Email", "Phone")
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------")

	rows, err := db.Query("SELECT username, name, date, email, phone FROM users")
	if err != nil {
		log.Fatalf("❌ 查詢失敗: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var username, name, date, email, phone string
		if err := rows.Scan(&username, &name, &date, &email, &phone); err != nil {
			log.Printf("⚠️  讀取行失敗: %v", err)
			continue
		}
		fmt.Printf("%-64s %-20s %-15s %-30s %-15s\n", username, name, date, email, phone)
	}

	fmt.Println("")
	count := 0
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	fmt.Printf("總共 %d 位用戶\n", count)
}

func queryInsurers(db *sql.DB) {
	fmt.Println("==========================================")
	fmt.Println("📋 所有保險業者資料")
	fmt.Println("==========================================")
	fmt.Printf("%-64s %-30s %-20s %-30s %-15s\n", "Insurer ID (雜湊)", "Company Name", "Contact Person", "Email", "Phone")
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------")

	rows, err := db.Query("SELECT insurer_id, company_name, contact_person, email, phone FROM insurers")
	if err != nil {
		log.Fatalf("❌ 查詢失敗: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var insurerID, companyName, contactPerson, email, phone string
		if err := rows.Scan(&insurerID, &companyName, &contactPerson, &email, &phone); err != nil {
			log.Printf("⚠️  讀取行失敗: %v", err)
			continue
		}
		fmt.Printf("%-64s %-30s %-20s %-30s %-15s\n", insurerID, companyName, contactPerson, email, phone)
	}

	fmt.Println("")
	count := 0
	db.QueryRow("SELECT COUNT(*) FROM insurers").Scan(&count)
	fmt.Printf("總共 %d 家保險業者\n", count)
}

func countRecords(db *sql.DB) {
	var userCount, insurerCount int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	db.QueryRow("SELECT COUNT(*) FROM insurers").Scan(&insurerCount)

	fmt.Println("==========================================")
	fmt.Println("📊 資料庫統計")
	fmt.Println("==========================================")
	fmt.Printf("用戶數量:      %d\n", userCount)
	fmt.Printf("保險業者數量:  %d\n", insurerCount)
	fmt.Printf("總計:         %d\n", userCount+insurerCount)
}

