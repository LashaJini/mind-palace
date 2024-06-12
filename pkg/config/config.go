package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/lashajini/mind-palace/cli/errors"
	"github.com/lashajini/mind-palace/pkg/constants"
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
	errors.Handle(err)

	mindPalaceUser, err := CurrentUser()
	errors.Handle(err)

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

func MindPalacePath(homePrefix bool) string {
	userHome := ""
	if homePrefix {
		userHome, _ = os.UserHomeDir()
	}
	mindPalaceRoot := filepath.Join(userHome, constants.MIND_PALACE_ROOT)

	return mindPalaceRoot
}

func UserPath(user string, homePrefix bool) string {
	return filepath.Join(MindPalacePath(homePrefix), user)
}

func InfoPath(homePrefix bool) string {
	return filepath.Join(MindPalacePath(homePrefix), constants.MIND_PALACE_INFO)
}

func OriginalResourcePath(user string, homePrefix bool) string {
	return filepath.Join(UserPath(user, homePrefix), constants.MIND_PALACE_RESOURCES, constants.MIND_PALACE_ORIGINAL)
}

// relative to user home dir
func OriginalResourceRelativePath(user string) string {
	return filepath.Join(UserPath(user, false), constants.MIND_PALACE_RESOURCES, constants.MIND_PALACE_ORIGINAL)
}

// including user home dir
func OriginalResourceFullPath(user string) string {
	return filepath.Join(UserPath(user, true), constants.MIND_PALACE_RESOURCES, constants.MIND_PALACE_ORIGINAL)
}

func MemoryPath(user string, homePrefix bool) string {
	return filepath.Join(UserPath(user, homePrefix), constants.MIND_PALACE_MEMORIES)
}

func UserConfigPath(user string, homePrefix bool) string {
	return filepath.Join(UserPath(user, homePrefix), constants.MIND_PALACE_CONFIG)
}
