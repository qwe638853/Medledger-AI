package service

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"unicode"

	"go_server/database"
	fc "go_server/fabric"
	pb "go_server/proto"
	ut "go_server/utils"
	wl "go_server/wallet"

	"github.com/hyperledger/fabric-ca/api"
)

// HandleRegisterUser 處理用戶註冊邏輯 + 寫入 SQLite + Fabric CA 註冊
func HandleRegisterUser(ctx context.Context, req *pb.RegisterUserRequest, wallet wl.WalletInterface) (*pb.RegisterResponse, error) {
	log.Printf("收到用戶註冊請求: %v", req)

	// ✅ 基本欄位驗證
	if req.UserId == "" || req.Password == "" || req.Name == "" || req.Date == "" || req.Email == "" || req.Phone == "" {
		return &pb.RegisterResponse{Success: false, Message: "所有欄位皆為必填"}, nil
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return &pb.RegisterResponse{Success: false, Message: "Email 格式錯誤"}, nil
	}
	for _, c := range req.Phone {
		if !unicode.IsDigit(c) {
			return &pb.RegisterResponse{Success: false, Message: "電話號碼只能是數字"}, nil
		}
	}

	// ✅ SQLite 查重
	exists, err := database.IsUserExists(req.UserId)
	if err != nil {
		return &pb.RegisterResponse{Success: false, Message: "查詢使用者時出錯"}, nil
	}
	if exists {
		return &pb.RegisterResponse{Success: false, Message: "帳號已存在"}, nil
	}

	// ✅ 呼叫 Fabric CA 註冊帳號（使用 api.RegistrationRequest）
	err = fc.RegisterUser(
		"http://localhost:7054",
		"../orgs/org1.example.com/users/org1-admin/msp/signcerts/cert.pem",
		"../orgs/org1.example.com/users/org1-admin/msp/keystore/server.key",
		api.RegistrationRequest{
			Name:        req.UserId,
			Secret:      req.Password,
			Type:        "client",
			Affiliation: "org1.department1",
			Attributes: []api.Attribute{
				{Name: "role", Value: "patient", ECert: true},
			},
		},
	)
	if err != nil {
		log.Printf("❌ Fabric CA 註冊失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "Fabric 註冊失敗"}, nil
	}

	// ✅ 產生私鑰與 CSR
	privKey, csrPEM, err := fc.GenerateCSR(req.UserId)
	if err != nil {
		log.Printf("❌ 產生私鑰或 CSR 失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "無法產生憑證"}, nil
	}

	// ✅ 建立使用者資料夾並儲存檔案
	baseDir := filepath.Join("msp-data", "users", req.UserId)
	os.MkdirAll(filepath.Join(baseDir, "keystore"), 0700)
	os.MkdirAll(filepath.Join(baseDir, "signcerts"), 0700)
	os.MkdirAll(filepath.Join(baseDir, "csr"), 0700)

	csrPath := filepath.Join(baseDir, "csr", "csr.pem")
	err = fc.SaveCSRToFile(csrPEM, csrPath)
	if err != nil {
		log.Printf("❌ 寫入 CSR 失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "儲存 CSR 失敗"}, nil
	}

	keyPath := filepath.Join(baseDir, "keystore", "key.pem")
	err = fc.SavePrivateKeyToFile(privKey, keyPath)
	if err != nil {
		log.Printf("❌ 寫入私鑰失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "儲存私鑰失敗"}, nil
	}

	// ✅ Enroll 產生證書
	enrollReq := fc.EnrollRequest{
		Certificate_request: string(csrPEM),
	}

	certPem, err := fc.EnrollUser("http://localhost:7054", req.UserId, req.Password, enrollReq)
	if err != nil {
		log.Fatalf("Enroll 失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "Enroll 憑證註冊失敗"}, nil
	}

	certPath := filepath.Join(baseDir, "signcerts", "cert.pem")
	err = fc.SaveCertToFile(certPem, certPath)
	if err != nil {
		log.Printf("❌ 寫入證書失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "儲存證書失敗"}, nil
	}

	err = wallet.PutFile(req.UserId, certPath, keyPath, "Org1MSP")
	if err != nil {
		log.Printf("wallet save error: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "儲存錢包失敗"}, nil
	}

	// ✅ 寫入 SQLite
	err = database.InsertUser(req.UserId, req.Password, req.Name, req.Date, req.Email, req.Phone)
	if err != nil {
		log.Printf("❌ 寫入資料庫失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "寫入資料庫失敗"}, nil
	}

	return &pb.RegisterResponse{Success: true, Message: "用戶註冊成功"}, nil
}

// HandleRegisterInsurer 處理保險業者註冊邏輯 + 寫入 SQLite + Fabric CA 註冊
func HandleRegisterInsurer(ctx context.Context, req *pb.RegisterInsurerRequest, wallet wl.WalletInterface) (*pb.RegisterResponse, error) {
	log.Printf("收到保險業者註冊請求: %v", req)

	// ✅ 基本欄位驗證
	if req.InsurerId == "" || req.Password == "" || req.CompanyName == "" || req.ContactPerson == "" || req.Email == "" || req.Phone == "" {
		return &pb.RegisterResponse{Success: false, Message: "所有欄位皆為必填"}, nil
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return &pb.RegisterResponse{Success: false, Message: "Email 格式錯誤"}, nil
	}
	for _, c := range req.Phone {
		if !unicode.IsDigit(c) {
			return &pb.RegisterResponse{Success: false, Message: "電話號碼只能是數字"}, nil
		}
	}
	log.Printf("嘗試尋找保險業者ID: '%s'", req.InsurerId)

	// ✅ SQLite 查重
	exists, err := database.IsInsurerExists(req.InsurerId)
	if err != nil {
		return &pb.RegisterResponse{Success: false, Message: "查詢保險業者時出錯"}, nil
	}
	if exists {
		return &pb.RegisterResponse{Success: false, Message: "保險業者帳號已存在"}, nil
	}
	log.Printf("保險業者ID查詢結果: 存在=%v, 錯誤=%v", exists, err)
	// ✅ 呼叫 Fabric CA 註冊帳號（使用 api.RegistrationRequest）
	err = fc.RegisterUser(
		"http://localhost:7054",
		"../orgs/org1.example.com/users/org1-admin/msp/signcerts/cert.pem",
		"../orgs/org1.example.com/users/org1-admin/msp/keystore/server.key",
		api.RegistrationRequest{
			Name:        req.InsurerId,
			Secret:      req.Password,
			Type:        "client",
			Affiliation: "org1.department2",
			Attributes: []api.Attribute{
				{Name: "role", Value: "insurer", ECert: true},
			},
		},
	)
	if err != nil {
		log.Printf("❌ Fabric CA 註冊失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "Fabric 註冊失敗"}, nil
	}

	// ✅ 產生私鑰與 CSR
	privKey, csrPEM, err := fc.GenerateCSR(req.InsurerId)
	if err != nil {
		log.Printf("❌ 產生私鑰或 CSR 失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "無法產生憑證"}, nil
	}

	// ✅ 建立保險業者資料夾並儲存檔案
	baseDir := filepath.Join("msp-data", "insurers", req.InsurerId)
	os.MkdirAll(filepath.Join(baseDir, "keystore"), 0700)
	os.MkdirAll(filepath.Join(baseDir, "signcerts"), 0700)
	os.MkdirAll(filepath.Join(baseDir, "csr"), 0700)

	csrPath := filepath.Join(baseDir, "csr", "csr.pem")
	err = fc.SaveCSRToFile(csrPEM, csrPath)
	if err != nil {
		log.Printf("❌ 寫入 CSR 失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "儲存 CSR 失敗"}, nil
	}

	keyPath := filepath.Join(baseDir, "keystore", "key.pem")
	err = fc.SavePrivateKeyToFile(privKey, keyPath)
	if err != nil {
		log.Printf("❌ 寫入私鑰失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "儲存私鑰失敗"}, nil
	}

	// ✅ Enroll 產生證書
	enrollReq := fc.EnrollRequest{
		Certificate_request: string(csrPEM),
	}

	certPem, err := fc.EnrollUser("http://localhost:7054", req.InsurerId, req.Password, enrollReq)
	if err != nil {
		log.Fatalf("Enroll 失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "Enroll 憑證註冊失敗"}, nil
	}

	certPath := filepath.Join(baseDir, "signcerts", "cert.pem")
	err = fc.SaveCertToFile(certPem, certPath)
	if err != nil {
		log.Printf("❌ 寫入證書失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "儲存證書失敗"}, nil
	}

	err = wallet.PutFile(req.InsurerId, certPath, keyPath, "Org1MSP")
	if err != nil {
		log.Printf("wallet save error: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "儲存錢包失敗"}, nil
	}

	// ✅ 寫入 SQLite
	err = database.InsertInsurer(req.InsurerId, req.Password, req.CompanyName, req.ContactPerson, req.Email, req.Phone)
	if err != nil {
		log.Printf("❌ 寫入資料庫失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "寫入資料庫失敗"}, nil
	}
	log.Printf("保險業者原始ID: %s, 雜湊後ID (存入資料庫): %s", req.InsurerId, database.HashString(req.InsurerId))
	log.Printf("保險業者註冊成功: %s", req.InsurerId)

	return &pb.RegisterResponse{Success: true, Message: "保險業者註冊成功"}, nil
}

func HandleLogin(ctx context.Context, req *pb.LoginRequest, w wl.WalletInterface) (*pb.LoginResponse, error) {
	log.Printf("Received Login: %v", req)

	if req.UserId == "" || req.Password == "" {
		return &pb.LoginResponse{Success: false, Message: "帳號或密碼錯誤"}, nil
	}

	// 先檢查是否為保險業者
	insurerPw, err := database.GetInsurerPassword(req.UserId)
	if err != nil {
		log.Printf("查詢保險業者密碼錯誤: %v", err)
	}
	log.Printf("保險業者密碼查詢結果: 密碼=%s, 錯誤=%v", insurerPw, err)
	
	// 比對雜湊後的密碼
	hashedPassword := database.HashString(req.Password)
	if err == nil && insurerPw == hashedPassword {
		// 保險業者登入成功
		log.Printf("✅ 保險業者密碼驗證成功: %s", req.UserId)
		if !w.Exists(req.UserId) {
			log.Printf("❌ 保險業者錢包不存在: %s", req.UserId)
			return &pb.LoginResponse{Success: false, Message: "錢包不存在"}, nil
		}

		token, err := ut.GenerateJWT(req.UserId)
		if err != nil {
			return &pb.LoginResponse{Success: false, Message: "產生 token 失敗"}, nil
		}

		return &pb.LoginResponse{
			Success: true,
			Message: "保險業者登入成功",
			Token:   token,
		}, nil
	}

	// 再檢查是否為普通用戶帳號
	log.Printf("檢查是否為普通用戶帳號: %s", req.UserId)
	dbPw, err := database.GetUserPassword(req.UserId)
	if err != nil {
		log.Printf("查詢普通用戶密碼錯誤: %v", err)
	}
	log.Printf("普通用戶密碼查詢結果: 密碼=%s, 錯誤=%v", dbPw, err)
	
	// 比對雜湊後的密碼
	if err != nil || dbPw != hashedPassword {
		log.Printf("❌ 密碼驗證失敗: 用戶密碼=%s, 輸入密碼雜湊=%s", dbPw, hashedPassword)
		return &pb.LoginResponse{Success: false, Message: "帳號或密碼錯誤"}, nil
	}

	if !w.Exists(req.UserId) {
		log.Printf("❌ 錢包不存在: %s", req.UserId)
		return &pb.LoginResponse{Success: false, Message: "錢包不存在"}, nil
	}

	log.Printf("✅ 普通用戶密碼驗證成功: %s", req.UserId)
	token, err := ut.GenerateJWT(req.UserId)
	if err != nil {
		return &pb.LoginResponse{Success: false, Message: "產生 token 失敗"}, nil
	}

	return &pb.LoginResponse{
		Success: true,
		Message: "登入成功",
		Token:   token,
	}, nil
}
