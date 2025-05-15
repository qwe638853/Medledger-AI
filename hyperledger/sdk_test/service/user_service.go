package service

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"unicode"

	"sdk_test/database"
	fc "sdk_test/fabric"
	pb "sdk_test/proto"
	ut "sdk_test/utils"
	wl "sdk_test/wallet"

	"github.com/hyperledger/fabric-ca/api"
)

// HandleRegister 處理註冊邏輯 + 寫入 SQLite + Fabric CA 註冊
func HandleRegister(ctx context.Context, req *pb.RegisterRequest, wallet wl.WalletInterface) (*pb.RegisterResponse, error) {
	log.Printf("FFFFFReceived Register: %v", req)

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

	return &pb.RegisterResponse{Success: true, Message: "註冊成功"}, nil
}

func HandleLogin(ctx context.Context, req *pb.LoginRequest, w wl.WalletInterface) (*pb.LoginResponse, error) {
	log.Printf("Received Login: %v", req)

	if req.UserId == "" || req.Password == "" {
		return &pb.LoginResponse{Success: false, Message: "帳號或密碼錯誤"}, nil
	}

	dbPw, err := database.GetUserPassword(req.UserId)
	if err != nil {
		return &pb.LoginResponse{Success: false, Message: "查詢使用者時出錯"}, nil
	}
	if dbPw != req.Password {
		return &pb.LoginResponse{Success: false, Message: "帳號或密碼錯誤"}, nil
	}

	if !w.Exists(req.UserId) {
		log.Printf("❌ 錢包不存在: %s", req.UserId)
		return &pb.LoginResponse{Success: false, Message: "錢包不存在"}, nil
	}

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
