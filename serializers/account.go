package serializers

import (
	"context"
	"errors"
	"golang-ethereum-example-api/pkg/util"
	"math/big"
	"net/http"
)

type GetBalanceRequest struct {
	Address string `uri:"address" validate:"required"`
}

type GetBalanceResponse struct {
	Address  string     `json:"address,omitempty"`
	EthValue *big.Float `json:"ethValue,omitempty"`
}

func (r *GetBalanceRequest) Validate(ctx context.Context) *util.ErrorInfo {
	errInfo := validate(ctx, r)
	if errInfo != nil {
		return errInfo
	}
	isValid := addressValidationRegex.MatchString(r.Address)
	if !isValid {
		return &util.ErrorInfo{
			HttpCode: http.StatusBadRequest,
			Message:  util.InvalidAccountAddressErrorMessage,
			Err:      errors.New("invalid address"),
		}
	}
	return nil
}

type CreateAccountResponse struct {
	Address    string `json:"accountId,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
}
