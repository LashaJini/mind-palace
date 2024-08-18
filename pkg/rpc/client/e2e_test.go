package rpcclient

import (
	"context"
	"testing"

	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/errors"
	"github.com/lashajini/mind-palace/pkg/mpuser"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

type E2ETestSuite struct {
	suite.Suite
	client *Client
	db     *database.MindPalaceDB
	user   string
	ctx    context.Context
}

// ensures test user exists
// ensures rpc server is up
// inits db
// init vdb is up
func (s *E2ETestSuite) SetupSuite() {
	s.ctx = context.Background()
	defer revert()

	s.user = "test_user"
	userCfg, err := mpuser.CreateNewUser(s.user)
	errors.On(err).Panic()

	cfg := common.NewConfig()

	s.client = NewClient(cfg, userCfg)
	err = s.client.Ping(s.ctx)
	errors.On(err).Panic()

	db := database.InitDB(cfg)
	err = db.CreateSchema(s.user)
	errors.On(err).Panic()

	schemas, err := db.ListMPSchemas()
	errors.On(err).Panic()

	sqlTemplates := []common.SQLTemplate{}
	for _, schema := range schemas {
		sqlTemplate := common.SQLTemplate{
			Namespace: schema,
		}

		sqlTemplates = append(sqlTemplates, sqlTemplate)
	}

	inMemorySource, ups := database.Up(cfg, sqlTemplates)
	database.CommitMigration(cfg, inMemorySource, ups)

	err = s.client.VDBPing(s.ctx)
	errors.On(err).Panic()
}

func (s *E2ETestSuite) Test_default_addon() {
	t := s.T()

	assert.True(t, true)
}

// deletes test user
// drop vdb
func (s *E2ETestSuite) TearDownSuite() {
	common.Log.Info().Msg("starting E2E test teardown")

	err := mpuser.DeleteUser(s.user)
	errors.On(err).Panic()

	common.Log.Info().Msgf("user '%s' mindpalace directory deleted", s.user)

	s.client.VDBDrop(s.ctx)
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}

func revert() {
	if r := recover(); r != nil {
		common.Log.Info().Msgf("RECOVERED FROM: %v", r)
	}
}
