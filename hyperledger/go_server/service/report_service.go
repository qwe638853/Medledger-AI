package service

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "log"
    "strconv"
    "time"

    "go_server/database"
    fc "go_server/fabric"
    sw "go_server/secure/wrap"
    pb "go_server/proto"
    ut "go_server/utils"
    wl "go_server/wallet"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/types/known/emptypb"
)

/**
 * @notice 上傳健檢報告：驗證請求並提交鏈上交易
 * @dev 由 JWT 提取使用者 → 建立 Gateway → Submit UploadReport 交易
 * @param ctx 請求上下文
 * @param req 上傳報告請求（reportId, userId, testResultsJson）
 * @param wallet 錢包介面
 * @param builder Gateway 建構器
 * @return *pb.UploadReportResponse 結果, error 內部錯誤
 */
// HandleUploadReport 驗證請求 → 存 SQLite → 調用 Fabric
func HandleUploadReport(
	ctx context.Context,
	req *pb.UploadReportRequest,
	wallet wl.WalletInterface, builder fc.GWBuilder) (*pb.UploadReportResponse, error) {

	// 取得JWT 裡面的userID
	userID, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("[Debug] UploadReport userID=%s", userID)
    entry, ok := wallet.GetResolved(userID)
    if !ok {
        log.Printf("[UploadReport] GetResolved 失敗 userID=%s (可能缺 certUri/signerUri 或 Vault 授權/連線問題)", userID)
        return nil, status.Error(codes.PermissionDenied, "錢包不存在")
    }

    log.Printf("[Debug] UploadReport args: reportID=%s, patientHash=%s, dataLen=%d",
		req.ReportId, req.UserId, len(req.TestResultsJson))

	// 依使用者身分建立 Gateway + Contract
	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, err
	}
	defer gw.Close()

	sum := sha256.Sum256([]byte(req.UserId))
	hashedUserID := hex.EncodeToString(sum[:])
	log.Printf("[Debug] 查詢患者雜湊: %s", hashedUserID)

    // Transit：改為 clinic-only 加密，平台不再持有解包能力
    tw, err := sw.NewTransitWrapperFromEnv(); if err != nil { return nil, status.Error(codes.Internal, "Vault 初始化失敗") }
    // 統一使用 hash 值（與資料庫和註冊邏輯一致）
    clinicHash := database.HashString(userID)
    // label 固定為 "clinic"，實際身分由 baseKey（clinic-<hash>-wrap）決定
    envJSON, err := tw.EncryptReportTransitClinicOnly(ctx, []byte(req.TestResultsJson), "clinic", "clinic-"+clinicHash)
    if err != nil { return nil, status.Error(codes.Internal, "資料加密失敗") }
    // 上傳階段預先授權病患：僅在病患已註冊時才追加 patient wrapped key
    // 檢查方式：確認用戶在資料庫中已註冊（避免自動建立 wrap key）
    userExists, checkErr := database.IsUserExists(req.UserId)
    if checkErr == nil && userExists {
        // 用戶已註冊，嘗試添加 patient wrapped key（使用 hash）
        patientHash := database.HashString(req.UserId)
        if envJSON2, aerr := tw.AddRecipientTransitFrom(ctx, envJSON, "clinic", "clinic-"+clinicHash, "patient", "user-"+patientHash); aerr != nil {
            log.Printf("[Warn] 追加病患 wrapped key 失敗（Vault 權限不足或 wrap key 不存在）userId=%s err=%v — 將以 clinic-only 續傳", req.UserId, aerr)
        } else {
            envJSON = envJSON2
            log.Printf("[Debug] 已追加病患 wrapped key: user-%s-wrap", patientHash)
        } 
    } else {
        log.Printf("[Info] 病患尚未註冊 userId=%s，跳過追加 wrapped key（將以 clinic-only 續傳）", req.UserId)
    }
    // 呼叫鏈碼
    log.Printf("[Debug] 準備調用 SubmitTransaction: UploadReport (encrypted)")
    log.Printf("[Debug] 參數 - ReportID: %s, PatientHash: %s, EncryptedSize: %d bytes", 
        req.ReportId, hashedUserID, len(envJSON))

    result, err := contract.SubmitTransaction(
        "UploadReport",
        req.ReportId,
        hashedUserID,
        string(envJSON),
    )
	
	if err != nil {
		log.Printf("[Error] SubmitTransaction 失敗: %v", err)
		fc.PrintGatewayError(err) // 看錯誤細節
		return nil, status.Error(codes.Internal, "鏈上交易失敗")
	}
	
	log.Printf("[Debug] SubmitTransaction 成功完成, 結果: %s", string(result))

	return &pb.UploadReportResponse{
		Success: true, Message: "上傳成功",
	}, nil
}

// BackfillPatientAccessAfterRegister
// 說明：用於「病患註冊完成後」背景補授權，為該病患的歷史報告追加 patient wrapped key
// 邏輯：以病患身份查詢自己的報告 → 對沒有 patient wrapped key 的報告，
//      以診所 unwrap 來源（clinic-<hash>-wrap）重包到 user-<hash>-wrap → 提交 UpdateReport
// 注意：統一使用 hash 值建立 Transit key（與註冊邏輯一致）
func BackfillPatientAccessAfterRegister(ctx context.Context, userID string, wallet wl.WalletInterface, builder fc.GWBuilder) {
    defer func() { recover() }()
    entry, ok := wallet.GetResolved(userID)
    if !ok || entry == nil || entry.Cert == nil {
        log.Printf("[Backfill] 缺少必要身份或金鑰，略過包鍵 user=%v", ok)
        return
    }
    // 以用戶身分建立合約
    contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
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

        // 以診所 unwrap 來源重包 patient
        tw, err := sw.NewTransitWrapperFromEnv(); if err != nil { log.Printf("[Backfill] Vault 初始化失敗: %v", err); return }
        // 統一使用 hash 值（與註冊邏輯一致）
        clinicHash := database.HashString(m.ClinicID)
        patientHash := database.HashString(userID)
        updated, err := tw.AddRecipientTransitFrom(context.Background(), envRaw, "clinic", "clinic-"+clinicHash, "patient", "user-"+patientHash)
        if err != nil { log.Printf("[Backfill] AddRecipientTransitFrom 失敗: %v", err); continue }

        // 提交鏈上更新
        if _, err := contract.SubmitTransaction("UpdateReport", m.ReportID, string(updated)); err != nil {
            log.Printf("[Backfill] UpdateReport 失敗 id=%s err=%v", m.ReportID, err)
            continue
        }
        log.Printf("[Backfill] 已為報告 %s 加入 patient wrapped key", m.ReportID)
    }
}


// HandleRequestAccess 處理保險業者請求授權
/**
 * @notice 授權請求：保險業者向病患申請訪問報告
 * @dev 檢查保險業者身份 → 建立 Gateway → Submit RequestAccess 交易
 * @param ctx 請求上下文
 * @param req 授權請求（reportId, patientId, reason, expiry）
 * @param wallet 錢包介面
 * @param builder Gateway 建構器
 * @return *pb.RequestAccessResponse 結果, error 內部錯誤
 */
func HandleRequestAccess(
	ctx context.Context,
	req *pb.RequestAccessRequest,
	wallet wl.WalletInterface, 
	builder fc.GWBuilder) (*pb.RequestAccessResponse, error) {
	log.Printf("[Debug] HandleRequestAccess", req)
	// 取得JWT中的使用者ID（應為保險業者）
	requesterId, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析授權資訊")
	}

	// 檢查是否為有效的保險業者
	_, err = database.GetInsurerPassword(requesterId)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "只有保險業者可以申請授權")
	}

	// 檢查請求內容
	if req.ReportId == "" || req.PatientId == "" || req.Reason == "" {
		return nil, status.Error(codes.InvalidArgument, "必須提供報告ID、病患ID和申請原因")
	}

	// 設定過期時間，若未提供則預設30天
	expiry := req.Expiry
	if expiry == 0 {
		expiry = time.Now().Unix() + 30*24*60*60 // 30天
	}

	// 取得合約
    entry, ok := wallet.GetResolved(requesterId)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, err
	}
	defer gw.Close()

	// 計算病患雜湊
	sum := sha256.Sum256([]byte(req.PatientId))
	patientHash := hex.EncodeToString(sum[:])

	// 呼叫鏈碼
	_, err = contract.SubmitTransaction(
		"RequestAccess",
		req.ReportId,
		patientHash,
		req.Reason,
		strconv.FormatInt(expiry, 10),
	)
	if err != nil {
		fc.PrintGatewayError(err)
		return nil, status.Error(codes.Internal, "授權請求失敗")
	}

	return &pb.RequestAccessResponse{
		Success: true,
		Message: "授權請求已送出",
	}, nil
}

type rawAccessRequest struct {
	DocType      string `json:"docType"`
	RequestID    string `json:"requestId"`
	ReportID     string `json:"reportId"`
	PatientHash  string `json:"patientHash"`
	RequesterHash string `json:"requesterHash"`  // 統一使用 hash（用於 Transit key 命名）
	Reason       string `json:"reason"`
	RequestedAt  int64  `json:"requestedAt"`
	Expiry       int64  `json:"expiry"`
	Status       string `json:"status"`
}
// HandleListAccessRequests 列出病患的所有授權請求
/**
 * @notice 列出病患的所有待處理授權請求
 * @dev 病患身份 → Evaluate ListPendingAccessRequests
 * @param ctx 請求上下文
 * @param _ 空請求
 * @param wallet 錢包介面
 * @param builder Gateway 建構器
 * @return *pb.ListAccessRequestsResponse 清單, error 內部錯誤
 */
func HandleListAccessRequests(
	ctx context.Context,
	_ *emptypb.Empty,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (*pb.ListAccessRequestsResponse, error) {

	// 取得JWT中的使用者ID
	userID, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析授權資訊")
	}

    entry, ok := wallet.GetResolved(userID)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, err
	}
	defer gw.Close()

	// 呼叫鏈碼
	result, err := contract.EvaluateTransaction("ListPendingAccessRequests")
	if err != nil {
		fc.PrintGatewayError(err)
		return nil, status.Error(codes.Internal, "查詢失敗")
	}
	log.Printf("[Debug] 查詢到授權請求: %s", string(result))

	
	var raws []rawAccessRequest
	if err := json.Unmarshal(result, &raws); err != nil {
		return nil, status.Error(codes.Internal, "解析結果失敗")
	}

	var requests []*pb.AccessRequest
	for _, r := range raws {
		// 從資料庫獲取保險業者資訊
		// 注意：這裡的 RequesterHash 已經是雜湊值，直接使用
		log.Printf("[Debug] 嘗試使用鏈碼返回的 RequesterHash 查詢資料庫: %s", r.RequesterHash)
		insurer, err := database.GetInsurerByHash(r.RequesterHash)
		if err != nil {
			log.Printf("[Warning] 無法獲取保險業者資訊: %v", err)
			// 繼續處理其他請求
			continue
		}

		requests = append(requests, &pb.AccessRequest{
			RequestId:     r.RequestID,
			ReportId:      r.ReportID,
			PatientHash:   r.PatientHash,
			RequesterHash: r.RequesterHash,
			RequesterName: insurer.Name,
			CompanyName:   insurer.CompanyName,
			Reason:        r.Reason,
			RequestedAt:   r.RequestedAt,
			Expiry:        r.Expiry,
			Status:        r.Status,
		})
	}

	log.Printf("[Debug] 查詢到授權請求: %v", requests)

	return &pb.ListAccessRequestsResponse{
		Requests: requests,
	}, nil
}

// HandleApproveAccessRequest 處理授權請求的批准
/**
 * @notice 病患批准授權請求
 * @dev 病患身份 → Submit ApproveAndAuthorizeAccess
 * @param ctx 請求上下文
 * @param req 批准請求（requestId）
 * @param wallet 錢包介面
 * @param builder Gateway 建構器
 * @return *pb.ApproveAccessRequestResponse 結果, error 內部錯誤
 */
func HandleApproveAccessRequest(
	ctx context.Context,
	req *pb.ApproveAccessRequestRequest,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (*pb.ApproveAccessRequestResponse, error) {

	userID, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析授權資訊")
	}

    entry, ok := wallet.GetResolved(userID)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, err
	}
	defer gw.Close()

    // 在批准前先取得該 request 的必要資訊（reportId、requesterHash），避免批准後從 PENDING 列表取不到
    var pendingReq rawAccessRequest
    {
        result, qErr := contract.EvaluateTransaction("ListPendingAccessRequests")
        if qErr != nil {
            log.Printf("[Warn] 取得待處理授權請求失敗，將直接批准 requestId=%s err=%v", req.RequestId, qErr)
        } else {
            var raws []rawAccessRequest
            if uErr := json.Unmarshal(result, &raws); uErr != nil {
                log.Printf("[Warn] 解析待處理清單失敗，將直接批准 requestId=%s err=%v", req.RequestId, uErr)
            } else {
                for _, r := range raws {
                    if r.RequestID == req.RequestId {
                        pendingReq = r
                        break
                    }
                }
            }
        }
    }

    // 呼叫鏈碼：批准請求
    log.Printf("[Debug] 批准授權請求: %s", req.RequestId)
    _, err = contract.SubmitTransaction(
        "ApproveAndAuthorizeAccess", 
        req.RequestId,
    )
    if err != nil {
        fc.PrintGatewayError(err)
        return nil, status.Error(codes.Internal, "更新授權狀態失敗")
    }

    // 最佳努力：批准後追加 insurer 的 wrapped key
    if pendingReq.ReportID != "" && pendingReq.RequesterHash != "" {
        // 統一使用 hash 值（與註冊邏輯一致）
        insurerHash := pendingReq.RequesterHash
        
        // 2) 取得 clinicId（由病患的報告 meta）
        metasRaw, mErr := contract.EvaluateTransaction("ListMyReportMeta")
        if mErr != nil {
            log.Printf("[Warn] 取得報告 meta 失敗: %v", mErr)
        } else {
            var metas []struct { ReportID string `json:"reportId"`; ClinicID string `json:"clinicId"` }
            if u2 := json.Unmarshal(metasRaw, &metas); u2 != nil {
                log.Printf("[Warn] 解析報告 meta 失敗: %v", u2)
            } else {
                var clinicID string
                for _, m := range metas { if m.ReportID == pendingReq.ReportID { clinicID = m.ClinicID; break } }
                if clinicID == "" { 
                    log.Printf("[Warn] 找不到報告的 clinicId reportId=%s", pendingReq.ReportID)
                } else {
                    // 計算診所 hash（與註冊時一致）
                    clinicHash := database.HashString(clinicID)
                    // 計算病患 hash（與註冊時一致）
                    patientHash := database.HashString(userID)
                    
                    // 3) 讀取報告 Envelope
                    envRaw, rErr := contract.EvaluateTransaction("ReadMyReport", pendingReq.ReportID)
                    if rErr != nil {
                        log.Printf("[Warn] 讀取報告失敗 reportId=%s err=%v", pendingReq.ReportID, rErr)
                    } else {
                        // 4) 優先以病患 unwrap → 以 insurer 重包；若無 patient 包鍵則回退用診所 unwrap
                        tw, iw := sw.NewTransitWrapperFromEnv(); if iw != nil { log.Printf("[Warn] Vault 初始化失敗: %v", iw)
                        } else {
                            // 檢查是否已有 patient wrapped key
                            var envChk struct{ WrappedKeys map[string]any `json:"wrappedKeys"` }
                            usePatient := false
                            if jerr := json.Unmarshal(envRaw, &envChk); jerr == nil {
                                if envChk.WrappedKeys != nil { if _, ok := envChk.WrappedKeys["patient"]; ok { usePatient = true } }
                            }
                            unwrapLabel := "clinic"
                            unwrapBase := "clinic-"+clinicHash
                            if usePatient {
                                unwrapLabel = "patient"
                                unwrapBase = "user-"+patientHash
                            }
                            updated, wErr := tw.AddRecipientTransitFrom(ctx, envRaw, unwrapLabel, unwrapBase, "insurer", "insurer-"+insurerHash)
                            if wErr != nil {
                                log.Printf("[Warn] 追加 insurer 包鍵失敗 reportId=%s err=%v (unwrap via %s:%s)", pendingReq.ReportID, wErr, unwrapLabel, unwrapBase)
                            } else {
                                if _, sErr := contract.SubmitTransaction("UpdateReport", pendingReq.ReportID, string(updated)); sErr != nil {
                                    log.Printf("[Warn] UpdateReport 失敗 reportId=%s err=%v", pendingReq.ReportID, sErr)
                                } else {
                                    log.Printf("[Info] 已為報告 %s 追加 insurer 包鍵: insurer-%s-wrap（unwrap via %s）", pendingReq.ReportID, insurerHash, unwrapLabel)
                                }
                            }
                        }
                    }
                }
            }
        }
    } else {
        log.Printf("[Warn] 缺少必要資訊，略過追加包鍵 requestId=%s", req.RequestId)
    }

    return &pb.ApproveAccessRequestResponse{ Success: true, Message: "已批准授權請求" }, nil
}

// HandleRejectAccessRequest 處理授權請求的拒絕
/**
 * @notice 病患拒絕授權請求
 * @dev 病患身份 → Submit RejectAccessRequest
 * @param ctx 請求上下文
 * @param req 拒絕請求（requestId）
 * @param wallet 錢包介面
 * @param builder Gateway 建構器
 * @return *pb.RejectAccessRequestResponse 結果, error 內部錯誤
 */
func HandleRejectAccessRequest(
	ctx context.Context,
	req *pb.RejectAccessRequestRequest,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (*pb.RejectAccessRequestResponse, error) {

	userID, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析授權資訊")
	}

    entry, ok := wallet.GetResolved(userID)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, err
	}
	defer gw.Close()

	// 呼叫鏈碼
	_, err = contract.SubmitTransaction(
		"RejectAccessRequest",
		req.RequestId,
	)
	if err != nil {
		fc.PrintGatewayError(err)
		return nil, status.Error(codes.Internal, "拒絕授權請求失敗")
	}

	return &pb.RejectAccessRequestResponse{
		Success: true,
		Message: "已拒絕授權請求",
	}, nil
}

// HandleListAuthorizedReports 獲取已授權的報告列表
/**
 * @notice 保險業者查看自己可存取的授權報告列表
 * @dev 保險業者身份 → Evaluate ListAuthorizedReports → 整合 DB 顯示名稱
 * @param ctx 請求上下文
 * @param _ 空請求
 * @param wallet 錢包介面
 * @param builder Gateway 建構器
 * @return *pb.ListAuthorizedReportsResponse 清單, error 內部錯誤
 */
func HandleListAuthorizedReports(
	ctx context.Context,
	_ *emptypb.Empty,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (*pb.ListAuthorizedReportsResponse, error) {

	// 取得JWT中的使用者ID（保險業者）
	insurerId, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析授權資訊")
	}

	// 檢查是否為有效的保險業者
	_, err = database.GetInsurerPassword(insurerId)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "只有保險業者可以存取此資料")
	}

	// 取得保險業者錢包
    entry, ok := wallet.GetResolved(insurerId)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	// 連接區塊鏈
	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, status.Error(codes.Internal, "區塊鏈連接失敗")
	}
	defer gw.Close()

	// 呼叫智能合約方法
	result, err := contract.EvaluateTransaction("ListAuthorizedReports")
	if err != nil {
		fc.PrintGatewayError(err)
		return nil, status.Error(codes.Internal, "查詢授權報告失敗")
	}

	// 解析鏈碼回傳的JSON結果
	var rawList []map[string]interface{}
	if err := json.Unmarshal(result, &rawList); err != nil {
		return nil, status.Errorf(codes.Internal, "回傳格式錯誤: %v", err)
	}

	// 轉換為 protobuf 格式
	var reports []*pb.AuthorizedReport
	for _, r := range rawList {
		// 從資料庫獲取病患資訊（將 PatientHash 轉換為真實姓名）
		log.Printf("[Debug] 嘗試使用鏈碼返回的 PatientHash 查詢資料庫: %s", r["patientHash"])
		user, err := database.GetUserByHash(r["patientHash"].(string))
		var patientName string
		if err != nil {
			log.Printf("[Warning] 無法獲取病患資訊: %v", err)
			patientName = "未知病患" // 如果查詢失敗，顯示預設值
		} else {
			patientName = user.Name
		}

		// 將時間戳轉換為日期字串，並處理 nil 的情況
		var createdAt, expiry int64
		
		if r["createdAt"] != nil {
			createdAt = int64(r["createdAt"].(float64))
		} else {
			createdAt = time.Now().Unix()
		}
		
		if r["expiry"] != nil {
			expiry = int64(r["expiry"].(float64))
		} else {
			// 如果沒有設定過期時間，預設為創建時間加上 30 天
			expiry = createdAt + (30 * 24 * 60 * 60)
		}
		
		date := time.Unix(createdAt, 0).Format("2006-01-02")
		expiryDate := time.Unix(expiry, 0).Format("2006-01-02")
		
		report := &pb.AuthorizedReport{
			ReportId:    r["reportId"].(string),
			PatientId:   r["patientHash"].(string),
			PatientName: patientName, // 添加病患真實姓名
			Date:        date,
			Expiry:      expiryDate,
		}
		reports = append(reports, report)
	}

	log.Printf("[Info] 已授權報告: %v", reports)

	return &pb.ListAuthorizedReportsResponse{
		Reports: reports,
	}, nil
}

// HandleListReportMetaByPatientID 獲取特定病患的報告元數據 (不含健檢數據)
/**
 * @notice 保險業者按病患ID查詢報告元數據（不含內容）
 * @dev 保險業者身份 → Evaluate ListReportMetaByPatientID
 * @param ctx 請求上下文
 * @param req 病患ID請求
 * @param wallet 錢包介面
 * @param builder Gateway 建構器
 * @return *pb.ListReportMetaResponse 清單, error 內部錯誤
 */
func HandleListReportMetaByPatientID(
	ctx context.Context,
	req *pb.PatientIDRequest,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (*pb.ListReportMetaResponse, error) {

	// 取得JWT中的使用者ID（保險業者）
	insurerId, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析授權資訊")
	}

	// 檢查是否為有效的保險業者
	_, err = database.GetInsurerPassword(insurerId)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "只有保險業者可以查詢病患報告元數據")
	}

	// 檢查請求
	if req.PatientId == "" {
		return nil, status.Error(codes.InvalidArgument, "必須提供病患ID")
	}

	// 取得保險業者錢包
    entry, ok := wallet.GetResolved(insurerId)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	// 連接區塊鏈
	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, status.Error(codes.Internal, "區塊鏈連接失敗")
	}
	defer gw.Close()

	// 呼叫智能合約方法
	result, err := contract.EvaluateTransaction("ListReportMetaByPatientID", req.PatientId)
	if err != nil {
		fc.PrintGatewayError(err)
		return nil, status.Error(codes.Internal, "查詢病患報告元數據失敗")
	}

	// 解析鏈碼回傳的JSON結果
	type rawReportMeta struct {
		ReportID  string `json:"reportId"`
		ClinicID  string `json:"clinicId"`
		CreatedAt int64  `json:"createdAt"`
	}

	var rawList []rawReportMeta
	if err := json.Unmarshal(result, &rawList); err != nil {
		return nil, status.Errorf(codes.Internal, "回傳格式錯誤: %v", err)
	}

	// 轉換為 protobuf 格式
	var reports []*pb.ReportMeta
	for _, r := range rawList {
		reports = append(reports, &pb.ReportMeta{
			ReportId:  r.ReportID,
			ClinicId:  r.ClinicID,
			CreatedAt: r.CreatedAt,
		})
	}

	log.Printf("[Info] 查詢到病患 %s 的報告元數據 %d 筆", req.PatientId, len(reports))
	log.Printf("[Info] 數據: %v", reports)
	return &pb.ListReportMetaResponse{
		Reports: reports,
	}, nil
}

// ViewAuthorizedReport 實現保險業者讀取授權報告的服務
/**
 * @notice 保險業者讀取已授權的報告內容
 * @dev 保險業者身份 → Evaluate ReadAuthorizedReport
 * @param ctx 請求上下文
 * @param req 請求（reportId, userId）
 * @param wallet 錢包介面
 * @param builder Gateway 建構器
 * @return *pb.ViewAuthorizedReportResponse 結果, error 內部錯誤
 */
func HandleViewAuthorizedReport(
	ctx context.Context,
	req *pb.ViewAuthorizedReportRequest,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (*pb.ViewAuthorizedReportResponse, error) {

	log.Printf("[Debug] HandleViewAuthorizedReport %s", req)
	// 取得JWT中的使用者ID（保險業者）
	insurerId, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析授權資訊")
	}

	// 檢查是否為有效的保險業者
	_, err = database.GetInsurerPassword(insurerId)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "只有保險業者可以查詢病患報告元數據")
	}

	// 檢查請求
	if req.ReportId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "必須提供報告ID和病患ID")
	}

	// 取得保險業者錢包
    entry, ok := wallet.GetResolved(insurerId)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	// 連接區塊鏈
	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, status.Error(codes.Internal, "區塊鏈連接失敗")
	}
	defer gw.Close()

	log.Printf("[Debug] HandleViewAuthorizedReport %s", req)
    // 呼叫智能合約方法
    result, err := contract.EvaluateTransaction("ReadAuthorizedReport", req.UserId, req.ReportId)
	if err != nil {
		fc.PrintGatewayError(err)
		return nil, status.Error(codes.Internal, "查詢病患報告元數據失敗")
	}

	
	log.Printf("[Info] 查詢到報告: %s", string(result))

    // 以保險業者身分解密回傳
    tw, tErr := sw.NewTransitWrapperFromEnv(); if tErr != nil { return nil, status.Error(codes.Internal, "Vault 初始化失敗") }
    // 統一使用 hash 值（與註冊邏輯一致）
    insurerHash := database.HashString(insurerId)
    pt, dErr := tw.DecryptReportTransit(ctx, result, "insurer", "insurer-"+insurerHash)
    if dErr != nil {
        return nil, status.Error(codes.Internal, "解密失敗")
    }
    return &pb.ViewAuthorizedReportResponse{ Success: true, ResultJson: string(pt) }, nil
}

// HandleListMyAccessRequests 處理保險業者查看自己發出的授權請求
/**
 * @notice 保險業者查看自己發出的授權請求
 * @dev 保險業者身份 → Evaluate ListMyAccessRequests → 補充病患名稱
 * @param ctx 請求上下文
 * @param _ 空請求
 * @param wallet 錢包介面
 * @param builder Gateway 建構器
 * @return *pb.ListMyAccessRequestsResponse 清單, error 內部錯誤
 */
func HandleListMyAccessRequests(
	ctx context.Context,
	_ *emptypb.Empty,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (*pb.ListMyAccessRequestsResponse, error) {

	// 取得JWT中的使用者ID（保險業者）
	insurerId, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析授權資訊")
	}

	// 檢查是否為有效的保險業者
	_, err = database.GetInsurerPassword(insurerId)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "只有保險業者可以查看授權請求")
	}

    // 取得保險業者錢包
    entry, ok := wallet.GetResolved(insurerId)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	// 連接區塊鏈
	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, err
	}
	defer gw.Close()

	// 呼叫智能合約方法
	result, err := contract.EvaluateTransaction("ListMyAccessRequests")
	if err != nil {
		fc.PrintGatewayError(err)
		return nil, status.Error(codes.Internal, "查詢授權請求失敗")
	}

	// 解析鏈碼回傳的JSON結果
	var raws []rawAccessRequest
	if err := json.Unmarshal(result, &raws); err != nil {
		return nil, status.Error(codes.Internal, "解析結果失敗")
	}

	// 轉換為 protobuf 格式
	var requests []*pb.AccessRequest
	for _, r := range raws {
		// 從資料庫獲取用戶資訊（將 PatientHash 轉換為真實姓名）
		log.Printf("[Debug] 嘗試使用鏈碼返回的 PatientHash 查詢資料庫: %s", r.PatientHash)
		user, err := database.GetUserByHash(r.PatientHash)
		var patientName string
		if err != nil {
			log.Printf("[Warning] 無法獲取用戶資訊: %v", err)
			patientName = "未知用戶" // 如果查詢失敗，顯示預設值
		} else {
			patientName = user.Name
		}

		requests = append(requests, &pb.AccessRequest{
			RequestId:     r.RequestID,
			ReportId:      r.ReportID,
			PatientHash:   r.PatientHash,
			RequesterHash: r.RequesterHash,
			PatientName:   patientName, // 添加用戶真實姓名
			Reason:        r.Reason,
			RequestedAt:   r.RequestedAt,
			Expiry:        r.Expiry,
			Status:        r.Status,
		})
	}

	return &pb.ListMyAccessRequestsResponse{
		Success: true,
		Requests: requests,
	}, nil
}

// 添加中間結構以匹配鏈碼的 AuthTicket 結構
/**
 * @notice 中間結構：對應鏈碼端的 AuthTicket JSON
 * @dev 僅供本檔案內 JSON 反序列化使用
 */
type rawAuthTicket struct {
	DocType     string `json:"docType"`
	PatientHash string `json:"patientHash"`
	TargetHash  string `json:"targetHash"`
	ReportID    string `json:"reportId"`
	GrantedAt   int64  `json:"grantedAt"`
	Expiry      int64  `json:"expiry"`
}

func HandleListMyAuthorizedTickets(
	ctx context.Context,
	_ *emptypb.Empty,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (*pb.ListAuthorizedTicketsResponse, error) {
	
	userID, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析授權資訊")
	}

    entry, ok := wallet.GetResolved(userID)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, err
	}
	defer gw.Close()

	result, err := contract.EvaluateTransaction("ListMyAuthorizedTickets")
    if err != nil {
		fc.PrintGatewayError(err)
		return nil, status.Error(codes.Internal, "查詢授權請求失敗")
	}

	// 先解析為中間結構
	var raws []rawAuthTicket
	if err := json.Unmarshal(result, &raws); err != nil {
		return nil, status.Error(codes.Internal, "解析結果失敗")
	}

	// 轉換為 proto 結構
	var tickets []*pb.AuthTicket
	for _, r := range raws {
		// 從資料庫獲取保險業者資訊
		insurer, err := database.GetInsurerByHash(r.TargetHash)
		if err != nil {
			log.Printf("[Warning] 無法獲取保險業者資訊: %v", err)
			// 如果無法獲取保險業者資訊，仍然添加票據，但不包含名稱資訊
			tickets = append(tickets, &pb.AuthTicket{
				PatientHash: r.PatientHash,
				TargetHash:  r.TargetHash,
				ReportId:    r.ReportID,
				GrantedAt:   r.GrantedAt,
				Expiry:      r.Expiry,
			})
			continue
		}

		tickets = append(tickets, &pb.AuthTicket{
			PatientHash:  r.PatientHash,
			TargetHash:   r.TargetHash,
			ReportId:     r.ReportID,
			GrantedAt:    r.GrantedAt,
			Expiry:       r.Expiry,
			RequesterName: insurer.Name,
			CompanyName:   insurer.CompanyName,
		})
	}
	log.Printf("[Info] 查詢到授權票據: %v", tickets)
	return &pb.ListAuthorizedTicketsResponse{
		Success: true,
		Tickets: tickets,
	}, nil
}

// HandleListMyReportMeta 處理病患查詢自己的報告 meta
func HandleListMyReportMeta(
	ctx context.Context, _ *emptypb.Empty,
	wallet wl.WalletInterface, builder fc.GWBuilder) (*pb.ListMyReportMetaResponse, error) {
	
	userID, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析 JWT")
	}
	log.Printf("[Debug] HandleListMyReportMeta %s", userID)

    entry, ok := wallet.GetResolved(userID)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, err
	}
	defer gw.Close()

	// 呼叫鏈碼方法
	result, err := contract.EvaluateTransaction("ListMyReportMeta")
	if err != nil {
		fc.PrintGatewayError(err)
		return nil, status.Error(codes.Internal, "查詢失敗")
	}

	// 解析返回的 JSON
	type rawReportMeta struct {
		ReportID  string `json:"reportId"`
		ClinicID  string `json:"clinicId"`
		CreatedAt int64  `json:"createdAt"`
	}

	var rawList []rawReportMeta
	if err := json.Unmarshal(result, &rawList); err != nil {
		return nil, status.Errorf(codes.Internal, "回傳格式錯誤: %v", err)
	}

	// 轉換為 protobuf 格式
	var reports []*pb.ReportMeta
	for _, r := range rawList {
		reports = append(reports, &pb.ReportMeta{
			ReportId:  r.ReportID,
			ClinicId:  r.ClinicID,
			CreatedAt: r.CreatedAt,
		})
	}

	return &pb.ListMyReportMetaResponse{
		Reports: reports,
	}, nil
}

// HandleReadMyReport 處理病患讀取自己的完整報告內容
func HandleReadMyReport(
	ctx context.Context,
	req *pb.ReadMyReportRequest,
	wallet wl.WalletInterface, builder fc.GWBuilder) (*pb.ReadMyReportResponse, error) {

	userID, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析 JWT")
	}

    entry, ok := wallet.GetResolved(userID)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, err
	}
	defer gw.Close()

    // 呼叫鏈碼方法
    result, err := contract.EvaluateTransaction("ReadMyReport", req.ReportId)
	if err != nil {
		fc.PrintGatewayError(err)
		return nil, status.Error(codes.Internal, "讀取報告失敗")
	}

    // 以病患身分解密回傳給前端
    tw, tErr := sw.NewTransitWrapperFromEnv()
    if tErr != nil {
        return nil, status.Error(codes.Internal, "Vault 初始化失敗")
    }
    // 統一使用 hash 值（與註冊邏輯一致）
    patientHash := database.HashString(userID)
    pt, dErr := tw.DecryptReportTransit(ctx, result, "patient", "user-"+patientHash)
    if dErr != nil {
        return nil, status.Error(codes.Internal, "解密失敗")
    }
    return &pb.ReadMyReportResponse{ Success: true, ResultJson: string(pt) }, nil
}	