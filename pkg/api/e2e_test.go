package api

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/lashajini/mind-palace/pkg/addons"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/errors"
	"github.com/lashajini/mind-palace/pkg/mpuser"
	rpcclient "github.com/lashajini/mind-palace/pkg/rpc/client"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

type E2ETestSuite struct {
	suite.Suite
	client  *rpcclient.Client
	db      *database.MindPalaceDB
	user    string
	cfg     *common.Config
	userCfg *mpuser.Config
	schema  string
	ctx     context.Context

	input_text     []byte
	input_filepath string
}

// ensures test user exists
// ensures rpc server is up
// inits db
// init vdb is up
func (s *E2ETestSuite) SetupSuite() {
	s.ctx = context.Background()
	defer revert(s)

	var err error
	s.user = "test_user"
	s.userCfg, err = mpuser.CreateNewUser(s.user)
	errors.On(err).Panic()

	common.UpdateMindPalaceInfo(common.MindPalaceInfo{CurrentUser: s.user})

	s.cfg = common.NewConfig()

	s.client = rpcclient.NewClient(s.cfg, s.userCfg)
	err = s.client.Ping(s.ctx)
	errors.On(err).Panic()

	s.db = database.InitDB(s.cfg)
	s.schema, err = s.db.CreateSchema(s.user)
	errors.On(err).Panic()

	schemas, err := s.db.ListMPSchemas()
	errors.On(err).Panic()

	sqlTemplates := common.NewSQLTemplates(schemas)

	inMemorySource, ups := database.Up(s.cfg, sqlTemplates)
	database.CommitMigration(s.cfg, inMemorySource, ups)

	err = s.client.VDBPing(s.ctx)
	errors.On(err).Panic()

	f, err := os.CreateTemp("", "mindpalace_test_file")
	errors.On(err).Panic()

	s.input_text = []byte("The tower is 324 metres (1,063 ft) tall, about the same height as an 81-storey building, and the tallest structure in Paris. Its base is square, measuring 125 metres (410 ft) on each side. During its construction, the Eiffel Tower surpassed the Washington Monument to become the tallest man-made structure in the world, a title it held for 41 years until the Chrysler Building in New York City was finished in 1930. It was the first structure to reach a height of 300 metres. Due to the addition of a broadcasting aerial at the top of the tower in 1957, it is now taller than the Chrysler Building by 5.2 metres (17 ft). Excluding transmitters, the Eiffel Tower is the second tallest free-standing structure in France after the Millau Viaduct.")
	_, err = f.Write(s.input_text)
	errors.On(err).Panic()

	s.input_filepath = f.Name()
	common.Log.Info().Msgf("temp file created '%s'", s.input_filepath)
}

func (s *E2ETestSuite) Test_add_with_default_addon() {
	t := s.T()
	defer s.TearDownSubTest()

	Add(s.input_filepath)

	var total_memory int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.memory", s.schema)).Scan(&total_memory)

	assert.Equal(t, 1, total_memory)
}

func (s *E2ETestSuite) Test_add_with_default_and_keywords_addons_joined() {
	t := s.T()
	defer s.TearDownSubTest()

	err := s.userCfg.EnableAddon(&addons.KeywordsAddonInstance)
	assert.NoError(t, err)

	Add(s.input_filepath)

	var total_keywords int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.keyword", s.schema)).Scan(&total_keywords)

	assert.Greater(t, total_keywords, 0)

	var total_chunks int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.chunk", s.schema)).Scan(&total_chunks)

	assert.Greater(t, total_chunks, 0)
}

func (s *E2ETestSuite) Test_add_with_default_and_summary_addons_joined() {
	t := s.T()
	defer s.TearDownSubTest()

	err := s.userCfg.EnableAddon(&addons.SummaryAddonInstance)
	assert.NoError(t, err)

	Add(s.input_filepath)

	var total_summary int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.summary", s.schema)).Scan(&total_summary)

	assert.Equal(t, total_summary, 1)
}

func (s *E2ETestSuite) Test_add_with_all_addons_joined() {
	t := s.T()
	defer s.TearDownSubTest()

	err := s.userCfg.EnableAddon(&addons.KeywordsAddonInstance)
	assert.NoError(t, err)
	err = s.userCfg.EnableAddon(&addons.SummaryAddonInstance)
	assert.NoError(t, err)

	Add(s.input_filepath)

	var total_memory int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.memory", s.schema)).Scan(&total_memory)

	assert.Equal(t, total_memory, 1)

	var total_summary int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.summary", s.schema)).Scan(&total_summary)

	assert.Equal(t, total_summary, 1)

	var total_keywords int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.keyword", s.schema)).Scan(&total_keywords)

	assert.Greater(t, total_keywords, 0)

	var total_chunks int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.chunk", s.schema)).Scan(&total_chunks)

	assert.Greater(t, total_chunks, 0)
}

func (s *E2ETestSuite) Test_add_with_default_and_keywords_addons_single() {
	t := s.T()
	defer s.TearDownSubTest()

	serverCfg := make(map[string]string)
	serverCfg["available_tokens"] = "1"
	err := s.client.SetConfig(s.ctx, serverCfg)
	assert.NoError(t, err)

	err = s.userCfg.EnableAddon(&addons.KeywordsAddonInstance)
	assert.NoError(t, err)

	Add(s.input_filepath)

	var total_memory int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.memory", s.schema)).Scan(&total_memory)

	assert.Equal(t, 1, total_memory)

	var total_keywords int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.keyword", s.schema)).Scan(&total_keywords)

	assert.Greater(t, total_keywords, 0)

	var total_chunks int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.chunk", s.schema)).Scan(&total_chunks)

	assert.Greater(t, total_chunks, 0)
}

func (s *E2ETestSuite) Test_add_with_default_and_summary_addons_single() {
	t := s.T()
	defer s.TearDownSubTest()

	serverCfg := make(map[string]string)
	serverCfg["available_tokens"] = "1"
	err := s.client.SetConfig(s.ctx, serverCfg)
	assert.NoError(t, err)

	err = s.userCfg.EnableAddon(&addons.SummaryAddonInstance)
	assert.NoError(t, err)

	Add(s.input_filepath)

	var total_memory int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.memory", s.schema)).Scan(&total_memory)

	assert.Equal(t, 1, total_memory)

	var total_summary int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.summary", s.schema)).Scan(&total_summary)

	assert.Equal(t, total_summary, 1)
}

func (s *E2ETestSuite) Test_add_with_all_addons_single() {
	t := s.T()
	defer s.TearDownSubTest()

	serverCfg := make(map[string]string)
	serverCfg["available_tokens"] = "1"
	err := s.client.SetConfig(s.ctx, serverCfg)
	assert.NoError(t, err)

	err = s.userCfg.EnableAddon(&addons.KeywordsAddonInstance)
	assert.NoError(t, err)

	err = s.userCfg.EnableAddon(&addons.SummaryAddonInstance)
	assert.NoError(t, err)

	Add(s.input_filepath)

	var total_memory int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.memory", s.schema)).Scan(&total_memory)

	assert.Equal(t, 1, total_memory)

	var total_summary int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.summary", s.schema)).Scan(&total_summary)

	assert.Equal(t, total_summary, 1)

	var total_keywords int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.keyword", s.schema)).Scan(&total_keywords)

	assert.Greater(t, total_keywords, 0)

	var total_chunks int
	s.db.DB().QueryRow(fmt.Sprintf("SELECT count(*) FROM %s.chunk", s.schema)).Scan(&total_chunks)

	assert.Greater(t, total_chunks, 0)
}

func (s *E2ETestSuite) TearDownSubTest() {
	s.db.DB().Exec(fmt.Sprintf("DELETE FROM %s.memory", s.schema))
	s.db.DB().Exec(fmt.Sprintf("DELETE FROM %s.keyword", s.schema))

	s.client.SetConfig(s.ctx, nil)

	common.RemoveAllFiles(common.OriginalResourceFullPath(s.user))

	for _, addon := range addons.List {
		s.userCfg.DisableAddon(addon)
	}
}

// deletes test user
// drop vdb
func (s *E2ETestSuite) TearDownSuite() {
	var err error
	common.Log.Info().Msg("starting E2E test teardown")

	err = mpuser.DeleteUser(s.user)
	if err != nil {
		common.Log.Warn().Msgf("could not delete user '%s', reason: %v", s.user, err)
	} else {
		common.Log.Info().Msgf("user '%s' deleted", s.user)
	}
	common.UpdateMindPalaceInfo(common.MindPalaceInfo{CurrentUser: ""})

	err = s.client.VDBDrop(s.ctx)
	if err != nil {
		common.Log.Warn().Msgf("could not drop vector database, reason: %v", err)
	} else {
		common.Log.Info().Msgf("vector database dropped")
	}

	err = os.Remove(s.input_filepath)
	if err != nil {
		common.Log.Warn().Msgf("could not delete temp file '%s', reason: %v", s.input_filepath, err)
	} else {
		common.Log.Info().Msgf("temp file '%s' deleted", s.input_filepath)
	}
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}

func revert(s *E2ETestSuite) {
	if r := recover(); r != nil {
		common.Log.Info().Msgf("RECOVERED FROM: %v", r)

		s.TearDownSuite()
		panic(r)
	}
}
