package flow

import (
	"ella.wallet-backend/internal/config"

	"github.com/onflow/flow-go-sdk/access/http"
)

type FlowClient struct {
	accessNode *http.Client
}

func NewClient() (*FlowClient, error) {
	var flowClient *http.Client
	var flowClientErr error

	switch config.AppConfig.FlowNetwork {
	case "emulator":
		flowClient, flowClientErr = http.NewClient(http.EmulatorHost)
	case "testnet":
		flowClient, flowClientErr = http.NewClient(http.TestnetHost)
	case "mainnet":
		flowClient, flowClientErr = http.NewClient(http.MainnetHost)
	}

	if flowClientErr != nil {
		return nil, flowClientErr
	}

	client := &FlowClient{
		accessNode: flowClient,
	}

	return client, nil
}
