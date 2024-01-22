package geth_client

import (
	"golang-ethereum-example-api/pkg/settings"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

var client *ethclient.Client

// todo:logger ile patlat
func Setup(ethereumClient *settings.EthereumClient) {
	c, err := ethclient.Dial(ethereumClient.Url)
	if err != nil {
		log.Fatal("Failed to connect to the Ethereum client: %v", err)
	}
	client = c
}

func GetClient() *ethclient.Client {
	return client
}
