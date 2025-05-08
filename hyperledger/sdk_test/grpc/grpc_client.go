package mygrpc

import (
	"context"
	"log"
	"time"

	pb "sdk_test/proto" // 改成你的 package 名，通常跟proto package health對應

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewHealthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 測試 UploadReport
	res, err := client.UploadReport(ctx, &pb.UploadReportRequest{
		ReportId:        "report001",
		PatientHash:     "hash123",
		TestResultsJson: `{"Glu-AC": "95 mg/dL", "HbA1c": "5.3%"}`,
	})
	if err != nil {
		log.Fatalf("could not upload report: %v", err)
	}
	log.Printf("Upload Response: %s", res.Message)

	// 測試 ReadReport
	readRes, err := client.ReadReport(ctx, &pb.ReadReportRequest{
		ReportId: "report001",
	})
	if err != nil {
		log.Fatalf("could not read report: %v", err)
	}
	log.Printf("Report Content: %s", readRes.ReportContent)
}
