package services

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"golang-ethereum-example-api/pkg/logging"
	"golang-ethereum-example-api/pkg/util"
	"golang-ethereum-example-api/serializers"
	"math"
	"math/big"
	"net/http"
)

type AccountService interface {
	GetBalance(ctx context.Context, request serializers.GetBalanceRequest) (*serializers.GetBalanceResponse, *util.ErrorInfo)
	CreateAccount() (*serializers.CreateAccountResponse, *util.ErrorInfo)
}

type accountService struct {
	client *ethclient.Client
	logger *logging.LogWrapper
}

func NewAccountService(client *ethclient.Client, logger *logging.LogWrapper) AccountService {
	return &accountService{client: client, logger: logger}
}

func (s *accountService) GetBalance(ctx context.Context, request serializers.GetBalanceRequest) (*serializers.GetBalanceResponse, *util.ErrorInfo) {
	address := common.HexToAddress(request.Address)
	balance, err := s.client.BalanceAt(ctx, address, nil)
	if err != nil {
		s.logger.Error("GetBalance getting balance error", zap.Error(err), zap.String("address", request.Address))
		return nil, &util.ErrorInfo{
			HttpCode: http.StatusInternalServerError,
			Message:  util.InternalServiceErrorMessage,
			Err:      err,
		}
	}

	var fbalance = new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return &serializers.GetBalanceResponse{
		EthValue: ethValue,
		Address:  request.Address,
	}, nil
}

func (s *accountService) CreateAccount() (*serializers.CreateAccountResponse, *util.ErrorInfo) {
	privateKeyECDSA, err := crypto.GenerateKey()
	if err != nil {
		s.logger.Error("CreateAccount generateKey error", zap.Error(err))
		return nil, &util.ErrorInfo{
			HttpCode: http.StatusInternalServerError,
			Message:  util.InternalServiceErrorMessage,
			Err:      err,
		}
	}
	privateKeyBytes := crypto.FromECDSA(privateKeyECDSA)
	privateKeyStr := hexutil.Encode(privateKeyBytes)

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		s.logger.Error("CreateAccount publicKey convert to ECDSA error", zap.Error(err))
		return nil, &util.ErrorInfo{
			HttpCode: http.StatusInternalServerError,
			Message:  util.InternalServiceErrorMessage,
			Err:      errors.New("public key convert error"),
		}
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return &serializers.CreateAccountResponse{
		Address:    address,
		PrivateKey: privateKeyStr,
	}, nil
}
