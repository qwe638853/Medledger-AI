package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"sdk_test/database"
	fc "sdk_test/fabric"
	pb "sdk_test/proto"
	ut "sdk_test/utils"
	wl "sdk_test/wallet"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

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
	entry, ok := wallet.Get(userID)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	log.Printf("[Debug] UploadReport args: reportID=%s, patientHash=%s, testResult=%s",
		req.ReportId, req.UserId, req.TestResultsJson)

	// 依使用者身分建立 Gateway + Contract
	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, err
	}
	defer gw.Close()

	sum := sha256.Sum256([]byte(req.UserId))
	hashedUserID := hex.EncodeToString(sum[:])
	log.Printf("[Debug] 查詢患者雜湊: %s", hashedUserID)

	// 呼叫鏈碼
	_, err = contract.SubmitTransaction(
		"UploadReport",
		req.ReportId,
		hashedUserID,
		req.TestResultsJson,
	)
	if err != nil {
		fc.PrintGatewayError(err) // 看錯誤細節
		return nil, status.Error(codes.Internal, "鏈上交易失敗")
	}

	return &pb.UploadReportResponse{
		Success: true, Message: "上傳成功",
	}, nil
}

// HandleListMyReports 呼叫鏈碼查詢病人自己的報告
func HandleListMyReports(
	ctx context.Context, _ *emptypb.Empty,
	wallet wl.WalletInterface, builder fc.GWBuilder) (*pb.ListMyReportsResponse, error) {

	userID, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析 JWT")
	}

	entry, ok := wallet.Get(userID)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, err
	}
	defer gw.Close()

	sum := sha256.Sum256([]byte(userID))
	hashedUserID := hex.EncodeToString(sum[:])
	log.Printf("[Debug] 查詢患者雜湊: %s", hashedUserID)

	// 5. EvaluateTransaction 傳入 hashedUserID 給鏈碼使用
	result, err := contract.EvaluateTransaction("ListMyReports", hashedUserID)
	if err != nil {
		fc.PrintGatewayError(err)
		return nil, status.Error(codes.Internal, "查詢失敗")
	}
	// 建立中介結構以對應 camelCase JSON 欄位
	type rawReport struct {
		ReportID    string `json:"reportId"`
		ClinicID    string `json:"clinicId"`
		PatientHash string `json:"patientHash"`
		ResultJson  string `json:"resultJson"`
		CreatedAt   int64  `json:"createdAt"`
	}

	var rawList []rawReport
	if err := json.Unmarshal(result, &rawList); err != nil {
		return nil, status.Errorf(codes.Internal, "回傳格式錯誤: %v", err)
	}

	// 映射成 protobuf 格式
	var reports []*pb.Report
	for i, r := range rawList {
		log.Printf("[Report #%d] ReportID: %s, PatientHash: %s, TestResults: %s",
			i, r.ReportID, r.PatientHash, r.ResultJson)

		reports = append(reports, &pb.Report{
			ReportId:    r.ReportID,
			ClinicId:    r.ClinicID,
			PatientHash: r.PatientHash,
			ResultJson:  r.ResultJson,
			CreatedAt:   r.CreatedAt,
		})
	}

	return &pb.ListMyReportsResponse{
		Reports: reports,
	}, nil
}

// HandleRequestAccess 處理保險業者請求授權
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
	entry, ok := wallet.Get(requesterId)
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
	RequesterHash string `json:"requesterHash"`
	Reason       string `json:"reason"`
	RequestedAt  int64  `json:"requestedAt"`
	Expiry       int64  `json:"expiry"`
	Status       string `json:"status"`
}
// HandleListAccessRequests 列出病患的所有授權請求
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

	entry, ok := wallet.Get(userID)
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
func HandleApproveAccessRequest(
	ctx context.Context,
	req *pb.ApproveAccessRequestRequest,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (*pb.ApproveAccessRequestResponse, error) {

	userID, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析授權資訊")
	}

	entry, ok := wallet.Get(userID)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, err
	}
	defer gw.Close()

	// 呼叫鏈碼
	log.Printf("[Debug] 批准授權請求: %s", req.RequestId)
	_, err = contract.SubmitTransaction(
		"ApproveAndAuthorizeAccess", 
		req.RequestId,
		"APPROVED",
	)
	if err != nil {
		fc.PrintGatewayError(err)
		return nil, status.Error(codes.Internal, "更新授權狀態失敗")
	}

	return &pb.ApproveAccessRequestResponse{
		Success: true,
		Message: "已批准授權請求",
	}, nil
}

// HandleRejectAccessRequest 處理授權請求的拒絕
func HandleRejectAccessRequest(
	ctx context.Context,
	req *pb.RejectAccessRequestRequest,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (*pb.RejectAccessRequestResponse, error) {

	userID, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析授權資訊")
	}

	entry, ok := wallet.Get(userID)
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
		"UpdateAccessRequestStatus",
		req.RequestId,
		"REJECTED",
	)
	if err != nil {
		fc.PrintGatewayError(err)
		return nil, status.Error(codes.Internal, "更新授權狀態失敗")
	}

	return &pb.RejectAccessRequestResponse{
		Success: true,
		Message: "已拒絕授權請求",
	}, nil
}

// HandleListAuthorizedReports 獲取已授權的報告列表
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
	entry, ok := wallet.Get(insurerId)
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
			ReportId:  r["reportId"].(string),
			PatientId: r["patientHash"].(string),
			Date:      date,
			Expiry:    expiryDate,
		}
		reports = append(reports, report)
	}

	log.Printf("[Info] 已授權報告: %v", reports)

	return &pb.ListAuthorizedReportsResponse{
		Reports: reports,
	}, nil
}

// HandleListReportMetaByPatientID 獲取特定病患的報告元數據 (不含健檢數據)
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
	entry, ok := wallet.Get(insurerId)
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
	entry, ok := wallet.Get(insurerId)
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

	return &pb.ViewAuthorizedReportResponse{
		Success: true,
		ResultJson: string(result),
	}, nil
}

// HandleListMyAccessRequests 處理保險業者查看自己發出的授權請求
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
	entry, ok := wallet.Get(insurerId)
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
		// 從資料庫獲取病患資訊（如果需要的話）
		requests = append(requests, &pb.AccessRequest{
			RequestId:     r.RequestID,
			ReportId:     r.ReportID,
			PatientHash:  r.PatientHash,
			RequesterHash: r.RequesterHash,
			Reason:       r.Reason,
			RequestedAt:  r.RequestedAt,
			Expiry:       r.Expiry,
			Status:       r.Status,
		})
	}

	return &pb.ListMyAccessRequestsResponse{
		Success: true,
		Requests: requests,
	}, nil
}

// 添加中間結構以匹配鏈碼的 AuthTicket 結構
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

	entry, ok := wallet.Get(userID)
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