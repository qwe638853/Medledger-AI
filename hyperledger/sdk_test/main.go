package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	db "sdk_test/database"
	fc "sdk_test/fabric"
	pb "sdk_test/proto"
	sc "sdk_test/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedHealthServiceServer
	Contract fc.FabricContract
}

// UploadReport
func (s *server) UploadReport(ctx context.Context, req *pb.UploadReportRequest) (*pb.UploadReportResponse, error) {
	log.Printf("Received UploadReport: %v", req)
	return &pb.UploadReportResponse{
		Success: true,
		Message: "Report Uploaded Successfully",
	}, nil
}

// ClaimReport
func (s *server) ClaimReport(ctx context.Context, req *pb.ClaimReportRequest) (*pb.ClaimReportResponse, error) {
	log.Printf("Received ClaimReport: %v", req)
	return &pb.ClaimReportResponse{
		Success: true,
		Message: "Report Claimed Successfully",
	}, nil
}

// ReadReport
func (s *server) ReadReport(ctx context.Context, req *pb.ReadReportRequest) (*pb.ReadReportResponse, error) {
	log.Printf("Received ReadReport: %v", req)
	return &pb.ReadReportResponse{
		Success:       true,
		ReportContent: "Fake report content...",
	}, nil
}

// Login
func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return sc.HandleLogin(ctx, req)
}

// Register
func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return sc.HandleRegister(ctx, req)
}

func main() {
	err := db.InitDB("database/user_data.sqlite")
	if err != nil {
		log.Fatalf("❌ SQLite 初始化失敗: %v", err)
	}

	fabric := fc.NewFabricContract()
	defer fabric.Gateway.Close()

	go startGrpcServer()     // 開 gRPC server
	startHttpGatewayServer() // 開 gRPC-Gateway server (HTTP server)
}

func startGrpcServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHealthServiceServer(grpcServer, &server{})

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

// 測試功能
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
