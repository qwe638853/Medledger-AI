package service

import (
    "crypto/ecdsa"
    "crypto/x509"
    "encoding/pem"
    "encoding/json"
    "context"
    "log"
    "os"
    "regexp"
    "unicode"
    "crypto/rand"
    "crypto/x509/pkix"

    "go_server/database"
    fc "go_server/fabric"
    wrap "go_server/secure/wrap"
    vs "go_server/vaultstore"
    pb "go_server/proto"
    ut "go_server/utils"
    wl "go_server/wallet"
    sg "go_server/secure/signer"

    "github.com/hyperledger/fabric-ca/api"
)

/**
 * @notice 用戶註冊：建立 SQLite 帳號並向 Fabric CA 註冊與 Enroll
 * @dev 驗證請求 → 查重 → CA Register → 產生金鑰與 CSR → Enroll → 錢包寫入 → DB 寫入
 * @param ctx 請求上下文
 * @param req 用戶註冊請求
 * @param wallet 錢包介面，用於寫入身份
 * @return *pb.RegisterResponse 成功與訊息, error 內部錯誤
 */
// HandleRegisterUser 處理用戶註冊邏輯 + 寫入 SQLite + Fabric CA 註冊
func HandleRegisterUser(ctx context.Context, req *pb.RegisterUserRequest, wallet wl.WalletInterface, builder fc.GWBuilder) (*pb.RegisterResponse, error) {
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
	log.Printf("嘗試尋找用戶ID: '%s'", req.UserId)

	// ✅ SQLite 查重
	exists, err := database.IsUserExists(req.UserId)
	if err != nil {
		return &pb.RegisterResponse{Success: false, Message: "查詢使用者時出錯"}, nil
	}
	if exists {
		return &pb.RegisterResponse{Success: false, Message: "帳號已存在"}, nil
	}
	

	// ✅ 呼叫 Fabric CA 註冊帳號（使用 api.RegistrationRequest）
	log.Printf("[DEBUG] 開始 Fabric CA 註冊，用戶ID: %s", req.UserId)
	log.Printf("[DEBUG] Fabric CA URL: http://localhost:7054")
	log.Printf("[DEBUG] 管理員證書路徑: ../orgs/org1.example.com/users/org1-admin/msp/signcerts/cert.pem")
	log.Printf("[DEBUG] 管理員私鑰路徑: ../orgs/org1.example.com/users/org1-admin/msp/keystore/server.key")
	
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
		log.Printf("[DEBUG] 註冊失敗詳細信息 - 用戶: %s, 錯誤類型: %T", req.UserId, err)
		return &pb.RegisterResponse{Success: false, Message: "Fabric 註冊失敗"}, nil
	}
	log.Printf("[DEBUG] ✅ Fabric CA 註冊成功，用戶ID: %s", req.UserId)

    // ✅ 以 TransitSigner 產生 CSR（私鑰不出庫）
    log.Printf("[DEBUG] 開始 Transit 產生 CSR，用戶ID: %s", req.UserId)
    store, err := vs.NewFromEnv()
    if err != nil { return &pb.RegisterResponse{Success:false, Message:"Vault 初始化失敗"}, nil }
    // 確保 Transit key 存在
    if err := store.EnsureTransitKey(ctx, "user-"+req.UserId); err != nil {
        return &pb.RegisterResponse{Success:false, Message:"建立使用者 Transit 金鑰失敗"}, nil
    }
    pub, err := store.TransitGetPublicKey(ctx, "user-"+req.UserId)
    if err != nil { return &pb.RegisterResponse{Success:false, Message:"取得公鑰失敗"}, nil }
    signerObj, err := sg.NewTransitSignerWithPublicKey(store, "user-"+req.UserId, pub)
    if err != nil { return &pb.RegisterResponse{Success:false, Message:"建立 Transit 簽章器失敗"}, nil }
    tmpl := x509.CertificateRequest{ Subject: pkix.Name{ CommonName: req.UserId } }
    csrDER, err := x509.CreateCertificateRequest(rand.Reader, &tmpl, signerObj)
    if err != nil { return &pb.RegisterResponse{Success:false, Message:"CSR 產生失敗"}, nil }
    csrPEM := pem.EncodeToMemory(&pem.Block{Type:"CERTIFICATE REQUEST", Bytes: csrDER})
    keyPEM := []byte("")

	// ✅ Enroll 產生證書
	enrollReq := fc.EnrollRequest{
		Certificate_request: string(csrPEM),
	}

    certPem, err := fc.EnrollUser("http://localhost:7054", req.UserId, req.Password, enrollReq)
	if err != nil {
		log.Fatalf("Enroll 失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "Enroll 憑證註冊失敗"}, nil
	}
    // ✅ 寫入 Vault + 錢包改為僅存引用（DB 不存證書與私鑰）
    if store, verr := vs.NewFromEnv(); verr != nil {
        log.Printf("[Vault] 初始化失敗（略過）：%v", verr)
    } else {
        if werr := store.WriteUserMaterial(ctx, req.UserId, csrPEM, keyPEM, certPem); werr != nil {
            log.Printf("[Vault] 寫入使用者材料失敗：%v", werr)
        } else {
            log.Printf("[Vault] ✅ 已寫入使用者材料至 Vault：%s", req.UserId)
        }
    }
    // signerUri 使用使用者專屬 Transit key：user-<id>
    err = wallet.PutReference(req.UserId, "Org1MSP", "transit://user-"+req.UserId, "kv://users/"+req.UserId)
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

    // 背景回填：為使用者既有報告包一份 patient 的 wrapped key（若存在報告）
    go func(userID string) {
        defer func() { recover() }()
        // 取得使用者與平台身份（從 Vault 補齊）
        userEntry, okU := wallet.GetResolved(userID)
        platformEntry, okP := wallet.GetResolved("platform")
        if !okU || !okP || userEntry == nil || platformEntry == nil || userEntry.Cert == nil {
            log.Printf("[Backfill] 缺少必要身份或金鑰，略過包鍵 user=%v platform=%v", okU, okP)
            return
        }

        // 確認憑證公鑰型別（不使用變數）
        if _, ok := userEntry.Cert.PublicKey.(*ecdsa.PublicKey); !ok {
            log.Printf("[Backfill] 用戶公鑰不是 ECDSA，略過 user=%s", userID)
            return
        }

        // 建立以用戶身分的合約
        contract, gw, err := builder.NewContract(userEntry.ID, userEntry.Signer)
        if err != nil { log.Printf("[Backfill] NewContract 失敗: %v", err); return }
        defer gw.Close()

        // 查詢 meta 列表
        metasRaw, err := contract.EvaluateTransaction("ListMyReportMeta")
        if err != nil { log.Printf("[Backfill] 查 meta 失敗: %v", err); return }
        var metas []struct {
            ReportID  string `json:"reportId"`
            ClinicID  string `json:"clinicId"`
            CreatedAt int64  `json:"createdAt"`
        }
        if err := json.Unmarshal(metasRaw, &metas); err != nil {
            log.Printf("[Backfill] 解析 meta 失敗: %v", err); return
        }
        if len(metas) == 0 { return }

        for _, m := range metas {
            // 讀取 envelope JSON
            envRaw, err := contract.EvaluateTransaction("ReadMyReport", m.ReportID)
            if err != nil { log.Printf("[Backfill] 讀取報告失敗 id=%s err=%v", m.ReportID, err); continue }

            // 若已存在 patient key 則略過
            var env struct{ WrappedKeys map[string]any `json:"wrappedKeys"` }
            if err := json.Unmarshal(envRaw, &env); err != nil { log.Printf("[Backfill] 解析 envelope 失敗: %v", err); continue }
            if env.WrappedKeys != nil {
                if _, exists := env.WrappedKeys["patient"]; exists { continue }
            }

    // Transit 模組：平台用 Transit decrypt 解出 dataKey，再用 Transit encrypt 為用戶包一份
    tw, err := wrap.NewTransitWrapperFromEnv(); if err != nil { log.Printf("[Backfill] Vault 初始化失敗: %v", err); return }
    updated, err := tw.AddRecipientTransit(context.Background(), envRaw, "patient", "user-"+userID)
    if err != nil { log.Printf("[Backfill] AddRecipientTransit 失敗: %v", err); return }
            if err != nil { log.Printf("[Backfill] AddRecipient 失敗 id=%s err=%v", m.ReportID, err); continue }

            // 以用戶身分提交鏈上更新
            if _, err := contract.SubmitTransaction("UpdateReport", m.ReportID, string(updated)); err != nil {
                log.Printf("[Backfill] UpdateReport 失敗 id=%s err=%v", m.ReportID, err)
                continue
            }
            log.Printf("[Backfill] 已為報告 %s 加入 patient wrapped key", m.ReportID)
        }
    }(req.UserId)

    return &pb.RegisterResponse{Success: true, Message: "用戶註冊成功"}, nil
}

/**
 * @notice 保險業者註冊：建立 SQLite 帳號並向 Fabric CA 註冊與 Enroll
 * @dev 驗證請求 → 查重 → CA Register → 產生金鑰與 CSR → Enroll → 錢包寫入 → DB 寫入
 * @param ctx 請求上下文
 * @param req 保險業者註冊請求
 * @param wallet 錢包介面，用於寫入身份
 * @return *pb.RegisterResponse 成功與訊息, error 內部錯誤
 */
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

    // ✅ 以 TransitSigner 產生 CSR（私鑰不出庫）
    store, err := vs.NewFromEnv()
    if err != nil { return &pb.RegisterResponse{Success:false, Message:"Vault 初始化失敗"}, nil }
    if err := store.EnsureTransitKey(ctx, "insurer-"+req.InsurerId); err != nil {
        return &pb.RegisterResponse{Success:false, Message:"建立保險業者 Transit 金鑰失敗"}, nil
    }
    pub, err := store.TransitGetPublicKey(ctx, "insurer-"+req.InsurerId)
    if err != nil { return &pb.RegisterResponse{Success:false, Message:"取得公鑰失敗"}, nil }
    signerObj, err := sg.NewTransitSignerWithPublicKey(store, "insurer-"+req.InsurerId, pub)
    if err != nil { return &pb.RegisterResponse{Success:false, Message:"建立 Transit 簽章器失敗"}, nil }
    tmpl := x509.CertificateRequest{ Subject: pkix.Name{ CommonName: req.InsurerId } }
    csrDER, err := x509.CreateCertificateRequest(rand.Reader, &tmpl, signerObj)
    if err != nil { return &pb.RegisterResponse{Success:false, Message:"CSR 產生失敗"}, nil }
    csrPEM := pem.EncodeToMemory(&pem.Block{Type:"CERTIFICATE REQUEST", Bytes: csrDER})
    keyPEM := []byte("")

	// ✅ Enroll 產生證書
	enrollReq := fc.EnrollRequest{
		Certificate_request: string(csrPEM),
	}

	certPem, err := fc.EnrollUser("http://localhost:7054", req.InsurerId, req.Password, enrollReq)
	if err != nil {
		log.Fatalf("Enroll 失敗: %v", err)
		return &pb.RegisterResponse{Success: false, Message: "Enroll 憑證註冊失敗"}, nil
	}

    // ✅ 寫入 Vault + 錢包改為僅存引用（DB 不存證書與私鑰）
    if store, verr := vs.NewFromEnv(); verr != nil {
        log.Printf("[Vault] 初始化失敗（略過）：%v", verr)
    } else {
        if werr := store.WriteInsurerMaterial(ctx, req.InsurerId, csrPEM, keyPEM, certPem); werr != nil {
            log.Printf("[Vault] 寫入保險業者材料失敗：%v", werr)
        } else {
            log.Printf("[Vault] ✅ 已寫入保險業者材料至 Vault：%s", req.InsurerId)
        }
    }
    err = wallet.PutReference(req.InsurerId, "Org1MSP", "transit://insurer-"+req.InsurerId, "kv://insurers/"+req.InsurerId)
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

/**
 * @notice 登入：驗證用戶或保險業者身份並簽發 JWT
 * @dev 依序檢查保險業者與一般用戶的雜湊密碼，比對成功後確認錢包存在並簽發 JWT
 * @param ctx 請求上下文
 * @param req 登入請求（user_id、password）
 * @param w 錢包介面，用於檢查身份是否存在
 * @return *pb.LoginResponse 結果與 JWT, error 內部錯誤
 */
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

// RegisterPlatformIdentity 幫「平台」向 Fabric CA 註冊並 Enroll，並寫入到錢包與 msp-data/platform
// - 使用者名稱固定為 "platform"
// - 密碼優先讀取環境變數 PLATFORM_CA_SECRET，預設為 "platformpw"
// - MSP 名稱預設 "Org1MSP"，如需更動可調整此函式
func RegisterPlatformIdentity(ctx context.Context, wallet wl.WalletInterface) error {
	platformID := "platform"
	platformSecret := os.Getenv("PLATFORM_CA_SECRET")
	if platformSecret == "" {
		platformSecret = "platformpw"
	}

	// 若錢包已存在，視為已完成，直接返回
	if wallet.Exists(platformID) {
		log.Printf("[Platform] 身分已存在於錢包，略過註冊/Enroll: %s", platformID)
		return nil
	}

	// 1) CA 註冊（以 org admin 身分）
	log.Printf("[Platform] 開始 Fabric CA 註冊: %s", platformID)
	if err := fc.RegisterUser(
		"http://localhost:7054",
		"../orgs/org1.example.com/users/org1-admin/msp/signcerts/cert.pem",
		"../orgs/org1.example.com/users/org1-admin/msp/keystore/server.key",
		api.RegistrationRequest{
			Name:        platformID,
			Secret:      platformSecret,
			Type:        "client",
			Affiliation: "org1.department1",
			Attributes: []api.Attribute{{Name: "role", Value: "platform", ECert: true}},
		},
	); err != nil {
		log.Printf("[Platform] ❌ CA 註冊失敗: %v", err)
		return err
	}

    // 2) 以 TransitSigner 產生 CSR（私鑰不出庫）
    store, err := vs.NewFromEnv()
    if err != nil { log.Printf("[Platform] ❌ Vault 初始化失敗: %v", err); return err }
    if err := store.EnsureTransitKey(ctx, "platform"); err != nil { log.Printf("[Platform] ❌ 建立 Transit 金鑰失敗: %v", err); return err }
    pub, err := store.TransitGetPublicKey(ctx, "platform")
    if err != nil { log.Printf("[Platform] ❌ 讀取 Transit 公鑰失敗: %v", err); return err }
    signerObj, err := sg.NewTransitSignerWithPublicKey(store, "platform", pub)
    if err != nil { log.Printf("[Platform] ❌ 建立 TransitSigner 失敗: %v", err); return err }
    tmpl := x509.CertificateRequest{ Subject: pkix.Name{ CommonName: platformID } }
    csrDER, err := x509.CreateCertificateRequest(rand.Reader, &tmpl, signerObj)
    if err != nil { log.Printf("[Platform] ❌ CSR 產生失敗: %v", err); return err }
    csrPEM := pem.EncodeToMemory(&pem.Block{Type:"CERTIFICATE REQUEST", Bytes: csrDER})
    keyPEM := []byte("")

	// 4) Enroll 取得憑證
    certPem, err := fc.EnrollUser("http://localhost:7054", platformID, platformSecret, fc.EnrollRequest{Certificate_request: string(csrPEM)})
	if err != nil {
		log.Printf("[Platform] ❌ Enroll 憑證失敗: %v", err)
		return err
	}
    // 寫入 Vault（platform 固定節點）
    if store, verr := vs.NewFromEnv(); verr != nil {
        log.Printf("[Vault] 初始化失敗（略過）：%v", verr)
    } else {
        if werr := store.WritePlatformMaterial(ctx, csrPEM, keyPEM, certPem); werr != nil {
            log.Printf("[Vault] 寫入 platform 材料失敗：%v", werr)
        } else {
            log.Printf("[Vault] ✅ 已寫入 platform 材料至 Vault")
        }
    }

	// 5) 寫入錢包（label 固定為 platform）
    if err := wallet.PutReference(platformID, "Org1MSP", "transit://platform", "kv://platform"); err != nil {
		log.Printf("[Platform] ❌ 寫入錢包失敗: %v", err)
		return err
	}

	log.Printf("[Platform] ✅ 註冊與 Enroll 成功，已寫入錢包: %s", platformID)
	return nil
}
