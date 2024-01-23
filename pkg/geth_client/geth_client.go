package geth_client

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"golang-ethereum-example-api/pkg/logging"
	"golang-ethereum-example-api/pkg/settings"
)

var client *ethclient.Client

func Setup(ethereumClient *settings.EthereumClient, logger *logging.LogWrapper) {
	c, err := ethclient.Dial(ethereumClient.Url)
	if err != nil {
		logger.Fatal("Failed to connect to the Ethereum client", zap.Error(err))
	}
	client = c
}

func GetClient() *ethclient.Client {
	return client
}
