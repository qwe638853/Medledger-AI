package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go_server/database"
	fc "go_server/fabric"
	pb "go_server/proto"
	sw "go_server/secure/wrap"
	ut "go_server/utils"
	wl "go_server/wallet"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"strings"
	"regexp"
)

type pythonBackendClient struct {
	ctx            context.Context
	cancel         context.CancelFunc
	conn           *grpc.ClientConn
	client         pb.HealthServiceClient
	timeoutSeconds int
}

func (p *pythonBackendClient) Close() {
	if p.cancel != nil {
		p.cancel()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
}

func newPythonBackendClient(logLabel string) (*pythonBackendClient, error) {
	pythonBackendAddr := os.Getenv("PYTHON_BACKEND_GRPC_ADDR")
	if pythonBackendAddr == "" {
		pythonBackendAddr = "localhost:50052"
	}

	timeoutSeconds := 1200
	if timeoutStr := os.Getenv("PYTHON_BACKEND_TIMEOUT_SECONDS"); timeoutStr != "" {
		if parsed, err := strconv.Atoi(timeoutStr); err == nil && parsed > 0 {
			timeoutSeconds = parsed
		}
	}

	log.Printf("[%s] 連接到 Python Backend: %s, 超時時間: %d 秒", logLabel, pythonBackendAddr, timeoutSeconds)

	pyCtx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	conn, err := grpc.DialContext(
		pyCtx,
		pythonBackendAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		cancel()
		return nil, status.Error(codes.Internal, "無法連接到 Python Backend: "+err.Error())
	}

	return &pythonBackendClient{
		ctx:            pyCtx,
		cancel:         cancel,
		conn:           conn,
		client:         pb.NewHealthServiceClient(conn),
		timeoutSeconds: timeoutSeconds,
	}, nil
}

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

	// 從 JWT 獲取當前用戶 ID（調用者身份）
	currentUserID, err := ut.ExtractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "無法解析用戶ID")
	}

	// 確定病患 ID（用於 Python Backend 分析）
	patientID := req.UserId
	if patientID == "" {
		// 如果未提供，且當前用戶是病患，則使用當前用戶 ID
		patientID = currentUserID
	}

	// 判斷當前用戶的 role
	// 先檢查是否為保險業者
	_, err = database.GetInsurerPassword(currentUserID)
	isInsurer := err == nil

	// 根據實際 role 決定讀取方法
	var reportData string
	if isInsurer {
		// 保險業者：使用 ReadAuthorizedReport
		if patientID == "" {
			return nil, status.Error(codes.InvalidArgument, "保險業者分析需要提供病患ID")
		}
		log.Printf("[GetHealthAnalysis] 保險業者讀取授權報告: report_id=%s, patient_id=%s", req.ReportId, patientID)
		reportData, err = readReportFromChain(ctx, req.ReportId, patientID, "insurer", currentUserID, wallet, builder)
	} else {
		// 普通用戶：使用 ReadMyReport（讀取自己的報告）
		log.Printf("[GetHealthAnalysis] 病患讀取自己的報告: report_id=%s", req.ReportId)
		reportData, err = readReportFromChain(ctx, req.ReportId, currentUserID, "user", currentUserID, wallet, builder)
		// 更新 patientID 為當前用戶 ID（用於 Python Backend）
		patientID = currentUserID
	}
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
	// 根據實際 role 決定分析類型
	actualAnalysisType := "user"
	if isInsurer {
		actualAnalysisType = "insurer"
	}
	log.Printf("[GetHealthAnalysis] 調用 Python Backend 進行分析: analysis_type=%s, patient_id=%s", actualAnalysisType, patientID)
	analysisResult, err := callPythonAnalysisService(ctx, req.ReportId, patientID, reportData, actualAnalysisType)
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

	if actualAnalysisType == "user" && analysisResult.UserAnalysis != nil {
		response.UserAnalysis = analysisResult.UserAnalysis
	} else if actualAnalysisType == "insurer" && analysisResult.InsurerAnalysis != nil {
		response.InsurerAnalysis = analysisResult.InsurerAnalysis
	}

	log.Printf("[GetHealthAnalysis] 分析完成: report_id=%s", req.ReportId)
	return response, nil
}

/**
 * @notice 從鏈上讀取報告數據
 * @dev 根據分析類型使用不同的讀取方式：
 *      - user: 使用病患身份讀取自己的報告（ReadMyReport）
 *      - insurer: 使用保險業者身份讀取已授權的報告（ReadAuthorizedReport）
 * @param ctx 請求上下文
 * @param reportId 報告ID
 * @param patientID 病患ID（用於 ReadAuthorizedReport 的 patientHash 參數，或 ReadMyReport 的調用者）
 * @param analysisType 分析類型（user 或 insurer）
 * @param callerID 調用者ID（用於查找錢包和建立 Gateway）
 * @param wallet 錢包介面
 * @param builder Gateway 建構器
 * @return string 報告數據 JSON, error 錯誤
 */
func readReportFromChain(
	ctx context.Context,
	reportId string,
	patientID string,
	analysisType string,
	callerID string,
	wallet wl.WalletInterface,
	builder fc.GWBuilder) (string, error) {

	var entry *wl.Entry
	var ok bool
	var contract *client.Contract
	var gw *client.Gateway
	var err error

	// 根據分析類型決定使用哪個身份
	if analysisType == "insurer" {
		// 保險業者：使用調用者 ID（保險業者）查找錢包
		// 檢查是否為有效的保險業者
		_, err = database.GetInsurerPassword(callerID)
		if err != nil {
			return "", status.Error(codes.PermissionDenied, "只有保險業者可以查詢病患報告")
		}

		// 取得保險業者錢包
		entry, ok = wallet.GetResolved(callerID)
		if !ok {
			return "", status.Error(codes.PermissionDenied, "錢包不存在")
		}

		// 連接區塊鏈
		contract, gw, err = builder.NewContract(entry.ID, entry.Signer)
		if err != nil {
			return "", status.Error(codes.Internal, "區塊鏈連接失敗")
		}
		defer gw.Close()

		// 呼叫智能合約方法 ReadAuthorizedReport（需要 patientHash 和 reportId）
		if patientID == "" {
			return "", status.Error(codes.InvalidArgument, "保險分析需要提供病患ID")
		}
		// 檢查 patientID 是否已經是 hash 值（64 個字符的十六進制字符串）
		// 如果是，直接使用；如果不是，進行 hash
		var patientHash string
		if len(patientID) == 64 && isHexString(patientID) {
			// 已經是 hash 值，直接使用（與 HandleViewAuthorizedReport 保持一致）
			patientHash = patientID
			log.Printf("[readReportFromChain] 使用已提供的 hash 值: %s", patientHash)
		} else {
			// 是原始ID，需要進行 hash
			patientHash = database.HashString(patientID)
			log.Printf("[readReportFromChain] 計算 hash 值: %s -> %s", patientID, patientHash)
		}
		result, err := contract.EvaluateTransaction("ReadAuthorizedReport", patientHash, reportId)
		if err != nil {
			fc.PrintGatewayError(err)
			// 檢查是否為授權相關錯誤
			errMsg := err.Error()
			if strings.Contains(errMsg, "access denied") || strings.Contains(errMsg, "access expired") {
				return "", status.Error(codes.PermissionDenied, "尚未獲得該報告的授權，請先申請授權並等待病患批准")
			}
			return "", status.Error(codes.Internal, "讀取授權報告失敗: "+err.Error())
		}

		// 以保險業者身分解密
		tw, tErr := sw.NewTransitWrapperFromEnv()
		if tErr != nil {
			return "", status.Error(codes.Internal, "Vault 初始化失敗")
		}
		insurerHash := database.HashString(callerID)
		pt, dErr := tw.DecryptReportTransit(ctx, result, "insurer", "insurer-"+insurerHash)
		if dErr != nil {
			return "", status.Error(codes.Internal, "解密失敗: "+dErr.Error())
		}

		return string(pt), nil
	} else {
		// 病患：使用病患身份讀取自己的報告
		entry, ok = wallet.GetResolved(callerID)
		if !ok {
			return "", status.Error(codes.PermissionDenied, "錢包不存在")
		}

		contract, gw, err = builder.NewContract(entry.ID, entry.Signer)
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

		patientHash := database.HashString(callerID)
		pt, dErr := tw.DecryptReportTransit(ctx, result, "patient", "user-"+patientHash)
		if dErr != nil {
			return "", status.Error(codes.Internal, "解密失敗")
		}

		return string(pt), nil
	}
}

// isHexString 檢查字符串是否為有效的十六進制字符串
func isHexString(s string) bool {
	matched, _ := regexp.MatchString("^[0-9a-fA-F]+$", s)
	return matched
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

	backendClient, err := newPythonBackendClient("callPythonAnalysisService")
	if err != nil {
		return nil, err
	}
	defer backendClient.Close()

	pyCtx := backendClient.ctx
	timeoutSeconds := backendClient.timeoutSeconds
	client := backendClient.client

	// 根據分析類型調用對應的方法
	if analysisType == "user" {
		// 調用用戶分析（使用新的 context）
		userResp, err := client.AnalyzeHealthReportForUser(pyCtx, &pb.AnalyzeHealthReportRequest{
			ReportId:        reportId,
			PatientId:       patientId,
			TestResultsJson: testResultsJson,
		})
		if err != nil {
			if status.Code(err) == codes.DeadlineExceeded || pyCtx.Err() == context.DeadlineExceeded {
				return nil, status.Error(codes.DeadlineExceeded, fmt.Sprintf("Python Backend 分析超時（超過 %d 秒）", timeoutSeconds))
			}
			return nil, status.Error(status.Code(err), "用戶分析失敗: "+err.Error())
		}

		return &pb.HealthAnalysisResponse{
			Success:      true,
			Message:      "分析完成",
			UserAnalysis: userResp,
		}, nil

	} else if analysisType == "insurer" {
		// 調用保險業者分析（使用新的 context）
		insurerResp, err := client.AnalyzeHealthReportForInsurer(pyCtx, &pb.AnalyzeHealthReportRequest{
			ReportId:        reportId,
			PatientId:       patientId,
			TestResultsJson: testResultsJson,
		})
		if err != nil {
			if status.Code(err) == codes.DeadlineExceeded || pyCtx.Err() == context.DeadlineExceeded {
				return nil, status.Error(codes.DeadlineExceeded, fmt.Sprintf("Python Backend 分析超時（超過 %d 秒）", timeoutSeconds))
			}
			return nil, status.Error(status.Code(err), "保險業者分析失敗: "+err.Error())
		}

		return &pb.HealthAnalysisResponse{
			Success:         true,
			Message:         "分析完成",
			InsurerAnalysis: insurerResp,
		}, nil
	}

	return nil, status.Error(codes.InvalidArgument, "無效的分析類型")
}

func HandleAnalyzeHealthReportForUser(ctx context.Context, req *pb.AnalyzeHealthReportRequest) (*pb.UserHealthAnalysisResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "請求不可為空")
	}
	if req.ReportId == "" {
		return nil, status.Error(codes.InvalidArgument, "report_id 不能為空")
	}
	if req.PatientId == "" {
		return nil, status.Error(codes.InvalidArgument, "patient_id 不能為空")
	}
	if req.TestResultsJson == "" {
		return nil, status.Error(codes.InvalidArgument, "test_results_json 不能為空")
	}

	backendClient, err := newPythonBackendClient("HandleAnalyzeHealthReportForUser")
	if err != nil {
		return nil, err
	}
	defer backendClient.Close()

	resp, err := backendClient.client.AnalyzeHealthReportForUser(backendClient.ctx, req)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded || backendClient.ctx.Err() == context.DeadlineExceeded {
			return nil, status.Error(codes.DeadlineExceeded, fmt.Sprintf("Python Backend 分析超時（超過 %d 秒）", backendClient.timeoutSeconds))
		}
		return nil, status.Error(status.Code(err), "Python Backend 分析失敗: "+err.Error())
	}

	return resp, nil
}

func HandleAnalyzeHealthReportForInsurer(ctx context.Context, req *pb.AnalyzeHealthReportRequest) (*pb.InsurerHealthAnalysisResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "請求不可為空")
	}
	if req.ReportId == "" {
		return nil, status.Error(codes.InvalidArgument, "report_id 不能為空")
	}
	if req.PatientId == "" {
		return nil, status.Error(codes.InvalidArgument, "patient_id 不能為空")
	}
	if req.TestResultsJson == "" {
		return nil, status.Error(codes.InvalidArgument, "test_results_json 不能為空")
	}

	backendClient, err := newPythonBackendClient("HandleAnalyzeHealthReportForInsurer")
	if err != nil {
		return nil, err
	}
	defer backendClient.Close()

	resp, err := backendClient.client.AnalyzeHealthReportForInsurer(backendClient.ctx, req)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded || backendClient.ctx.Err() == context.DeadlineExceeded {
			return nil, status.Error(codes.DeadlineExceeded, fmt.Sprintf("Python Backend 分析超時（超過 %d 秒）", backendClient.timeoutSeconds))
		}
		return nil, status.Error(status.Code(err), "Python Backend 分析失敗: "+err.Error())
	}

	return resp, nil
}

func HandleParseDocument(ctx context.Context, req *pb.ParseDocumentRequest) (*pb.ParseDocumentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "請求不可為空")
	}
	if len(req.FileContent) == 0 {
		return nil, status.Error(codes.InvalidArgument, "file_content 不能為空")
	}
	if req.FileType == "" {
		return nil, status.Error(codes.InvalidArgument, "file_type 不能為空")
	}

	backendClient, err := newPythonBackendClient("HandleParseDocument")
	if err != nil {
		return nil, err
	}
	defer backendClient.Close()

	resp, err := backendClient.client.ParseDocument(backendClient.ctx, req)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded || backendClient.ctx.Err() == context.DeadlineExceeded {
			return nil, status.Error(codes.DeadlineExceeded, fmt.Sprintf("Python Backend 文件解析超時（超過 %d 秒）", backendClient.timeoutSeconds))
		}
		return nil, status.Error(status.Code(err), "Python Backend 文件解析失敗: "+err.Error())
	}

	return resp, nil
}
