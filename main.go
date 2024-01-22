package main

import (
	"fmt"
	ethereumClient "golang-ethereum-example-api/pkg/geth_client"
	"golang-ethereum-example-api/pkg/logging"
	"golang-ethereum-example-api/pkg/settings"
)

func init() {
	settings.Setup()
	logConfig := logging.LogConfig{
		Name:     settings.AppSettings.AppName,
		Level:    settings.AppSettings.LogLevel,
		RunMode:  settings.ServerSettings.RunMode,
		Hostname: settings.AppSettings.Hostname,
		Encoding: settings.AppSettings.LogEncoding,
	}
	logging.Setup(logConfig)
	ethereumClient.Setup(settings.EthereumClientSettings)
}

func main() {
	fmt.Println("Hello")
}
