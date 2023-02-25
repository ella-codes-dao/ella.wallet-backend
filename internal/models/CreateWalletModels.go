package models

type CreateWalletRequest struct {
	RecoveryPublicKey string `json:"recoveryPublicKey" binding:"required"`
	AccountPublicKey  string `json:"accountPublicKey" binding:"required"`
	DevicePublicKey   string `json:"devicePublicKey" binding:"required"`
}

type CreateWalletResponse struct {
	Status string `json:"status"`
	MSG    string `json:"msg,ommitempty"`
	TxId   string `json:"txId,ommitempty"`
}
