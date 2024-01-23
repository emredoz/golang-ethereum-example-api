package controller

import (
	"github.com/gin-gonic/gin"
	"golang-ethereum-example-api/serializers"
	"golang-ethereum-example-api/services"
	"net/http"
)

type AccountController struct {
	Service services.AccountService
}
type AccountControllerConfig struct {
	R       *gin.Engine
	Service services.AccountService
}

func NewAccountController(c *AccountControllerConfig) {
	accountController := &AccountController{
		Service: c.Service,
	}

	api := c.R.Group("/api/v1")
	api.GET("/account/:address/balance", accountController.GetBalance)
	api.POST("/account", accountController.CreateAccount)
}

func (s *AccountController) GetBalance(c *gin.Context) {
	serializer := serializers.Serializer{C: c}
	ctx := c.Request.Context()
	var request serializers.GetBalanceRequest

	errorInfo := serializer.ShouldBindUri(&request)
	if errorInfo != nil {
		serializer.ErrorResponse(errorInfo)
		return
	}

	errorInfo = request.Validate(ctx)
	if errorInfo != nil {
		serializer.ErrorResponse(errorInfo)
		return
	}

	response, errInfo := s.Service.GetBalance(ctx, request)
	if errInfo != nil {
		serializer.ErrorResponse(errInfo)
		return
	}

	serializer.SuccessfulResponse(http.StatusOK, response)
}

func (s *AccountController) CreateAccount(c *gin.Context) {
	serializer := serializers.Serializer{C: c}
	response, errInfo := s.Service.CreateAccount()
	if errInfo != nil {
		serializer.ErrorResponse(errInfo)
		return
	}

	serializer.SuccessfulResponse(http.StatusOK, response)
}
