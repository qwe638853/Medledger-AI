package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	db "go_server/database"
	fc "go_server/fabric"
	pb "go_server/proto"
	sc "go_server/service"
	wl "go_server/wallet"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

    "github.com/joho/godotenv"
)

type server struct {
	pb.UnimplementedHealthServiceServer
	Wallet  *wl.Wallet // ← 注入
	Builder fc.GWBuilder
}

/**
 * @notice 上傳健檢報告 API
 * @dev 轉呼叫 service.HandleUploadReport
 * @param ctx 請求上下文
 * @param req 上傳請求
 * @return *pb.UploadReportResponse 結果, error 內部錯誤
 */
// UploadReport
func (s *server) UploadReport(ctx context.Context, req *pb.UploadReportRequest) (*pb.UploadReportResponse, error) {
	return sc.HandleUploadReport(ctx, req, s.Wallet, s.Builder)
}

/**
 * @notice 登入 API
 * @dev 轉呼叫 service.HandleLogin
 * @param ctx 請求上下文
 * @param req 登入請求
 * @return *pb.LoginResponse 結果, error 內部錯誤
 */
// Login
func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return sc.HandleLogin(ctx, req, s.Wallet)
}

/**
 * @notice 用戶註冊 API
 * @dev 轉呼叫 service.HandleRegisterUser
 * @param ctx 請求上下文
 * @param req 註冊請求
 * @return *pb.RegisterResponse 結果, error 內部錯誤
 */
// 實現新的註冊方法
func (s *server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterResponse, error) {
    return sc.HandleRegisterUser(ctx, req, s.Wallet, s.Builder)
}

/**
 * @notice 保險業者註冊 API
 * @dev 轉呼叫 service.HandleRegisterInsurer
 * @param ctx 請求上下文
 * @param req 註冊請求
 * @return *pb.RegisterResponse 結果, error 內部錯誤
 */
func (s *server) RegisterInsurer(ctx context.Context, req *pb.RegisterInsurerRequest) (*pb.RegisterResponse, error) {
	return sc.HandleRegisterInsurer(ctx, req, s.Wallet)
}



/**
 * @notice 病患查詢自己的報告 meta API
 * @dev 轉呼叫 service.HandleListMyReportMeta
 * @param ctx 請求上下文
 * @param in 空請求
 * @return *pb.ListMyReportMetaResponse 清單, error 內部錯誤
 */
// 新增 ListMyReportMeta API 方法
func (s *server) ListMyReportMeta(ctx context.Context, in *emptypb.Empty) (*pb.ListMyReportMetaResponse, error) {
	log.Printf("[Debug] ListMyReportMeta")
	return sc.HandleListMyReportMeta(ctx, in, s.Wallet, s.Builder)
}

/**
 * @notice 病患讀取自己的報告 API
 * @dev 轉呼叫 service.HandleReadMyReport
 * @param ctx 請求上下文
 * @param req 讀取請求
 * @return *pb.ReadMyReportResponse 結果, error 內部錯誤
 */
// 新增 ReadMyReport API 方法
func (s *server) ReadMyReport(ctx context.Context, req *pb.ReadMyReportRequest) (*pb.ReadMyReportResponse, error) {
	log.Printf("[Debug] ReadMyReport")
	return sc.HandleReadMyReport(ctx, req, s.Wallet, s.Builder)
}

/**
 * @notice 保險業者提出授權請求 API
 * @dev 轉呼叫 service.HandleRequestAccess
 * @param ctx 請求上下文
 * @param req 授權請求
 * @return *pb.RequestAccessResponse 結果, error 內部錯誤
 */
// 新增 RequestAccess API 方法
func (s *server) RequestAccess(ctx context.Context, req *pb.RequestAccessRequest) (*pb.RequestAccessResponse, error) {
	return sc.HandleRequestAccess(ctx, req, s.Wallet, s.Builder)
}

/**
 * @notice 病患列出待處理授權請求 API
 * @dev 轉呼叫 service.HandleListAccessRequests
 * @param ctx 請求上下文
 * @param in 空請求
 * @return *pb.ListAccessRequestsResponse 清單, error 內部錯誤
 */
// 新增 ListAccessRequests API 方法
func (s *server) ListAccessRequests(ctx context.Context, in *emptypb.Empty) (*pb.ListAccessRequestsResponse, error) {
	return sc.HandleListAccessRequests(ctx, in, s.Wallet, s.Builder)
}

/**
 * @notice 病患批准授權請求 API
 * @dev 轉呼叫 service.HandleApproveAccessRequest
 * @param ctx 請求上下文
 * @param req 批准請求
 * @return *pb.ApproveAccessRequestResponse 結果, error 內部錯誤
 */
// 新增 ApproveAccessRequest API 方法
func (s *server) ApproveAccessRequest(ctx context.Context, req *pb.ApproveAccessRequestRequest) (*pb.ApproveAccessRequestResponse, error) {
	return sc.HandleApproveAccessRequest(ctx, req, s.Wallet, s.Builder)
}

/**
 * @notice 病患拒絕授權請求 API
 * @dev 轉呼叫 service.HandleRejectAccessRequest
 * @param ctx 請求上下文
 * @param req 拒絕請求
 * @return *pb.RejectAccessRequestResponse 結果, error 內部錯誤
 */
// 新增 RejectAccessRequest API 方法
func (s *server) RejectAccessRequest(ctx context.Context, req *pb.RejectAccessRequestRequest) (*pb.RejectAccessRequestResponse, error) {
	return sc.HandleRejectAccessRequest(ctx, req, s.Wallet, s.Builder)
}



/**
 * @notice 保險業者查看已授權報告列表 API
 * @dev 轉呼叫 service.HandleListAuthorizedReports
 * @param ctx 請求上下文
 * @param in 空請求
 * @return *pb.ListAuthorizedReportsResponse 清單, error 內部錯誤
 */
// 新增 ListAuthorizedReports API 方法
func (s *server) ListAuthorizedReports(ctx context.Context, in *emptypb.Empty) (*pb.ListAuthorizedReportsResponse, error) {
	return sc.HandleListAuthorizedReports(ctx, in, s.Wallet, s.Builder)
}

/**
 * @notice 保險業者按病患ID查報告 meta API
 * @dev 轉呼叫 service.HandleListReportMetaByPatientID
 * @param ctx 請求上下文
 * @param req 病患ID請求
 * @return *pb.ListReportMetaResponse 清單, error 內部錯誤
 */
// 新增 ListReportMetaByPatientID API 方法
func (s *server) ListReportMetaByPatientID(ctx context.Context, req *pb.PatientIDRequest) (*pb.ListReportMetaResponse, error) {
	return sc.HandleListReportMetaByPatientID(ctx, req, s.Wallet, s.Builder)
}

/**
 * @notice 保險業者讀取已授權的報告 API
 * @dev 轉呼叫 service.HandleViewAuthorizedReport
 * @param ctx 請求上下文
 * @param req 讀取請求
 * @return *pb.ViewAuthorizedReportResponse 結果, error 內部錯誤
 */
func (s *server) ViewAuthorizedReport(ctx context.Context, req *pb.ViewAuthorizedReportRequest) (*pb.ViewAuthorizedReportResponse, error) {
	return sc.HandleViewAuthorizedReport(ctx, req, s.Wallet, s.Builder)
}

/**
 * @notice 保險業者查看自己發出的授權請求 API
 * @dev 轉呼叫 service.HandleListMyAccessRequests
 * @param ctx 請求上下文
 * @param in 空請求
 * @return *pb.ListMyAccessRequestsResponse 清單, error 內部錯誤
 */
func (s *server) ListMyAccessRequests(ctx context.Context, in *emptypb.Empty) (*pb.ListMyAccessRequestsResponse, error) {
	return sc.HandleListMyAccessRequests(ctx, in, s.Wallet, s.Builder)
}

/**
 * @notice 病患查看自己授權票據 API
 * @dev 轉呼叫 service.HandleListMyAuthorizedTickets
 * @param ctx 請求上下文
 * @param in 空請求
 * @return *pb.ListAuthorizedTicketsResponse 清單, error 內部錯誤
 */
func (s *server) ListMyAuthorizedTickets(ctx context.Context, in *emptypb.Empty) (*pb.ListAuthorizedTicketsResponse, error) {
	return sc.HandleListMyAuthorizedTickets(ctx, in, s.Wallet, s.Builder)
}

/**
 * @notice 程式進入點：初始化 DB、Peer、Gateway 並啟動 gRPC/HTTP 服務
 * @dev 啟動前測試 Gateway 連線可用性
 * @return 無（阻塞直到服務結束）
 */
func main() {
    // 讀取 .env（若不存在則忽略）
    if err := godotenv.Load(); err != nil {
        log.Printf("[env] 未找到 .env，略過載入：%v", err)
    } else {
        log.Printf("[env] 已載入 .env")
    }
	err := db.InitDB("database/user_data.sqlite")
	if err != nil {
		log.Fatalf("❌ SQLite 初始化失敗: %v", err)
	}

	w := wl.New()

	// 平台身分已不再需要，略過註冊/Enroll

	// ③ 建 PeerConnector (只做一次)
	log.Println("🔗 正在連接到 Peer 節點...")
	peer, err := fc.NewPeer(
		"localhost:7051",
		"../orgs/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt",
		"peer1.org1.example.com",
	)

	if err != nil {
		log.Fatalf("❌ Peer 連線失敗: %v", err)
	}
	log.Println("✅ Peer 連線成功建立")

	// ④ 建 Gateway Builder
	builder := fc.GWBuilder{
		Peer:    peer,
		Channel: "channel1",
		CCName:  "health",
	}

	// 測試Gateway連線
	log.Println("🧪 測試 Gateway 連線...")
	if err := testGatewayConnection(builder, w); err != nil {
		log.Printf("⚠️ Gateway 連線測試失敗: %v", err)
	} else {
		log.Println("✅ Gateway 連線測試成功")
	}

	go startGrpcServer(w, builder) // 開 gRPC server
	startHttpGatewayServer()       // 開 gRPC-Gateway server (HTTP server)
}

/**
 * @notice 測試 Gateway 連線可用性
 * @dev 從錢包取任一身份建立 Gateway，嘗試 Evaluate 簡單鏈碼函式
 * @param builder Gateway 建構器
 * @param wallet 錢包
 * @return error 測試失敗原因
 */
// 添加Gateway連線測試函數
func testGatewayConnection(builder fc.GWBuilder, wallet *wl.Wallet) error {
    // 嘗試使用現有的用戶身份測試連線
    entries, err := wallet.List()
    if err != nil {
        return fmt.Errorf("無法列出錢包條目: %w", err)
    }
    if len(entries) == 0 {
        log.Println("⚠️ 錢包中沒有用戶身份，跳過Gateway測試")
        return nil
    }

    // 遍歷挑選第一個可成功解析的身份
    var picked *wl.Entry
    for _, uid := range entries {
        if e, ok := wallet.GetResolved(uid); ok && e != nil && e.ID != nil && e.Signer != nil {
            picked = e
            break
        }
    }
    if picked == nil {
        log.Println("⚠️ 錢包中無可用身份（可能缺 certUri 或 Transit 金鑰未備妥），跳過Gateway測試")
        return nil
    }

    contract, gw, err := builder.NewContract(picked.ID, picked.Signer)
    if err != nil {
        return fmt.Errorf("無法建立Gateway: %w", err)
    }
    defer gw.Close()

	// 嘗試評估一個簡單的chaincode函數
	_, err = contract.EvaluateTransaction("ListMyReportMeta")
	if err != nil {
		return fmt.Errorf("chaincode 調用失敗: %w", err)
	}

	return nil
}

/**
 * @notice 啟動 gRPC 服務
 * @dev 監聽 :50051 並註冊 HealthServiceServer
 * @param wallet 錢包
 * @param builder Gateway 建構器
 */
func startGrpcServer(wallet *wl.Wallet, builder fc.GWBuilder) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHealthServiceServer(grpcServer, &server{Wallet: wallet, Builder: builder})

	log.Println("gRPC server is running at :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

/**
 * @notice 簡易 CORS 中介層
 * @dev 允許指定 Origin，處理 OPTIONS 預檢請求
 * @param h 下一層 HTTP handler
 * @return http.Handler 包裝後的 handler
 */
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			// 只允許特定 Origin，不要全部 "*"
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// 如果是預檢請求 (OPTIONS)，直接返回 200，不然請求會被擋
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

/**
 * @notice 啟動 HTTP gRPC-Gateway 服務
 * @dev 監聽 :8080，透過 gRPC-Gateway 轉發到 :50051
 */
func startHttpGatewayServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterHealthServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	// 🎯 加上 CORS handler
	handler := allowCORS(mux)

	log.Println("HTTP server listening at :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}

/*
func testUploadClaimRead(contract *client.Contract) {
	testResults := map[string]string{
		"Glu-AC": "95 mg/dL",
		"HbA1c":  "5.3 %",
		"LDL-C":  "125 mg/dL",
	}
	testJSON, _ := json.Marshal(testResults)

	fmt.Println("Submitting UploadReport transaction...")
	_, err := contract.SubmitTransaction("UploadReport", "report001", "patientHash123", string(testJSON))
	if err != nil {
		fc.PrintGatewayErrorDetails(err)
		log.Fatalf("Failed to submit UploadReport transaction: %v", err)
	}
	fmt.Println("✅ Report uploaded successfully")

	fmt.Println("Submitting ClaimReport transaction...")
	_, err = contract.SubmitTransaction("ClaimReport", "report001")
	if err != nil {
		log.Fatalf("Failed to claim report: %v", err)
		fc.PrintGatewayErrorDetails(err)
	}
	fmt.Println("✅ Claimed successfully")

	fmt.Println("Evaluating ReadReport transaction...")
	result, err := contract.EvaluateTransaction("ReadReport", "report001")
	if err != nil {
		log.Fatalf("Failed to read report: %v", err)
	}
	fmt.Println("📄 Report:")
	fmt.Println(string(result))
}
*/
