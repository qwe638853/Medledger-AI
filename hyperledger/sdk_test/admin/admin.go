package admin

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sdk_test/database"
	fc "sdk_test/fabric"
	wl "sdk_test/wallet"
)

func main() {
	// ✅ 健檢中心帳號資訊（可從設定檔讀取，或 CLI 輸入）
	userId := "clinic001"
	password := "clinicpass"
	name := "百花健檢中心"
	date := "2025-05-13"
	email := "clinic001@example.com"
	phone := "049-1234567"

	// ✅ 檢查是否已存在
	exists, err := database.IsUserExists(userId)
	if err != nil {
		log.Fatalf("查詢資料庫失敗: %v", err)
	}
	if exists {
		log.Fatalf("此帳號已存在: %s", userId)
	}

	// ✅ Fabric CA 註冊
	err = fc.RegisterUser(
		"http://localhost:7054",
		"../orgs/org1.example.com/users/org1-admin/msp/signcerts/cert.pem",
		"../orgs/org1.example.com/users/org1-admin/msp/keystore/server.key",
		fc.RegisterRequest{
			ID:          userId,
			Secret:      password,
			Affiliation: "org1.department1",
			Type:        "client",
		})
	if err != nil {
		log.Fatalf("Fabric 註冊失敗: %v", err)
	}
	fmt.Println("✅ CA 註冊成功")

	// ✅ 產生 CSR & 金鑰
	privKey, csrPEM, err := fc.GenerateCSR(userId)
	if err != nil {
		log.Fatalf("產生 CSR 失敗: %v", err)
	}

	// ✅ 檔案儲存路徑
	// ✅ 建立使用者資料夾
	baseDir := filepath.Join("msp-data", "users", userId)
	os.MkdirAll(filepath.Join(baseDir, "keystore"), 0700)
	os.MkdirAll(filepath.Join(baseDir, "signcerts"), 0700)
	os.MkdirAll(filepath.Join(baseDir, "csr"), 0700)

	// ✅ 儲存 CSR
	csrPath := filepath.Join(baseDir, "csr", "csr.pem")
	err = fc.SaveCSRToFile(csrPEM, csrPath)
	if err != nil {
		log.Printf("❌ 寫入 CSR 失敗: %v", err)
		return

	}

	// ✅ 儲存私鑰
	keyPath := filepath.Join(baseDir, "keystore", "key.pem")
	err = fc.SavePrivateKeyToFile(privKey, keyPath)
	if err != nil {
		log.Printf("❌ 寫入私鑰失敗: %v", err)
		return
	}

	// ✅ Enroll
	err = fc.EnrollUser("http://localhost:7054", userId, password, fc.EnrollRequest{
		Certificate_request: string(csrPEM),
	})
	if err != nil {
		log.Fatalf("Enroll 失敗: %v", err)
	}
	fmt.Println("✅ Enroll 憑證成功")

	// ✅ 寫入 wallet
	w := wl.New()
	err = w.PutFile(userId, csrPath, keyPath, "Org1MSP")
	if err != nil {
		log.Fatalf("錢包寫入失敗: %v", err)
	}

	// ✅ 寫入 SQLite
	err = database.InsertUser(userId, password, name, date, email, phone)
	if err != nil {
		log.Fatalf("資料庫寫入失敗: %v", err)
	}

	fmt.Println("🎉 健檢中心帳號建立完成！")
}
