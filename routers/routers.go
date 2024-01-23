package routers

import (
	"github.com/gin-gonic/gin"
	"golang-ethereum-example-api/controller"
	ethereumClient "golang-ethereum-example-api/pkg/geth_client"
	"golang-ethereum-example-api/pkg/logging"
	"golang-ethereum-example-api/pkg/settings"
	"golang-ethereum-example-api/services"
)

func BuildRouter() *gin.Engine {
	router := newRouter()
	logger := logging.GetLogger()

	accountService := services.NewAccountService(ethereumClient.GetClient(), logger)
	controller.NewAccountController(&controller.AccountControllerConfig{
		R: router, Service: accountService})

	transferService := services.NewTransferService(ethereumClient.GetClient(), settings.EthereumClientSettings, logger)
	controller.NewTransferController(&controller.TransferControllerConfig{
		R: router, Service: transferService,
	})
	return router
}

func newRouter() *gin.Engine {
	r := gin.New()
	return r
}
