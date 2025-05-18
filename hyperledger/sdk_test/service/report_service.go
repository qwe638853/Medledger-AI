package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"
	"strconv"

	"sdk_test/database"
	fc "sdk_test/fabric"
	pb "sdk_test/proto"
	ut "sdk_test/utils"
	wl "sdk_test/wallet"

	"github.com/google/uuid"
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

	// 生成唯一請求ID
	requestId := uuid.New().String()

	// 檢查請求內容
	if req.ReportId == "" || req.PatientId == "" || req.Reason == "" {
		return nil, status.Error(codes.InvalidArgument, "必須提供報告ID、病患ID和申請原因")
	}

	// 設定過期時間，若未提供則預設30天
	expiry := req.Expiry
	if expiry == 0 {
		expiry = time.Now().Unix() + 30*24*60*60 // 30天
	}

	// 儲存申請記錄到SQLite
	now := time.Now().Unix()
	err = database.InsertAccessRequest(
		requestId,
		req.ReportId,
		req.PatientId,
		requesterId,
		req.Reason,
		now,
		expiry,
		"PENDING", // 初始狀態
	)
	if err != nil {
		log.Printf("❌ 儲存授權請求失敗: %v", err)
		return nil, status.Error(codes.Internal, "無法儲存授權請求")
	}

	return &pb.RequestAccessResponse{
		Success:   true,
		RequestId: requestId,
	}, nil
}

// HandleListAccessRequests 列出病患的所有授權請求
func HandleListAccessRequests(
	ctx context.Context,
	_ *emptypb.Empty,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (*pb.ListAccessRequestsResponse, error) {

	// 取得JWT中的使用者ID
	patientId, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析授權資訊")
	}

	// 從數據庫查詢該用戶的授權請求
	requests, err := database.GetAccessRequestsForUser(patientId)
	if err != nil {
		log.Printf("❌ 查詢授權請求失敗: %v", err)
		return nil, status.Error(codes.Internal, "無法查詢授權請求")
	}

	// 轉換為protobuf格式
	var accessRequests []*pb.AccessRequest
	for _, req := range requests {
		accessRequests = append(accessRequests, &pb.AccessRequest{
			RequestId:    req["request_id"].(string),
			ReportId:     req["report_id"].(string),
			PatientHash:  "",
			TargetHash:   req["requester_id"].(string),
			Reason:       req["reason"].(string),
			RequestedAt:  req["requested_at"].(int64),	
			Expiry:       req["expiry"].(int64),
			Status:       req["status"].(string),
		})
	}

	return &pb.ListAccessRequestsResponse{
		Requests: accessRequests,
	}, nil
}

// HandleApproveAccessRequest 批准授權請求
func HandleApproveAccessRequest(
	ctx context.Context,
	req *pb.ApproveAccessRequestRequest,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (*pb.ApproveAccessRequestResponse, error) {

	// 取得JWT中的使用者ID
	patientId, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析授權資訊")
	}

	log.Printf("[Debug] ApproveAccessRequest patientId=%s, requestId=%s", patientId, req.RequestId)
	// 檢查請求ID
	if req.RequestId == "" {
		return nil, status.Error(codes.InvalidArgument, "必須提供請求ID")
	}

	// 獲取授權請求詳情
	requestDetails, err := database.GetAccessRequestById(req.RequestId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "找不到該授權請求")
	}

	// 驗證請求是否屬於當前用戶
	if requestDetails["patient_id"].(string) != patientId {
		return nil, status.Error(codes.PermissionDenied, "無權批准此授權請求")
	}

	// 檢查請求狀態
	if requestDetails["status"].(string) != "PENDING" {
		return nil, status.Error(codes.FailedPrecondition, "只能批准待處理的授權請求")
	}

	// 更新授權請求狀態
	err = database.UpdateAccessRequestStatus(req.RequestId, "APPROVED")
	if err != nil {
		log.Printf("❌ 更新授權請求狀態失敗: %v", err)
		return nil, status.Error(codes.Internal, "無法更新授權請求狀態")
	}

	// 調用區塊鏈授權操作
	entry, ok := wallet.Get(patientId)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, status.Error(codes.Internal, "區塊鏈連接失敗")
	}
	defer gw.Close()

	reportId := requestDetails["report_id"].(string)
	requesterId := requestDetails["requester_id"].(string)
	expiry := requestDetails["expiry"].(int64)

	// 將用戶ID轉為雜湊
	sumPatient := sha256.Sum256([]byte(patientId))
	hashedPatientID := hex.EncodeToString(sumPatient[:])

	sumRequester := sha256.Sum256([]byte(requesterId))
	hashedRequesterID := hex.EncodeToString(sumRequester[:])

	log.Printf("[Debug] AuthorizeAccess reportId=%s, patientId=%s, requesterId=%s, expiry=%d", reportId, hashedPatientID, hashedRequesterID, expiry)	
	// 呼叫區塊鏈授權
	_, err = contract.SubmitTransaction(
		"AuthorizeAccess",
		reportId,
		hashedPatientID,
		hashedRequesterID,
		strconv.FormatInt(expiry, 10),
	)
	if err != nil {
		fc.PrintGatewayError(err)
		// 發生錯誤時回滾資料庫更新
		database.UpdateAccessRequestStatus(req.RequestId, "PENDING")
		return nil, status.Error(codes.Internal, "區塊鏈授權失敗")
	}

	return &pb.ApproveAccessRequestResponse{
		Success: true,
		Message: "已成功批准授權請求",
	}, nil
}

// HandleRejectAccessRequest 拒絕授權請求
func HandleRejectAccessRequest(
	ctx context.Context,
	req *pb.RejectAccessRequestRequest,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (*pb.RejectAccessRequestResponse, error) {

	// 取得JWT中的使用者ID
	patientId, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析授權資訊")
	}

	// 檢查請求ID
	if req.RequestId == "" {
		return nil, status.Error(codes.InvalidArgument, "必須提供請求ID")
	}

	// 獲取授權請求詳情
	requestDetails, err := database.GetAccessRequestById(req.RequestId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "找不到該授權請求")
	}

	// 驗證請求是否屬於當前用戶
	if requestDetails["patient_id"].(string) != patientId {
		return nil, status.Error(codes.PermissionDenied, "無權拒絕此授權請求")
	}

	// 檢查請求狀態
	if requestDetails["status"].(string) != "PENDING" {
		return nil, status.Error(codes.FailedPrecondition, "只能拒絕待處理的授權請求")
	}

	// 更新授權請求狀態
	err = database.UpdateAccessRequestStatus(req.RequestId, "REJECTED")
	if err != nil {
		log.Printf("❌ 更新授權請求狀態失敗: %v", err)
		return nil, status.Error(codes.Internal, "無法更新授權請求狀態")
	}

	return &pb.RejectAccessRequestResponse{
		Success: true,
		Message: "已成功拒絕授權請求",
	}, nil
}

// HandleGetInsurerDashboardStats 獲取保險業者儀表板統計資料
func HandleGetInsurerDashboardStats(
	ctx context.Context,
	_ *emptypb.Empty,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (*pb.InsurerDashboardStatsResponse, error) {

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

	// 從數據庫查詢統計數據
	// 1. 已授權報告數量（從區塊鏈獲取）
	authorizedReports, err := database.GetAuthorizedReportsForInsurer(insurerId)
	totalAuthorized := len(authorizedReports)
	if err != nil {
		log.Printf("❌ 獲取已授權報告時出錯: %v", err)
		// 出錯時不要中斷，繼續查詢其他統計資料
		totalAuthorized = 0
	}

	// 2. 待處理請求數量（從資料庫獲取）
	pendingRequests, err := database.GetPendingRequestsCountForInsurer(insurerId)
	if err != nil {
		log.Printf("❌ 獲取待處理請求數時出錯: %v", err)
		pendingRequests = 0
	}

	// 3. 授權病患數量（從資料庫獲取）
	totalPatients, err := database.GetAuthorizedPatientsCountForInsurer(insurerId)
	if err != nil {
		log.Printf("❌ 獲取授權病患數時出錯: %v", err)
		totalPatients = 0
	}

	return &pb.InsurerDashboardStatsResponse{
		TotalAuthorized: int32(totalAuthorized),
		PendingRequests: int32(pendingRequests),
		TotalPatients:   int32(totalPatients),
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

	// 從數據庫獲取已授權報告
	authorizedReportsData, err := database.GetAuthorizedReportsForInsurer(insurerId)
	if err != nil {
		log.Printf("❌ 獲取已授權報告時出錯: %v", err)
		return nil, status.Error(codes.Internal, "無法獲取已授權報告")
	}

	// 轉換為 protobuf 格式
	var reports []*pb.AuthorizedReport
	for _, r := range authorizedReportsData {
		report := &pb.AuthorizedReport{
			ReportId:  r["report_id"].(string),
			PatientId: r["patient_id"].(string),
			Content:   r["content"].(string),
			Date:      r["date"].(string),
			Expiry:    r["expiry"].(string),
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
