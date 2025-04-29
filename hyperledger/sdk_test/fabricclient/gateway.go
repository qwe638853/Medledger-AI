package fabricclient

import (
	"context"
	"crypto/x509"
	"errors" 
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"google.golang.org/grpc/status"      
	"google.golang.org/grpc/credentials"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"

)

//宣告全域變數
const (
	peerEndpoint      = "localhost:7051"
    peerHostOverride  = "peer1.org1.example.com"
    tlsCertPath       = "../orgs/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt"
    mspID             = "Org1MSP"
    certPath          = "../orgs/org1.example.com/users/org1-admin/msp/signcerts/cert.pem"
    keyPath           = "../orgs/org1.example.com/users/org1-admin/msp/keystore/"
    channelName       = "channel1"
    chaincodeName     = "health"
)


type FabricContract struct {
	Gateway  *client.Gateway
	Contract *client.Contract
}

func NewFabricContract() *FabricContract {

	// 建立grpc連線
	grpcConn := newGrpcConnection()


	// 建立身分
	id := newIdentity()
	sign := newSigner()

	// 建立gateway
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
        client.WithSubmitTimeout(5*time.Second),
        client.WithCommitStatusTimeout(1*time.Minute),
        client.WithClientConnection(grpcConn),
	)
	if err != nil {
		log.Fatalf("failed to create gateway connection: %v", err)
	}


    channel := gw.GetNetwork(channelName)
    contract := channel.GetContract(chaincodeName)

	return &FabricContract{
		Gateway:  gw,
		Contract: contract,
	}
    
}

func newGrpcConnection() *grpc.ClientConn {
	
	//讀取TLS憑證
    certificatePEM, err := os.ReadFile(tlsCertPath)
    if err != nil {
        panic(fmt.Errorf("failed to read TLS certificate file: %w", err))
    }

	//解析TLS憑證
    certificate, err := identity.CertificateFromPEM(certificatePEM)
    if err != nil {
        panic(err)
    }

	//建立憑證池
    certPool := x509.NewCertPool()
	//將憑證加入憑證池
    certPool.AddCert(certificate)
	//建立TLS憑證物件(等等grpc連線會用到)
    transportCredentials := credentials.NewClientTLSFromCert(certPool, peerHostOverride)
	//建立grpc連線
    connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
    if err != nil {
        panic(fmt.Errorf("failed to create gRPC connection: %w", err))
    }

    return connection
}


// 載入使用者身份憑證（X.509）
func newIdentity() *identity.X509Identity {
	//讀取身分憑證
    certPEM, err := os.ReadFile(certPath)
    if err != nil {
        panic(fmt.Errorf("failed to read cert file: %w", err))
    }
    //解碼PEM
    certificate, err := identity.CertificateFromPEM(certPEM)
    if err != nil {
        panic(err)
    }

	//建立X509Identity物件(表示身分)
    id, err := identity.NewX509Identity(mspID, certificate)
    if err != nil {
        panic(fmt.Errorf("failed to create identity: %w", err))
    }

    return id
}

// 載入私鑰並建立簽章者
func newSigner() identity.Sign {
    // keystore 資料夾中只有一個檔案（預設）
    files, err := ioutil.ReadDir(keyPath)
    if err != nil || len(files) == 0 {
        panic(fmt.Errorf("failed to read private key folder: %w", err))
    }

    // 拼出 keystore 資料夾內唯一的檔案路徑
    keyFile := keyPath + "/" + files[0].Name()

    // 讀取私鑰檔案內容
    keyPEM, err := os.ReadFile(keyFile)
    if err != nil {
        panic(fmt.Errorf("failed to read private key: %w", err))
    }

    privateKey, err := identity.PrivateKeyFromPEM(keyPEM)
    if err != nil {
        panic(err)
    }

    sign, err := identity.NewPrivateKeySign(privateKey)
    if err != nil {
        panic(err)
    }

    return sign
}

// 顯示 Gateway 錯誤詳細資訊
func PrintGatewayErrorDetails(err error) {
	var endorseErr *client.EndorseError
	var submitErr *client.SubmitError
	var commitStatusErr *client.CommitStatusError
	var commitErr *client.CommitError

	if errors.As(err, &endorseErr) {
		fmt.Printf("❌ Endorse error [txID: %s], gRPC status: %v\n→ %s\n", endorseErr.TransactionID, status.Code(endorseErr), endorseErr)
	} else if errors.As(err, &submitErr) {
		fmt.Printf("❌ Submit error [txID: %s], gRPC status: %v\n→ %s\n", submitErr.TransactionID, status.Code(submitErr), submitErr)
	} else if errors.As(err, &commitStatusErr) {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("⏱️ Commit timeout [txID: %s]: %s\n", commitStatusErr.TransactionID, commitStatusErr)
		} else {
			fmt.Printf("❌ CommitStatus error [txID: %s], gRPC status: %v\n→ %s\n", commitStatusErr.TransactionID, status.Code(commitStatusErr), commitStatusErr)
		}
	} else if errors.As(err, &commitErr) {
		fmt.Printf("❌ Commit failed [txID: %s], status: %d\n→ %s\n", commitErr.TransactionID, int32(commitErr.Code), err)
	} else {
		fmt.Printf("❓ Unexpected error type %T: %v\n", err, err)
	}

	// 額外印出 endorsement 細節
	statusErr := status.Convert(err)
	for _, d := range statusErr.Details() {
		if detail, ok := d.(*gateway.ErrorDetail); ok {
			fmt.Printf("🔍 Peer: %s\n🏢 MSP: %s\n💬 Msg: %s\n", detail.Address, detail.MspId, detail.Message)
		}
	}
}