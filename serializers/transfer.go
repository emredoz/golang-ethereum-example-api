package serializers

import (
	"context"
	"errors"
	"golang-ethereum-example-api/pkg/util"
	"net/http"
)

type SendEthereumRequest struct {
	FromAddress    string  `json:"fromAddress" validate:"required"`
	PrivateKey     string  `json:"privateKey" validate:"required"`
	ToAddress      string  `json:"toAddress" validate:"required"`
	EthereumAmount float64 `json:"ethereumAmount" validate:"required"`
}

func (r *SendEthereumRequest) Validate(ctx context.Context) *util.ErrorInfo {
	errInfo := validate(ctx, r)
	if errInfo != nil {
		return errInfo
	}
	isValidToAddress := addressValidationRegex.MatchString(r.ToAddress)
	isValidFromAddress := addressValidationRegex.MatchString(r.FromAddress)
	if !isValidFromAddress || !isValidToAddress {
		return &util.ErrorInfo{
			HttpCode: http.StatusBadRequest,
			Message:  util.InvalidAccountAddressErrorMessage,
			Err:      errors.New("invalid address"),
		}
	}
	return nil
}

type SendEthereumResponse struct {
	TransactionHash string `json:"transactionHash,omitempty"`
}
