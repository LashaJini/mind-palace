package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPC_SERVER_PORT int
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	grpcServerPort, _ := strconv.Atoi(os.Getenv("PYTHON_GRPC_SERVER_PORT"))

	return &Config{
		GRPC_SERVER_PORT: grpcServerPort,
	}
}
