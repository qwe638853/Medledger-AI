package grpc

import (
    "context"
    "log"
    "net"
    pb "sdk_test/proto" // 改成你的 package 名，通常跟proto package health對應
	"google.golang.org/grpc"

)

type server struct {
    pb.UnimplementedHealthServiceServer
}


// 2. 實作你 proto 裡面的 RPC function
func (s *server) UploadReport(ctx context.Context, req *pb.UploadReportRequest) (*pb.UploadReportResponse, error) {
    log.Printf("Received UploadReport: %v", req)
    return &pb.UploadReportResponse{Message: "Report Uploaded Successfully"}, nil
}

func (s *server) ClaimReport(ctx context.Context, req *pb.ClaimReportRequest) (*pb.ClaimReportResponse, error) {
    log.Printf("Received ClaimReport: %v", req)
    return &pb.ClaimReportResponse{Message: "Report Claimed Successfully"}, nil
}

func (s *server) ReadReport(ctx context.Context, req *pb.ReadReportRequest) (*pb.ReadReportResponse, error) {
    log.Printf("Received ReadReport: %v", req)
    return &pb.ReadReportResponse{ReportContent: "Fake report content..."}, nil
}

// 3. 啟動 gRPC server
func StartGrpcServer() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterHealthServiceServer(grpcServer, &server{})  // 這裡要註冊上面的 server

    log.Println("gRPC server is running at :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
