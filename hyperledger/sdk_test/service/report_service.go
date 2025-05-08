package service

import (
	"context"

	"sdk_test/fabric"
	fc "sdk_test/fabric"
	pb "sdk_test/proto"
	wl "sdk_test/wallet"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HandleUploadReport 驗證請求 → 存 SQLite → 調用 Fabric
func HandleUploadReport(
	ctx context.Context,
	req *pb.UploadReportRequest,
	wallet *wl.Wallet, builder fc.GWBuilder) (*pb.UploadReportResponse, error) {

	userID := req.UserId // 例如從 metadata 或 JWT 取得
	entry, ok := wallet.Get(userID)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "錢包不存在")
	}

	// 依使用者身分建立 Gateway + Contract
	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return nil, err
	}
	defer gw.Close()

	// 呼叫鏈碼
	_, err = contract.SubmitTransaction(
		"UploadReport",
		req.ReportId,
		req.PatientHash,
		req.TestResultsJson,
	)
	if err != nil {
		fabric.PrintGatewayError(err) // 看錯誤細節
		return nil, status.Error(codes.Internal, "鏈上交易失敗")
	}

	return &pb.UploadReportResponse{
		Success: true, Message: "上傳成功",
	}, nil
}
