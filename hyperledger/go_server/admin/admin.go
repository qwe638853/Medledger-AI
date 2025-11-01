package main

import (
    "fmt"
    "log"
    "context"

    db "go_server/database"
    fc "go_server/fabric"
    vs "go_server/vaultstore"
    wl "go_server/wallet"
    sg "go_server/secure/signer"

    "crypto/x509"
    "crypto/rand"
    "crypto/x509/pkix"
    "encoding/pem"

    "github.com/hyperledger/fabric-ca/api"
    "github.com/joho/godotenv"
)

func main() {
    // 載入 .env（若存在）
    if err := godotenv.Load(); err == nil {
        log.Printf("[env] 已載入 .env")
    } else {
        log.Printf("[env] 無 .env 或載入失敗：%v", err)
    }
	// ✅ 健檢中心帳號資訊（可從設s定檔讀取，或 CLI 輸入）
	err := db.InitDB("database/user_data.sqlite")
	if err != nil {
		log.Fatalf("❌ SQLite 初始化失敗: %v", err)
	}
	userId := "clinic2"
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

    // ✅ TransitSigner 產生 CSR（私鑰不出庫）
    store, err := vs.NewFromEnv()
    if err != nil { log.Fatalf("Vault 初始化失敗: %v", err) }
    // 統一使用 hash 值建立 Transit key（與資料庫設計一致）
    clinicHash := db.HashString(userId)
    // 確保 Transit 簽章金鑰存在（用於 CSR 和 Fabric 交易簽章）
    if err := store.EnsureTransitKey(context.Background(), "clinic-"+clinicHash); err != nil {
        log.Fatalf("建立 Transit 簽章金鑰失敗: %v", err)
    }
    // 確保 Transit wrap 金鑰存在（用於資料金鑰包裝/解包，與註冊時一併建立）
    if err := store.EnsureTransitKeyOfType(context.Background(), "clinic-"+clinicHash+"-wrap", "aes256-gcm96"); err != nil {
        log.Fatalf("建立 Transit wrap 金鑰失敗: %v", err)
    }
    pub, err := store.TransitGetPublicKey(context.Background(), "clinic-"+clinicHash)
    if err != nil { log.Fatalf("讀取 Transit 公鑰失敗: %v", err) }
    signerObj, err := sg.NewTransitSignerWithPublicKey(store, "clinic-"+clinicHash, pub)
    if err != nil { log.Fatalf("建立 TransitSigner 失敗: %v", err) }
    tmpl := x509.CertificateRequest{ Subject: pkix.Name{ CommonName: userId } }
    csrDER, err := x509.CreateCertificateRequest(rand.Reader, &tmpl, signerObj)
    if err != nil { log.Fatalf("CSR 產生失敗: %v", err) }
    csrPEM := pem.EncodeToMemory(&pem.Block{Type:"CERTIFICATE REQUEST", Bytes: csrDER})
    keyPEM := []byte("")

	// ✅ Enroll（用自己產生的 CSR）
	enrollReq := fc.EnrollRequest{
		Certificate_request: string(csrPEM),
	}
    certPem, err := fc.EnrollUser("http://localhost:7054", userId, password, enrollReq)
	if err != nil {
		log.Fatalf("Enroll 失敗: %v", err)
	}

    // ✅ 寫入 Vault（clinic 類別）
    if store, verr := vs.NewFromEnv(); verr != nil {
        log.Printf("[Vault] 初始化失敗（略過）：%v", verr)
    } else {
        if werr := store.WriteClinicMaterial(context.Background(), userId, csrPEM, keyPEM, certPem); werr != nil {
            log.Printf("[Vault] 寫入 clinic 材料失敗：%v", werr)
        } else {
            log.Printf("[Vault] ✅ 已寫入 clinic 材料至 Vault：%s", userId)
        }
    }
	// ✅ 寫入 wallet
	w := wl.New()
    // DB 僅存引用
    // signerUri 使用診所 Transit key：clinic-<hash>（統一使用 hash）
    err = w.PutReference(userId, "Org1MSP", "transit://clinic-"+clinicHash, "kv://clinics/"+userId)
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
