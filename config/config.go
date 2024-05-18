package config

import (
	"log"
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/lashajini/mind-palace/constants"
)

// developer config
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

func UserMindPalaceRoot(user string) string {
	userHome, _ := os.UserHomeDir()
	mindPalaceRoot := path.Join(userHome, constants.MIND_PALACE_ROOT, user)

	return mindPalaceRoot
}
