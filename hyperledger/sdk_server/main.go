// ca-service/main.go
package main

import (
	//"context"
	//"fmt"
	"log"

	//"google.golang.org/grpc"
	//"net"
	//"os"
	//"path/filepath"

	//ca "ca-service/proto"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type server struct {
	//ca.UnimplementedCAServiceServer
	mspClient *msp.Client
}

func main() {
	// Init Fabric SDK
	sdk, err := fabsdk.New(config.FromFile("../configtx.yaml"))
	if err != nil {
		log.Fatalf("failed to init SDK: %v", err)
	}
	defer sdk.Close()

	// Admin context
	adminContext := sdk.Context(fabsdk.WithUser("admin"), fabsdk.WithOrg("Org1"))

	// MSP client
	mspClient, err := msp.New(adminContext)
	if err != nil {
		log.Fatalf("failed to create MSP client: %v", err)
	}

	// 註冊單一使用者
	username := "User900"
	affiliation := "org1.department1"

	secret, err := mspClient.Register(&msp.RegistrationRequest{
		Name:        username,
		Affiliation: affiliation,
	})
	if err != nil {
		log.Fatalf("⚠️ 註冊失敗 (%s): %v", username, err)
	}

	log.Printf("✅ 使用者 %s 註冊成功，Secret: %s", username, secret)

	/* Start gRPC server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	ca.RegisterCAServiceServer(s, &server{mspClient: mspClient})
	reflection.Register(s)

	log.Println("✅ CA service gRPC server running on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
		*/
}
