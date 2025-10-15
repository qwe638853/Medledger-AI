package fabric

import (
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
)

/**
 * @notice Gateway 建構器
 * @dev 封裝 Peer 連線、頻道與合約名稱，用任意身份快速建立 Gateway 與 Contract
 */

type GWBuilder struct {
	Peer    *Peer  // 共用 gRPC 連線
	Channel string // 頻道名
	CCName  string // 合約名
}

/**
 * @notice 依身份建立即時 Gateway，並回傳合約物件
 * @dev 設定 Evaluate/Endorse/Submit/CommitStatus 的逾時時間
 * @param id X.509 身份
 * @param signer 簽章器
 * @return *client.Contract 合約, *client.Gateway Gateway, error 連線或建立失敗
 */
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
