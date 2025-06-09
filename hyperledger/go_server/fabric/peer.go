package fabric

import (
	"crypto/x509"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Peer 封裝對單一 peer 節點的 gRPC 連線
// 只需建立一次，可重複給多個 Gateway 使用

type Peer struct {
	conn *grpc.ClientConn
}

// NewPeer 讀 TLS 憑證並連線
func NewPeer(endpoint, tlsPath, hostOverride string) (*Peer, error) {
	pem, err := os.ReadFile(tlsPath)
	if err != nil {
		return nil, fmt.Errorf("read TLS cert: %w", err)
	}
	cert, err := identity.CertificateFromPEM(pem)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AddCert(cert)

	creds := credentials.NewClientTLSFromCert(pool, hostOverride)
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	return &Peer{conn: conn}, nil
}

func (p *Peer) Conn() *grpc.ClientConn { return p.conn }
