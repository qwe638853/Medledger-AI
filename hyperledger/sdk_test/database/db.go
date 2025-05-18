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

	// 授權請求表
	createAccessRequestStmt := `
	CREATE TABLE IF NOT EXISTS access_requests (
		request_id TEXT PRIMARY KEY,
		report_id TEXT,
		patient_id TEXT,
		requester_id TEXT,
		reason TEXT,
		requested_at INTEGER,
		expiry INTEGER,
		status TEXT
	);`

	_, err = DB.Exec(createAccessRequestStmt)
	if err != nil {
		return fmt.Errorf("建立授權請求資料表失敗: %v", err)
	}

	log.Println("✅ SQLite 初始化成功")
	return nil
}

// 查詢用戶帳號是否存在
func IsUserExists(username string) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	return count > 0, err
}

// 查詢保險業者帳號是否存在
func IsInsurerExists(insurerId string) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM insurers WHERE insurer_id = ?", insurerId).Scan(&count)
	return count > 0, err
}

// 新增用戶
func InsertUser(username, password, name, date, email, phone string) error {
	_, err := DB.Exec("INSERT INTO users(username, password, name, date, email, phone) VALUES (?, ?, ?, ?, ?, ?)",
		username, password, name, date, email, phone)
	return err
}

// 新增保險業者
func InsertInsurer(insurerId, password, companyName, contactPerson, email, phone string) error {
	_, err := DB.Exec("INSERT INTO insurers(insurer_id, password, company_name, contact_person, email, phone) VALUES (?, ?, ?, ?, ?, ?)",
		insurerId, password, companyName, contactPerson, email, phone)
	return err
}

// 取得用戶密碼
func GetUserPassword(username string) (string, error) {
	var password string
	err := DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

// 取得保險業者密碼
func GetInsurerPassword(insurerId string) (string, error) {
	var password string
	err := DB.QueryRow("SELECT password FROM insurers WHERE insurer_id = ?", insurerId).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

// 儲存授權請求
func InsertAccessRequest(requestId, reportId, patientId, requesterId, reason string, requestedAt, expiry int64, status string) error {
	_, err := DB.Exec(`
	INSERT INTO access_requests(request_id, report_id, patient_id, requester_id, reason, requested_at, expiry, status) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		requestId, reportId, patientId, requesterId, reason, requestedAt, expiry, status)
	return err
}

// 查詢用戶的授權請求
func GetAccessRequestsForUser(patientId string) ([]map[string]interface{}, error) {
	rows, err := DB.Query(`
	SELECT request_id, report_id, requester_id, reason, requested_at, expiry, status 
	FROM access_requests 
	WHERE patient_id = ?`, patientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var results []map[string]interface{}
	for rows.Next() {
		var requestId, reportId, requesterId, reason, status string
		var requestedAt, expiry int64
		
		if err := rows.Scan(&requestId, &reportId, &requesterId, &reason, &requestedAt, &expiry, &status); err != nil {
			return nil, err
		}
		
		result := map[string]interface{}{
			"request_id": requestId,
			"report_id": reportId,
			"requester_id": requesterId,
			"reason": reason,
			"requested_at": requestedAt,
			"expiry": expiry,
			"status": status,
		}
		results = append(results, result)
	}
	
	return results, nil
}

// 更新授權請求狀態
func UpdateAccessRequestStatus(requestId, status string) error {
	_, err := DB.Exec("UPDATE access_requests SET status = ? WHERE request_id = ?", status, requestId)
	return err
}

// 取得授權請求詳情
func GetAccessRequestById(requestId string) (map[string]interface{}, error) {
	var reportId, patientId, requesterId, reason, status string
	var requestedAt, expiry int64
	
	err := DB.QueryRow(`
	SELECT report_id, patient_id, requester_id, reason, requested_at, expiry, status 
	FROM access_requests WHERE request_id = ?`, requestId).Scan(
		&reportId, &patientId, &requesterId, &reason, &requestedAt, &expiry, &status)
	
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"request_id": requestId,
		"report_id": reportId,
		"patient_id": patientId,
		"requester_id": requesterId,
		"reason": reason,
		"requested_at": requestedAt,
		"expiry": expiry,
		"status": status,
	}, nil
}

// 獲取保險業者的待處理請求數量
func GetPendingRequestsCountForInsurer(insurerId string) (int, error) {
	var count int
	err := DB.QueryRow(`
	SELECT COUNT(*) FROM access_requests 
	WHERE requester_id = ? AND status = 'PENDING'`, insurerId).Scan(&count)
	
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// 獲取保險業者的已授權病患數量
func GetAuthorizedPatientsCountForInsurer(insurerId string) (int, error) {
	var count int
	err := DB.QueryRow(`
	SELECT COUNT(DISTINCT patient_id) FROM access_requests 
	WHERE requester_id = ? AND status = 'APPROVED'`, insurerId).Scan(&count)
	
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// 獲取保險業者的已授權報告
func GetAuthorizedReportsForInsurer(insurerId string) ([]map[string]interface{}, error) {
	rows, err := DB.Query(`
	SELECT ar.report_id, ar.patient_id, 'Report content will be fetched from blockchain' as content, 
	       strftime('%Y-%m-%d', datetime(ar.requested_at, 'unixepoch')) as date, 
	       strftime('%Y-%m-%d', datetime(ar.expiry, 'unixepoch')) as expiry
	FROM access_requests ar 
	WHERE ar.requester_id = ? AND ar.status = 'APPROVED'`, insurerId)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var reports []map[string]interface{}
	for rows.Next() {
		var reportId, patientId, content, date, expiry string
		
		if err := rows.Scan(&reportId, &patientId, &content, &date, &expiry); err != nil {
			return nil, err
		}
		
		report := map[string]interface{}{
			"report_id": reportId,
			"patient_id": patientId,
			"content": content,
			"date": date,
			"expiry": expiry,
		}
		reports = append(reports, report)
	}
	
	return reports, nil
}
