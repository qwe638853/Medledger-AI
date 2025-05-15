package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"

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
