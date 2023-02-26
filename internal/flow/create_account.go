package flow

import (
	"context"

	"ella.wallet-backend/internal/config"
	"ella.wallet-backend/internal/models"
	"ella.wallet-backend/internal/queue"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/templates"
)

type AccountObjext struct {
	TxId *flow.Identifier
	Job  *queue.Job
}

func (fc *FlowClient) CreateAccount(keys models.CreateWalletRequest) (*flow.Identifier, error) {
	ctx := context.Background()

	// Get Backend Account
	adminAddr := flow.HexToAddress(config.AppConfig.FlowBackendAddr)
	adminAccount, err := fc.AccessNode.GetAccount(ctx, adminAddr)
	if err != nil {
		return nil, err
	}

	// Decode Admin Keys
	adminAccountPublicKey := adminAccount.Keys[0]                                                                             // TODO: This should not be hardcoded here
	adminAccountPrivateKey, err := crypto.DecodePrivateKeyHex(adminAccountPublicKey.SigAlgo, config.AppConfig.FlowBackendKey) // TODO: This needs to be updated to use GCP for key management
	if err != nil {
		return nil, err
	}

	// Create Admin Signer
	AdminSigner, err := crypto.NewInMemorySigner(adminAccountPrivateKey, adminAccountPublicKey.HashAlgo)
	if err != nil {
		return nil, err
	}

	// convert new account keys
	var flowAccountKeys []*flow.AccountKey
	recoveryKey, err := createFlowAccountKey(keys.RecoveryPublicKey, 1000)
	if err != nil {
		return nil, err
	}
	flowAccountKeys = append(flowAccountKeys, recoveryKey)

	accountKey, err := createFlowAccountKey(keys.AccountPublicKey, 500)
	if err != nil {
		return nil, err
	}
	flowAccountKeys = append(flowAccountKeys, accountKey)

	deviceKey, err := createFlowAccountKey(keys.DevicePublicKey, 500)
	if err != nil {
		return nil, err
	}
	flowAccountKeys = append(flowAccountKeys, deviceKey)

	// generate an account creation createAccountTx
	createAccountTx, err := templates.CreateAccount(flowAccountKeys, nil, adminAddr)
	if err != nil {
		return nil, err
	}

	latestBlock, err := fc.AccessNode.GetLatestBlock(ctx, true)
	if err != nil {
		return nil, err
	}

	tx := createAccountTx.SetGasLimit(100).SetProposalKey(adminAddr, adminAccountPublicKey.Index, adminAccountPublicKey.SequenceNumber).SetPayer(adminAddr).SetReferenceBlockID(latestBlock.ID)

	err = tx.SignEnvelope(adminAddr, adminAccountPublicKey.Index, AdminSigner)
	if err != nil {
		return nil, err
	}

	err = fc.AccessNode.SendTransaction(ctx, *tx)
	if err != nil {
		return nil, err
	}

	txId := tx.ID()

	return &txId, nil
}

func createFlowAccountKey(key string, threshold int) (*flow.AccountKey, error) {
	// convert recovery key hex to crypto key
	pk, err := crypto.DecodePublicKeyHex(crypto.ECDSA_P256, key)
	if err != nil {
		return nil, err
	}

	// construct an account key from the public key
	accountKey := flow.NewAccountKey().
		SetPublicKey(pk).
		SetHashAlgo(crypto.SHA2_256). // pair this key with the SHA3_256 hashing algorithm
		SetWeight(threshold)          // give this key full signing weight

	return accountKey, nil
}
