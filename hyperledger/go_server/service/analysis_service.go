package service

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"go_server/database"
	fc "go_server/fabric"
	pb "go_server/proto"
	sw "go_server/secure/wrap"
	ut "go_server/utils"
	wl "go_server/wallet"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
 * @notice 處理前端請求健檢資料分析
 * @dev 從鏈上獲取報告數據 → 調用 Python Backend gRPC 服務進行分析 → 返回結果
 * @param ctx 請求上下文
 * @param req 分析請求（report_id, user_id, analysis_type）
 * @param wallet 錢包介面
 * @param builder Gateway 建構器
 * @return *pb.HealthAnalysisResponse 分析結果, error 內部錯誤
 */
func HandleGetHealthAnalysis(
	ctx context.Context,
	req *pb.AnalyzeReportRequest,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (*pb.HealthAnalysisResponse, error) {

	log.Printf("[GetHealthAnalysis] 收到分析請求: report_id=%s, user_id=%s, analysis_type=%s",
		req.ReportId, req.UserId, req.AnalysisType)

	// 驗證請求參數
	if req.ReportId == "" {
		return nil, status.Error(codes.InvalidArgument, "報告ID不能為空")
	}

	// 確定分析類型（默認為 user）
	analysisType := req.AnalysisType
	if analysisType == "" {
		analysisType = "user"
	}
	if analysisType != "user" && analysisType != "insurer" {
		return nil, status.Error(codes.InvalidArgument, "分析類型必須為 'user' 或 'insurer'")
	}

	// 從 JWT 獲取用戶 ID（如果未提供）
	userID := req.UserId
	if userID == "" {
		extractedID, err := ut.ExtractUserIDFromContext(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "無法解析用戶ID")
		}
		userID = extractedID
	}

	// 從鏈上讀取報告數據
	log.Printf("[GetHealthAnalysis] 從鏈上讀取報告: report_id=%s", req.ReportId)
	reportData, err := readReportFromChain(ctx, req.ReportId, userID, wallet, builder)
	if err != nil {
		log.Printf("[GetHealthAnalysis] 讀取報告失敗: %v", err)
		return nil, status.Error(codes.Internal, "讀取報告失敗: "+err.Error())
	}

	// 解析報告 JSON 數據
	var testResults map[string]interface{}
	if err := json.Unmarshal([]byte(reportData), &testResults); err != nil {
		log.Printf("[GetHealthAnalysis] 解析報告 JSON 失敗: %v", err)
		return nil, status.Error(codes.Internal, "解析報告數據失敗")
	}

	// 調用 Python Backend 的 gRPC 服務進行分析
	log.Printf("[GetHealthAnalysis] 調用 Python Backend 進行分析: analysis_type=%s", analysisType)
	analysisResult, err := callPythonAnalysisService(ctx, req.ReportId, userID, reportData, analysisType)
	if err != nil {
		log.Printf("[GetHealthAnalysis] Python Backend 分析失敗: %v", err)
		return &pb.HealthAnalysisResponse{
			Success: false,
			Message: "分析失敗: " + err.Error(),
		}, nil
	}

	// 構建回應
	response := &pb.HealthAnalysisResponse{
		Success: true,
		Message: "分析完成",
	}

	if analysisType == "user" && analysisResult.UserAnalysis != nil {
		response.UserAnalysis = analysisResult.UserAnalysis
	} else if analysisType == "insurer" && analysisResult.InsurerAnalysis != nil {
		response.InsurerAnalysis = analysisResult.InsurerAnalysis
	}

	log.Printf("[GetHealthAnalysis] 分析完成: report_id=%s", req.ReportId)
	return response, nil
}

/**
 * @notice 從鏈上讀取報告數據
 * @dev 類似 HandleReadMyReport 的邏輯，但不返回給前端，而是返回原始數據
 * @param ctx 請求上下文
 * @param reportId 報告ID
 * @param userID 用戶ID
 * @param wallet 錢包介面
 * @param builder Gateway 建構器
 * @return string 報告數據 JSON, error 錯誤
 */
func readReportFromChain(
	ctx context.Context,
	reportId string,
	userID string,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (string, error) {

	entry, ok := wallet.GetResolved(userID)
	if !ok {
		return "", status.Error(codes.PermissionDenied, "錢包不存在")
	}

	contract, gw, err := builder.NewContract(entry.ID, entry.Signer)
	if err != nil {
		return "", err
	}
	defer gw.Close()

	// 呼叫鏈碼方法讀取報告
	result, err := contract.EvaluateTransaction("ReadMyReport", reportId)
	if err != nil {
		fc.PrintGatewayError(err)
		return "", status.Error(codes.Internal, "讀取報告失敗")
	}

	// 解密報告數據
	tw, tErr := sw.NewTransitWrapperFromEnv()
	if tErr != nil {
		return "", status.Error(codes.Internal, "Vault 初始化失敗")
	}

	patientHash := database.HashString(userID)
	pt, dErr := tw.DecryptReportTransit(ctx, result, "patient", "user-"+patientHash)
	if dErr != nil {
		return "", status.Error(codes.Internal, "解密失敗")
	}

	return string(pt), nil
}

/**
 * @notice 調用 Python Backend 的 gRPC 服務進行分析
 * @dev 連接到 Python Backend 的 gRPC 服務（端口從環境變數讀取，默認 50052）
 * @param ctx 請求上下文
 * @param reportId 報告ID
 * @param patientId 病患ID
 * @param testResultsJson 測試結果 JSON
 * @param analysisType 分析類型（user 或 insurer）
 * @return *pb.HealthAnalysisResponse 分析結果, error 錯誤
 */
func callPythonAnalysisService(
	ctx context.Context,
	reportId string,
	patientId string,
	testResultsJson string,
	analysisType string) (*pb.HealthAnalysisResponse, error) {

	// 從環境變數獲取 Python Backend 的 gRPC 地址
	pythonBackendAddr := os.Getenv("PYTHON_BACKEND_GRPC_ADDR")
	if pythonBackendAddr == "" {
		pythonBackendAddr = "localhost:50052" // 默認端口（Python Backend 應該運行在不同的端口）
	}

	log.Printf("[callPythonAnalysisService] 連接到 Python Backend: %s", pythonBackendAddr)

	// 建立 gRPC 連線
	conn, err := grpc.NewClient(pythonBackendAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, status.Error(codes.Internal, "無法連接到 Python Backend: "+err.Error())
	}
	defer conn.Close()

	client := pb.NewHealthServiceClient(conn)

	// 根據分析類型調用對應的方法
	if analysisType == "user" {
		// 調用用戶分析
		userResp, err := client.AnalyzeHealthReportForUser(ctx, &pb.AnalyzeHealthReportRequest{
			ReportId:        reportId,
			PatientId:       patientId,
			TestResultsJson: testResultsJson,
		})
		if err != nil {
			return nil, status.Error(codes.Internal, "用戶分析失敗: "+err.Error())
		}

		return &pb.HealthAnalysisResponse{
			Success:      true,
			Message:      "分析完成",
			UserAnalysis: userResp,
		}, nil

	} else if analysisType == "insurer" {
		// 調用保險業者分析
		insurerResp, err := client.AnalyzeHealthReportForInsurer(ctx, &pb.AnalyzeHealthReportRequest{
			ReportId:        reportId,
			PatientId:       patientId,
			TestResultsJson: testResultsJson,
		})
		if err != nil {
			return nil, status.Error(codes.Internal, "保險業者分析失敗: "+err.Error())
		}

		return &pb.HealthAnalysisResponse{
			Success:         true,
			Message:         "分析完成",
			InsurerAnalysis: insurerResp,
		}, nil
	}

	return nil, status.Error(codes.InvalidArgument, "無效的分析類型")
}

