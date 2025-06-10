package fabric

import (
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
)

// GWBuilder 可用任意錢包身分產生 Gateway + Contract

type GWBuilder struct {
	Peer    *Peer  // 共用 gRPC 連線
	Channel string // 頻道名
	CCName  string // 合約名
}

// NewContract 依身份建立即時 Gateway，回 Contract 與 Gateway
func (b GWBuilder) NewContract(id *identity.X509Identity, signer identity.Sign) (*client.Contract, *client.Gateway, error) {
	gw, err := client.Connect(
		id,
		client.WithSign(signer),
		client.WithClientConnection(b.Peer.Conn()),
		client.WithEvaluateTimeout(10*time.Second),
		client.WithEndorseTimeout(30*time.Second),
		client.WithSubmitTimeout(30*time.Second),
		client.WithCommitStatusTimeout(2*time.Minute),
	)
	if err != nil {
		return nil, nil, err
	}
	ctr := gw.GetNetwork(b.Channel).GetContract(b.CCName)
	return ctr, gw, nil
}
