package service

import (
	"context"
	"encoding/json"
	"fmt"

	fc "sdk_test/fabric"
	pb "sdk_test/proto"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HandleUploadReport 驗證請求 → 存 SQLite → 調用 Fabric
func HandleUploadReport(
	ctx context.Context,
	req *pb.UploadReportRequest,
	contract *client.Contract) (*pb.UploadReportResponse, error) {

	// 1. 基本驗證 --------------------------------------------------------
	if req.ReportId == "" || req.PatientHash == "" || len(req.TestResultJson) == 0 {
		return nil, status.Error(codes.InvalidArgument, "ReportId, PatientHash 與 TestResultJson 皆必填")
	}

	// 2. 先把結果 JSON 解析，確認格式正確 -------------------------------
	var resultMap map[string]string
	if err := json.Unmarshal([]byte(req.TestResultJson), &resultMap); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "TestResultJson 不是合法 JSON: %v", err)
	}

	// 4. 上鏈 (Fabric Contract) -----------------------------------------
	_, err := contract.SubmitTransaction(
		"UploadReport",
		req.ReportId,
		req.PatientHash,
		req.TestResultJson,
	)
	if err != nil {
		fc.PrintGatewayErrorDetails(err)
		// ⚠️ 失敗時記得 Rollback (簡易做法: 刪 DB；正式環境建議用 Tx outbox)

		return nil, status.Errorf(codes.Internal, "鏈上交易失敗: %v", err)
	}

	// 5. 回傳成功 --------------------------------------------------------
	return &pb.UploadReportResponse{
		Success: true,
		Message: fmt.Sprintf("報告 %s 上傳完成", req.ReportId),
	}, nil
}
