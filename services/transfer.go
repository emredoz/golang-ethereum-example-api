package services

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"golang-ethereum-example-api/pkg/logging"
	"golang-ethereum-example-api/pkg/settings"
	"golang-ethereum-example-api/pkg/util"
	"golang-ethereum-example-api/serializers"
	"math/big"
	"net/http"
)

type TransferService interface {
	SendEthereum(ctx context.Context, request serializers.SendEthereumRequest) (*serializers.SendEthereumResponse, *util.ErrorInfo)
}

type transferService struct {
	client *ethclient.Client
	config *settings.EthereumClient
	logger *logging.LogWrapper
}

func NewTransferService(client *ethclient.Client, config *settings.EthereumClient, logger *logging.LogWrapper) TransferService {
	return &transferService{client: client, config: config, logger: logger}
}

func (s *transferService) SendEthereum(ctx context.Context, request serializers.SendEthereumRequest) (*serializers.SendEthereumResponse, *util.ErrorInfo) {
	fromAccount := common.HexToAddress(request.FromAddress)
	toAccount := common.HexToAddress(request.ToAddress)
	privateKey, err := crypto.HexToECDSA(request.PrivateKey[2:])
	if err != nil {
		s.logger.Error("TransferEthereum privateKey to ECDSA error", zap.Error(err))
		return nil, &util.ErrorInfo{
			HttpCode: http.StatusInternalServerError,
			Message:  util.InternalServiceErrorMessage,
			Err:      err,
		}
	}
	nonce, err := s.client.PendingNonceAt(ctx, fromAccount)
	if err != nil {
		s.logger.Error("TransferEthereum getting nonce error", zap.Error(err))
		return nil, &util.ErrorInfo{
			HttpCode: http.StatusInternalServerError,
			Message:  util.InternalServiceErrorMessage,
			Err:      err,
		}
	}
	amount, err := etherToWei(request.EthereumAmount)
	if err != nil {
		s.logger.Error("TransferEthereum etherToWei error", zap.Error(err))
		return nil, &util.ErrorInfo{
			HttpCode: http.StatusInternalServerError,
			Message:  util.InternalServiceErrorMessage,
			Err:      err,
		}
	}

	gasLimit := s.config.GasLimit

	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		s.logger.Error("TransferEthereum suggestGasPrice error", zap.Error(err))
		return nil, &util.ErrorInfo{
			HttpCode: http.StatusInternalServerError,
			Message:  util.InternalServiceErrorMessage,
			Err:      err,
		}
	}

	txData := &types.LegacyTx{
		Nonce:    nonce,
		To:       &toAccount,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     nil,
	}
	tx := types.NewTx(txData)

	/*		During the development phase, I utilized Ganache.
			I encountered an error with this chainId in Ganache.
			The detailed information regarding the error can be found here: https://github.com/trufflesuite/ganache/issues/4367.
			Therefore, I am currently obtaining the chainId statically due to this issue.
	chainID, err := s.client.NetworkID(ctx)
	if err != nil {
		s.logger.Error("TransferEthereum getting chainID error", zap.Error(err))
		return nil, &util.ErrorInfo{
			HttpCode: http.StatusInternalServerError,
			Message:  util.InternalServiceErrorMessage,
			Err:      err,
		}
	}
	*/
	chainID := new(big.Int).SetInt64(int64(1337))

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		s.logger.Error("TransferEthereum sign transaction error", zap.Error(err))
		return nil, &util.ErrorInfo{
			HttpCode: http.StatusInternalServerError,
			Message:  util.InternalServiceErrorMessage,
			Err:      err,
		}
	}

	err = s.client.SendTransaction(ctx, signedTx)
	if err != nil {
		s.logger.Warn("TransferEthereum send transaction error", zap.Error(err))
		return nil, &util.ErrorInfo{
			HttpCode: http.StatusInternalServerError,
			Message:  util.InternalServiceErrorMessage,
			Err:      err,
		}
	}
	response := &serializers.SendEthereumResponse{
		TransactionHash: signedTx.Hash().Hex(),
	}
	return response, nil
}

func etherToWei(ether float64) (*big.Int, error) {
	weiFloat := new(big.Float).SetFloat64(ether * 1e18)
	wei, err := new(big.Int).SetString(weiFloat.Text('f', 0), 10)
	if !err {
		return nil, errors.New("ether to wei convert error")
	}
	return wei, nil
}
