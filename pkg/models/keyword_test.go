package models

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/lashajini/mind-palace/pkg/config"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreatePostgresContainer(ctx context.Context, cfg *config.Config) (*postgres.PostgresContainer, error) {
	pattern := filepath.Join(cfg.MIGRATIONS_DIR, "*.up.sql")
	migrationFiles, _ := filepath.Glob(pattern)

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage(fmt.Sprintf("postgres:%s", cfg.DB_VERSION)),
		postgres.WithInitScripts(migrationFiles...),
		postgres.WithDatabase(cfg.DB_NAME),
		postgres.WithUsername(cfg.DB_USER),
		postgres.WithPassword(cfg.DB_PASS),

		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	return pgContainer, nil
}

type KeywordsTestSuite struct {
	suite.Suite
	pgContainer *postgres.PostgresContainer
	db          *database.MindPalaceDB
	cfg         *config.Config
	ctx         context.Context
}

// TODO: initialize user first
func (suite *KeywordsTestSuite) SetupSuite() {
	suite.cfg = config.NewConfig()
	suite.ctx = context.Background()

	pgContainer, err := CreatePostgresContainer(suite.ctx, suite.cfg)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer

	_, err = pgContainer.ConnectionString(suite.ctx, "sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	port, err := pgContainer.MappedPort(suite.ctx, "5432")
	if err != nil {
		log.Fatal(err)
	}

	suite.cfg.DB_PORT, err = strconv.Atoi(port.Port())
	if err != nil {
		log.Fatal(err)
	}

	suite.db = database.InitDB(suite.cfg)
}

func (suite *KeywordsTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("failed to terminate pgContainer: %s", err)
	}
}

func (suite *KeywordsTestSuite) Test_InsertManyKeywordsTx_success() {
	t := suite.T()
	tx := database.NewMultiInstruction(suite.ctx, suite.db.DB())

	err := tx.Begin()
	assert.NoError(t, err)

	keywords := []string{"keyword1", "keyword2", "keyword3"}
	keywordIDs, _ := InsertManyKeywordsTx(tx, keywords)
	expectedIDs := map[string]int{
		"keyword1": 1,
		"keyword2": 2,
		"keyword3": 3,
	}

	if len(keywordIDs) != len(expectedIDs) {
		t.Fatalf("Unexpected number of keyword IDs. Expected %d, got %d", len(expectedIDs), len(keywordIDs))
	}
	for keyword, id := range expectedIDs {
		if retrievedID, ok := keywordIDs[keyword]; !ok || retrievedID != id {
			t.Fatalf("Unexpected keyword ID for keyword %s. Expected %d, got %d", keyword, id, retrievedID)
		}
	}
}

func TestKeywordsTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordsTestSuite))
}
