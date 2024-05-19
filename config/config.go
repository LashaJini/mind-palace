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

	// vector database
	VDB_HOST string
	VDB_NAME string
	VDB_PORT int
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	mindPalaceUser, err := CurrentUser()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	grpcServerPort, _ := strconv.Atoi(os.Getenv("PYTHON_GRPC_SERVER_PORT"))

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := mindPalaceUser + os.Getenv("DB_NAME")
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	vdbHost := os.Getenv("VDB_HOST")
	vdbName := mindPalaceUser + os.Getenv("VDB_NAME")
	vdbPort, _ := strconv.Atoi(os.Getenv("VDB_PORT"))

	return &Config{
		GRPC_SERVER_PORT: grpcServerPort,

		DB_USER: dbUser,
		DB_PASS: dbPass,
		DB_NAME: dbName,
		DB_PORT: dbPort,

		VDB_HOST: vdbHost,
		VDB_NAME: vdbName,
		VDB_PORT: vdbPort,
	}
}

func (c *Config) VDBAddr() string {
	return fmt.Sprintf("%s:%d", c.VDB_HOST, c.VDB_PORT)
}

func (c *Config) DBAddr() string {
	return fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable", c.DB_USER, c.DB_PASS, c.DB_PORT, c.DB_NAME)
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
