package common

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// developer config
type Config struct {
	// grpc server
	PALACE_GRPC_SERVER_PORT int
	VDB_GRPC_SERVER_PORT    int
	LOG_GRPC_SERVER_PORT    int

	// database
	DB_USER              string
	DB_PASS              string
	DB_NAME              string
	DB_PORT              int
	DB_VERSION           string
	DB_DRIVER            string
	DB_DEFAULT_NAMESPACE string
	DB_SCHEMA_SUFFIX     string
	MIGRATIONS_DIR       string

	// vector database
	VDB_HOST string
	VDB_NAME string
	VDB_PORT int
}

func NewConfig() *Config {
	projectRoot := os.Getenv(PROJECT_ROOT)
	env := os.Getenv(MP_ENV)
	if !ENVS[env] {
		env = DEV_ENV
	}

	envFile := filepath.Join(projectRoot, fmt.Sprintf(".env.%s", env))

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	palaceGrpcServerPort, _ := strconv.Atoi(os.Getenv("PALACE_GRPC_SERVER_PORT"))
	vdbGrpcServerPort, _ := strconv.Atoi(os.Getenv("VDB_GRPC_SERVER_PORT"))
	logGrpcServerPort, _ := strconv.Atoi(os.Getenv("LOG_GRPC_SERVER_PORT"))

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	dbVersion := os.Getenv("DB_VERSION")
	dbDriver := os.Getenv("DB_DRIVER")
	dbDefaultNamespace := DB_DEFAULT_SCHEMA
	dbSchemaSuffix := DB_SCHEMA_SUFFIX
	migrationsDir := filepath.Join(projectRoot, os.Getenv("MIGRATIONS_DIR"))

	vdbHost := os.Getenv("VDB_HOST")
	vdbName := os.Getenv("VDB_NAME")
	vdbPort, _ := strconv.Atoi(os.Getenv("VDB_PORT"))

	return &Config{
		PALACE_GRPC_SERVER_PORT: palaceGrpcServerPort,
		VDB_GRPC_SERVER_PORT:    vdbGrpcServerPort,
		LOG_GRPC_SERVER_PORT:    logGrpcServerPort,

		DB_USER:              dbUser,
		DB_PASS:              dbPass,
		DB_NAME:              dbName,
		DB_PORT:              dbPort,
		DB_VERSION:           dbVersion,
		DB_DRIVER:            dbDriver,
		DB_DEFAULT_NAMESPACE: dbDefaultNamespace,
		DB_SCHEMA_SUFFIX:     dbSchemaSuffix,
		MIGRATIONS_DIR:       migrationsDir,

		VDB_HOST: vdbHost,
		VDB_NAME: vdbName,
		VDB_PORT: vdbPort,
	}
}

func (c *Config) VDBAddr() string {
	return fmt.Sprintf("%s:%d", c.VDB_HOST, c.VDB_PORT)
}

func (c *Config) DBAddr() string {
	return fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable&connect_timeout=5", c.DB_USER, c.DB_PASS, c.DB_PORT, c.DB_NAME)
}

func MindPalacePath(homePrefix bool) string {
	userHome := ""
	if homePrefix {
		userHome, _ = os.UserHomeDir()
	}

	var root string
	env := os.Getenv("MP_ENV")
	if env == TEST_ENV {
		root = _MIND_PALACE_TEST_PATH
	}
	root = filepath.Join(MIND_PALACE_ROOT, root)

	mindPalaceRoot := filepath.Join(userHome, root)

	return mindPalaceRoot
}

func UserPath(user string, homePrefix bool) string {
	return filepath.Join(MindPalacePath(homePrefix), user)
}

func InfoPath(homePrefix bool) string {
	return filepath.Join(MindPalacePath(homePrefix), MIND_PALACE_INFO)
}

func OriginalResourcePath(user string, homePrefix bool) string {
	return filepath.Join(UserPath(user, homePrefix), MIND_PALACE_RESOURCES, MIND_PALACE_ORIGINAL)
}

// relative to user home dir
func OriginalResourceRelativePath(user string) string {
	return filepath.Join(UserPath(user, false), MIND_PALACE_RESOURCES, MIND_PALACE_ORIGINAL)
}

// including user home dir
func OriginalResourceFullPath(user string) string {
	return filepath.Join(UserPath(user, true), MIND_PALACE_RESOURCES, MIND_PALACE_ORIGINAL)
}

func MemoryPath(user string, homePrefix bool) string {
	return filepath.Join(UserPath(user, homePrefix), MIND_PALACE_MEMORIES)
}

func UserConfigPath(user string, homePrefix bool) string {
	return filepath.Join(UserPath(user, homePrefix), MIND_PALACE_CONFIG)
}
