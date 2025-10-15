package fabric

import (
	"crypto/x509"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/**
 * @notice Peer 連線封裝
 * @dev 封裝對單一 peer 節點的 gRPC 連線，建議全域重用
 */

type Peer struct {
	conn *grpc.ClientConn
}

/**
 * @notice 建立到 peer 的 TLS gRPC 連線
 * @dev 讀取 TLS CA 憑證，建立憑證池，再以指定 SNI 連線
 * @param endpoint peer 位址，例如 localhost:7051
 * @param tlsPath TLS CA 憑證路徑
 * @param hostOverride SNI 主機名稱（通常為 peer FQDN）
 * @return *Peer 連線物件, error 連線失敗
 */
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
