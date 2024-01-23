package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	ethereumClient "golang-ethereum-example-api/pkg/geth_client"
	"golang-ethereum-example-api/pkg/logging"
	"golang-ethereum-example-api/pkg/settings"
	"golang-ethereum-example-api/routers"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	ethereumClient.Setup(settings.EthereumClientSettings, logging.GetLogger())
}

func main() {
	logger := logging.GetLogger()
	gin.SetMode(settings.ServerSettings.RunMode)
	if settings.ServerSettings.RunMode == gin.ReleaseMode {
		gin.DefaultWriter = io.Discard
	}

	router := routers.BuildRouter()
	readTimeout := settings.ServerSettings.ReadTimeout
	writeTimeout := settings.ServerSettings.WriteTimeout
	endPoint := fmt.Sprintf(":%d", settings.ServerSettings.HttpPort)
	server := &http.Server{
		Addr:           endPoint,
		Handler:        router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Server failed to start: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		_ = logger.Sync()
	}()
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	}
}
