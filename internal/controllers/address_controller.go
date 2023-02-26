package controllers

import (
	ctx "context"
	"log"
	"net/http"

	"ella.wallet-backend/internal/models"
	"ella.wallet-backend/internal/queue"
	"ella.wallet-backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/onflow/flow-go-sdk"
)

func GetAddress(context *gin.Context) {

}

func CreateAddress(context *gin.Context) {
	var accountRequest models.CreateWalletRequest
	if err := context.BindJSON(&accountRequest); err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var response models.CreateWalletResponse

	txId, err := services.Flow.CreateAccount(accountRequest)
	if err != nil {
		log.Println(err)
		response.Status = "error"
		response.MSG = "error creating transaction to create new Flow account"
		context.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	monitorTransaction := queue.Job{
		Name: "Monitor Create Account Transaction",
		Action: func() error {
			var errorCount = 0
		checkTx:
			for {
				result, err := services.Flow.AccessNode.GetTransactionResult(ctx.Background(), *txId)
				if err != nil {
					log.Println(err) // TODO: Logging!!!

					errorCount += 1
					if errorCount < 100 { // TODO: This is not a clean way to handle this, the transaction is not immediately available and you have to check several times before it stops erroring out...
						continue
					} else {
						return err
					}
				}

				wallet := models.Wallet{
					PublicKey: accountRequest.RecoveryPublicKey,
				}

				if result.Status == flow.TransactionStatusSealed {
					for _, event := range result.Events {
						if event.Type == flow.EventAccountCreated {
							accountCreatedEvent := flow.AccountCreatedEvent(event)
							wallet.Address = accountCreatedEvent.Address().Hex()
						}
					}

					services.DB.Create(&wallet)

					break checkTx
				}
			}

			return nil
		},
	}

	go services.Queues.Flow.AddJob(monitorTransaction)

	response.Status = "success"
	response.TxId = txId.Hex()

	context.JSON(http.StatusOK, response)
}
