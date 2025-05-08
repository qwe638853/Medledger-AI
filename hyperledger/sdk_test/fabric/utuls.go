package fabric

import (
	"context"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc/status"
)

// PrintGatewayError 詳細列印 Fabric Gateway 的錯誤
func PrintGatewayError(err error) {
	var (
		endorseErr      *client.EndorseError
		submitErr       *client.SubmitError
		commitStatusErr *client.CommitStatusError
		commitErr       *client.CommitError
	)

	switch {
	case errors.As(err, &endorseErr):
		fmt.Printf("❌ Endorse error [txID:%s] %v\n", endorseErr.TransactionID, endorseErr)
	case errors.As(err, &submitErr):
		fmt.Printf("❌ Submit error [txID:%s] %v\n", submitErr.TransactionID, submitErr)
	case errors.As(err, &commitStatusErr):
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("⏱️ Commit timeout [txID:%s]\n", commitStatusErr.TransactionID)
		}
	case errors.As(err, &commitErr):
		fmt.Printf("❌ Commit failed [txID:%s] status:%d %v\n", commitErr.TransactionID, commitErr.Code, commitErr)
	default:
		fmt.Printf("❓ Unexpected error %v\n", err)
	}

	// 印 endorsement 細節
	if s, ok := status.FromError(err); ok {
		for _, d := range s.Details() {
			if det, ok := d.(*gateway.ErrorDetail); ok {
				fmt.Printf("🔍 Peer:%s MSP:%s → %s\n", det.Address, det.MspId, det.Message)
			}
		}
	}
}
