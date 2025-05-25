package main

import (
	"context"
	"log"
	"net"
	"net/http"

	db "sdk_test/database"
	fc "sdk_test/fabric"
	pb "sdk_test/proto"
	sc "sdk_test/service"
	wl "sdk_test/wallet"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

)

type server struct {
	pb.UnimplementedHealthServiceServer
	Wallet  *wl.Wallet // ← 注入
	Builder fc.GWBuilder
}

// UploadReport
func (s *server) UploadReport(ctx context.Context, req *pb.UploadReportRequest) (*pb.UploadReportResponse, error) {
	return sc.HandleUploadReport(ctx, req, s.Wallet, s.Builder)
}

// Login
func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return sc.HandleLogin(ctx, req, s.Wallet)
}

// 實現新的註冊方法
func (s *server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterResponse, error) {
	return sc.HandleRegisterUser(ctx, req, s.Wallet)
}

func (s *server) RegisterInsurer(ctx context.Context, req *pb.RegisterInsurerRequest) (*pb.RegisterResponse, error) {
	return sc.HandleRegisterInsurer(ctx, req, s.Wallet)
}

func (s *server) ListMyReports(ctx context.Context, in *emptypb.Empty) (*pb.ListMyReportsResponse, error) {
	return sc.HandleListMyReports(ctx, in, s.Wallet, s.Builder)
}

// 新增 RequestAccess API 方法
func (s *server) RequestAccess(ctx context.Context, req *pb.RequestAccessRequest) (*pb.RequestAccessResponse, error) {
	return sc.HandleRequestAccess(ctx, req, s.Wallet, s.Builder)
}

// 新增 ListAccessRequests API 方法
func (s *server) ListAccessRequests(ctx context.Context, in *emptypb.Empty) (*pb.ListAccessRequestsResponse, error) {
	return sc.HandleListAccessRequests(ctx, in, s.Wallet, s.Builder)
}

// 新增 ApproveAccessRequest API 方法
func (s *server) ApproveAccessRequest(ctx context.Context, req *pb.ApproveAccessRequestRequest) (*pb.ApproveAccessRequestResponse, error) {
	return sc.HandleApproveAccessRequest(ctx, req, s.Wallet, s.Builder)
}

// 新增 RejectAccessRequest API 方法
func (s *server) RejectAccessRequest(ctx context.Context, req *pb.RejectAccessRequestRequest) (*pb.RejectAccessRequestResponse, error) {
	return sc.HandleRejectAccessRequest(ctx, req, s.Wallet, s.Builder)
}



// 新增 ListAuthorizedReports API 方法
func (s *server) ListAuthorizedReports(ctx context.Context, in *emptypb.Empty) (*pb.ListAuthorizedReportsResponse, error) {
	return sc.HandleListAuthorizedReports(ctx, in, s.Wallet, s.Builder)
}

// 新增 ListReportMetaByPatientID API 方法
func (s *server) ListReportMetaByPatientID(ctx context.Context, req *pb.PatientIDRequest) (*pb.ListReportMetaResponse, error) {
	return sc.HandleListReportMetaByPatientID(ctx, req, s.Wallet, s.Builder)
}

func (s *server) ViewAuthorizedReport(ctx context.Context, req *pb.ViewAuthorizedReportRequest) (*pb.ViewAuthorizedReportResponse, error) {
	return sc.HandleViewAuthorizedReport(ctx, req, s.Wallet, s.Builder)
}

func (s *server) ListMyAccessRequests(ctx context.Context, in *emptypb.Empty) (*pb.ListMyAccessRequestsResponse, error) {
	return sc.HandleListMyAccessRequests(ctx, in, s.Wallet, s.Builder)
}

func main() {
	err := db.InitDB("database/user_data.sqlite")
	if err != nil {
		log.Fatalf("❌ SQLite 初始化失敗: %v", err)
	}

	w := wl.New()

	// ③ 建 PeerConnector (只做一次)
	peer, err := fc.NewPeer(
		"localhost:7051",
		"../orgs/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt",
		"peer1.org1.example.com",
	)

	if err != nil {
		log.Fatal(err)
	}

	// ④ 建 Gateway Builder
	builder := fc.GWBuilder{
		Peer:    peer,
		Channel: "channel1",
		CCName:  "health",
	}

	go startGrpcServer(w, builder) // 開 gRPC server
	startHttpGatewayServer()       // 開 gRPC-Gateway server (HTTP server)
}

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
