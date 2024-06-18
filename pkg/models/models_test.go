package models

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/mpuser"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ModelsTestSuite struct {
	suite.Suite
	pgContainer      *postgres.PostgresContainer
	mpuser           string
	currentSchema    string
	tmpMigrationsDir string
	db               *database.MindPalaceDB
	cfg              *common.Config
	ctx              context.Context
}

func CreatePostgresContainer(ctx context.Context, cfg *common.Config, migrationFiles []string) (*postgres.PostgresContainer, error) {
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

// create temporary migrations dir for up migrations and inject into sql templates
func (suite *ModelsTestSuite) initMigrationFiles() []string {
	pattern := filepath.Join(suite.cfg.MIGRATIONS_DIR, "*.up.sql")
	migrationFiles, _ := filepath.Glob(pattern)
	suite.tmpMigrationsDir = filepath.Join(suite.cfg.MIGRATIONS_DIR, "..", "tmp_test_migrations")
	os.Mkdir(suite.tmpMigrationsDir, 0755)

	suite.currentSchema = suite.mpuser + suite.cfg.DB_SCHEMA_SUFFIX
	sqlTemplate := common.SQLTemplate{Namespace: suite.currentSchema}
	for _, migrationFile := range migrationFiles {
		var sqlBuffer bytes.Buffer
		err := sqlTemplate.Inject(&sqlBuffer, migrationFile)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(filepath.Join(suite.tmpMigrationsDir, filepath.Base(migrationFile)), sqlBuffer.Bytes(), 0644)
	}
	schemaQuery := fmt.Sprintf("CREATE SCHEMA %s", suite.currentSchema)
	err := os.WriteFile(filepath.Join(suite.tmpMigrationsDir, "000_init_user_schema.up.sql"), []byte(schemaQuery), 0644)
	if err != nil {
		log.Fatal(err)
	}

	pattern = filepath.Join(suite.tmpMigrationsDir, "*.up.sql")
	migrationFiles, _ = filepath.Glob(pattern)

	return migrationFiles
}

func (suite *ModelsTestSuite) SetupSuite() {
	suite.mpuser = common.TEST_USER
	err := mpuser.CreateMindPalace(suite.mpuser)
	if err != nil {
		log.Fatal(err)
	}
	err = common.UpdateMindPalaceInfo(common.MindPalaceInfo{CurrentUser: suite.mpuser})
	if err != nil {
		log.Fatal(err)
	}

	suite.cfg = common.NewConfig()
	suite.ctx = context.Background()

	pgContainer, err := CreatePostgresContainer(suite.ctx, suite.cfg, suite.initMigrationFiles())
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
	suite.db.SetSchema(suite.currentSchema)
}

func (suite *ModelsTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("failed to terminate pgContainer: %s", err)
	}

	err := os.RemoveAll(common.UserPath(suite.mpuser, true))
	if err != nil {
		log.Fatalf("failed to remove user mind palace: %s", err)
	}

	err = common.UpdateMindPalaceInfo(common.MindPalaceInfo{CurrentUser: ""})
	if err != nil {
		log.Fatalf("failed to update config: %s", err)
	}

	err = os.RemoveAll(suite.tmpMigrationsDir)
	if err != nil {
		log.Fatalf("failed to remove %s: %s", suite.tmpMigrationsDir, err)
	}
}

func TestModelsTestSuite(t *testing.T) {
	suite.Run(t, new(ModelsTestSuite))
}
