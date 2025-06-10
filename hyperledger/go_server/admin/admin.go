package main

import (

	"fmt"
	"log"

	"os"
	"path/filepath"

	db "go_server/database"
	fc "go_server/fabric"
	wl "go_server/wallet"

	"github.com/hyperledger/fabric-ca/api"
)

func main() {
	// ✅ 健檢中心帳號資訊（可從設s定檔讀取，或 CLI 輸入）
	err := db.InitDB("database/user_data.sqlite")
	if err != nil {
		log.Fatalf("❌ SQLite 初始化失敗: %v", err)
	}
	userId := "clinic000001"
	password := "clinicpass"
	name := "健檢中心1"
	date := "2025-05-13"
	email := "clinic001@example.com"
	phone := "049-1234567"

	// ✅ 檢查是否已存在
	exists, err := db.IsUserExists(userId)
	if err != nil {
		log.Fatalf("查詢資料庫失敗: %v", err)
	}
	if exists {
		log.Fatalf("此帳號已存在: %s", userId)
	}

	// ✅ Fabric CA 註冊（使用官方 API 結構）
	err = fc.RegisterUser(
		"http://localhost:7054",
		"../orgs/org1.example.com/users/org1-admin/msp/signcerts/cert.pem",
		"../orgs/org1.example.com/users/org1-admin/msp/keystore/server.key",
		api.RegistrationRequest{
			Name:        userId,
			Secret:      password,
			Type:        "client",
			Affiliation: "org1.department1",
			Attributes: []api.Attribute{
				{Name: "role", Value: "clinic", ECert: true},
				{Name: "clinicId", Value: userId, ECert: true},
			},
		},
	)
	if err != nil {
		log.Fatalf("Fabric 註冊失敗: %v", err)
	}
	fmt.Println("✅ CA 註冊成功")

	// ✅ 產生 CSR & 金鑰
	privKey, csrPEM, err := fc.GenerateCSR(userId)
	if err != nil {
		log.Fatalf("產生 CSR 失敗: %v", err)
	}

	// ✅ 建立使用者資料夾並儲存檔案
	baseDir := filepath.Join("msp-data", "clinic", userId)
	os.MkdirAll(filepath.Join(baseDir, "keystore"), 0700)
	os.MkdirAll(filepath.Join(baseDir, "signcerts"), 0700)
	os.MkdirAll(filepath.Join(baseDir, "csr"), 0700)

	csrPath := filepath.Join(baseDir, "csr", "csr.pem")
	err = fc.SaveCSRToFile(csrPEM, csrPath)
	if err != nil {
		log.Printf("❌ 寫入 CSR 失敗: %v", err)
		return
	}

	keyPath := filepath.Join(baseDir, "keystore", "key.pem")
	err = fc.SavePrivateKeyToFile(privKey, keyPath)
	if err != nil {
		log.Printf("❌ 寫入私鑰失敗: %v", err)
		return
	}

	// ✅ Enroll（用自己產生的 CSR）
	enrollReq := fc.EnrollRequest{
		Certificate_request: string(csrPEM),
	}
	certPem, err := fc.EnrollUser("http://localhost:7054", userId, password, enrollReq)
	if err != nil {
		log.Fatalf("Enroll 失敗: %v", err)
	}

	certPath := filepath.Join(baseDir, "signcerts", "cert.pem")
	err = fc.SaveCertToFile(certPem, certPath)
	if err != nil {
		log.Printf("❌ 寫入證書失敗: %v", err)
		return
	}
	// ✅ 寫入 wallet
	w := wl.New()
	err = w.PutFile(userId, certPath, keyPath, "Org1MSP")
	if err != nil {
		log.Fatalf("錢包寫入失敗: %v", err)
	}

	// ✅ 寫入 SQLite
	err = db.InsertUser(userId, password, name, date, email, phone)
	if err != nil {
		log.Fatalf("資料庫寫入失敗: %v", err)
	}

	fmt.Println("🎉 健檢中心帳號建立完成！")
}
