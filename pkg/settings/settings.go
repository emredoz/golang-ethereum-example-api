package settings

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type App struct {
	AppName     string `validate:"required"`
	LogLevel    string `validate:"required"`
	Hostname    string `validate:"required"`
	LogEncoding string `validate:"required"`
}

var AppSettings = &App{}

type Server struct {
	HttpPort     int           `validate:"required"`
	ReadTimeout  time.Duration `validate:"required"`
	WriteTimeout time.Duration `validate:"required"`
	RunMode      string        `validate:"required"`
}

var ServerSettings = &Server{}

type EthereumClient struct {
	Url      string `validate:"required"`
	GasLimit uint64 `validate:"required"`
}

var EthereumClientSettings = &EthereumClient{}

func Setup() {
	_ = godotenv.Load()
	validate := validator.New()

	AppSettings.Hostname, _ = os.Hostname()
	AppSettings.LogLevel = os.Getenv("LOG_LEVEL")
	AppSettings.AppName = os.Getenv("APP_NAME")
	AppSettings.LogEncoding = os.Getenv("LOG_ENCODING")
	err := validate.Struct(AppSettings)
	if err != nil {
		log.Fatalf("AppSettings settings missing err: %v", err)
	}

	EthereumClientSettings.Url = os.Getenv("ETHEREUM_URL")
	gasLimitStr := os.Getenv("ETHEREUM_GAS_LIMIT")
	gasLimit, err := strconv.Atoi(gasLimitStr)
	if err != nil {
		log.Fatal("ETHEREUM_GAS_LIMIT setting is not proper err: %v", err)
	}
	EthereumClientSettings.GasLimit = uint64(gasLimit)
	err = validate.Struct(EthereumClientSettings)
	if err != nil {
		log.Fatal("EthereumClient settings missing err: %v", err)
	}

	ServerSettings.HttpPort, _ = strconv.Atoi(os.Getenv("HTTP_PORT"))
	readTimeoutStr := os.Getenv("READ_TIMEOUT")
	ReadTimeout, err := strconv.Atoi(readTimeoutStr)
	if err != nil {
		log.Fatal("READ_TIMEOUT setting is not proper err: %v", err)
	}
	ServerSettings.ReadTimeout = time.Duration(ReadTimeout * 1000000000)
	writeTimeoutStr := os.Getenv("WRITE_TIMEOUT")
	WriteTimeout, err := strconv.Atoi(writeTimeoutStr)
	if err != nil {
		log.Fatal("WRITE_TIMEOUT setting is not proper err: %v", err)
	}
	ServerSettings.WriteTimeout = time.Duration(WriteTimeout * 1000000000)
	ServerSettings.HttpPort, _ = strconv.Atoi(os.Getenv("HTTP_PORT"))
	ServerSettings.RunMode = os.Getenv("RUN_MODE")
	err = validate.Struct(ServerSettings)
	if err != nil {
		log.Fatal("Server settings missing err: %v", err)
	}
}
