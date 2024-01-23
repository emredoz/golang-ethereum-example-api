package controller

import (
	"github.com/gin-gonic/gin"
	"golang-ethereum-example-api/serializers"
	"golang-ethereum-example-api/services"
	"net/http"
)

type TransferController struct {
	Service services.TransferService
}
type TransferControllerConfig struct {
	R       *gin.Engine
	Service services.TransferService
}

func NewTransferController(c *TransferControllerConfig) {
	accountController := &TransferController{
		Service: c.Service,
	}

	api := c.R.Group("/api/v1")
	api.POST("/transfer/send", accountController.SendEthereum)
}

func (s *TransferController) SendEthereum(c *gin.Context) {
	serializer := serializers.Serializer{C: c}
	ctx := c.Request.Context()
	var request serializers.SendEthereumRequest
	errInfo := serializer.ShouldBindJSON(&request)
	if errInfo != nil {
		serializer.ErrorResponse(errInfo)
		return
	}

	errorInfo := request.Validate(ctx)
	if errorInfo != nil {
		serializer.ErrorResponse(errorInfo)
		return
	}

	response, errInfo2 := s.Service.SendEthereum(ctx, request)
	if errInfo2 != nil {
		serializer.ErrorResponse(errInfo2)
		return
	}
	serializer.SuccessfulResponse(http.StatusOK, response)
}
