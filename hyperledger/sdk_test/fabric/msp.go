package fabric

import (
	"fmt"
	"log"

	// fabric sdk
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	// 補上 contextAPI 別名
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
)

// 註冊並 Enroll 用戶
func RegisterNewUser(adminContextProvider contextAPI.ClientProvider, username, secret, affiliation string) error {
	mspClient, err := msp.New(adminContextProvider)
	if err != nil {
		return fmt.Errorf("❌ 建立 MSP client 失敗: %w", err)
	}

	_, err = mspClient.Register(&msp.RegistrationRequest{
		Name:        username,
		Affiliation: affiliation,
		Secret:     secret,
	})
	if err != nil {
		return fmt.Errorf("❌ 註冊失敗: %w", err)
	}

	err = mspClient.Enroll(username, msp.WithSecret(secret))
	if err != nil {
		return fmt.Errorf("❌ Enroll 失敗: %w", err)
	}

	log.Printf("✅ 用戶 %s 註冊與 Enroll 成功", username)
	return nil
}

// 建立 MSP client
func NewMspClient(sdk *fabsdk.FabricSDK, orgName string) (*msp.Client, error) {
	ctxProvider := sdk.Context(fabsdk.WithOrg(orgName))
	return msp.New(ctxProvider)
}

// 初始化 Fabric SDK
func NewSDK(configPath string) (*fabsdk.FabricSDK, error) {
	sdk, err := fabsdk.New(config.FromFile(configPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create Fabric SDK: %w", err)
	}
	return sdk, nil
}
