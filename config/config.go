package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/lashajini/mind-palace/constants"
)

// developer config
type Config struct {
	// grpc server
	GRPC_SERVER_PORT int

	// database
	DB_USER string
	DB_PASS string
	DB_NAME string
	DB_PORT int
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	grpcServerPort, _ := strconv.Atoi(os.Getenv("PYTHON_GRPC_SERVER_PORT"))
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	return &Config{
		GRPC_SERVER_PORT: grpcServerPort,

		DB_USER: dbUser,
		DB_PASS: dbPass,
		DB_NAME: dbName,
		DB_PORT: dbPort,
	}
}

func MindPalacePath() string {
	userHome, _ := os.UserHomeDir()
	mindPalaceRoot := filepath.Join(userHome, constants.MIND_PALACE_ROOT)

	return mindPalaceRoot
}

func MindPalaceUserPath(user string) string {
	return filepath.Join(MindPalacePath(), user)
}

func MindPalaceInfoPath() string {
	return filepath.Join(MindPalacePath(), constants.MIND_PALACE_INFO)
}

func MindPalaceOriginalResourcePath(user string) string {
	return filepath.Join(MindPalaceUserPath(user), constants.MIND_PALACE_RESOURCES, constants.MIND_PALACE_ORIGINAL)
}

func MindPalaceMemoryPath(user string) string {
	return filepath.Join(MindPalaceUserPath(user), constants.MIND_PALACE_MEMORIES)
}

func MindPalaceUserConfigPath(user string) string {
	return filepath.Join(MindPalaceUserPath(user), constants.MIND_PALACE_CONFIG)
}
