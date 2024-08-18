package mpuser

import (
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/common"
	pb "github.com/lashajini/mind-palace/pkg/rpc/client/gen/proto"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/lashajini/mind-palace/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MpUserTestSuite struct {
	suite.Suite
	cfg *Config
}

func (suite *MpUserTestSuite) SetupSuite() {
	suite.cfg = NewUserConfig(common.TEST_USER)
}

type DefaultAddon struct {
	Name        string
	Description string
	InputTypes  types.IOTypes
	OutputTypes types.IOTypes
	Response    *pb.AddonResponse
}

func (a *DefaultAddon) Action(db *database.MindPalaceDB, memoryIDC chan uuid.UUID, args ...any) (err error) {
	return nil
}
func (a *DefaultAddon) Empty() bool                            { return true }
func (a *DefaultAddon) GetName() string                        { return a.Name }
func (a *DefaultAddon) GetDescription() string                 { return a.Description }
func (a *DefaultAddon) GetInputTypes() types.IOTypes           { return a.InputTypes }
func (a *DefaultAddon) GetOutputTypes() types.IOTypes          { return a.OutputTypes }
func (a *DefaultAddon) GetResponse() *pb.AddonResponse         { return a.Response }
func (a *DefaultAddon) SetResponse(response *pb.AddonResponse) { a.Response = response }
func (a DefaultAddon) String() string                          { return "" }

func (suite *MpUserTestSuite) Test_user_cant_disable_default_addon() {
	t := suite.T()

	_, err := CreateMindPalace(common.TEST_USER)
	if err != nil {
		log.Fatal(err)
	}
	err = common.UpdateMindPalaceInfo(common.MindPalaceInfo{CurrentUser: common.TEST_USER})
	if err != nil {
		log.Fatal(err)
	}

	defaultAddon := &DefaultAddon{Name: types.AddonDefault}
	err = suite.cfg.DisableAddon(defaultAddon)
	assert.NoError(t, err)
	assert.Contains(t, suite.cfg.Steps(), defaultAddon.GetName())
	assert.Equal(t, 1, len(suite.cfg.Steps()))
}

func (suite *MpUserTestSuite) TearDownSuite() {
	err := os.RemoveAll(common.UserPath(common.TEST_USER, true))
	if err != nil {
		log.Fatalf("failed to remove user mind palace: %s", err)
	}

	err = common.UpdateMindPalaceInfo(common.MindPalaceInfo{CurrentUser: ""})
	if err != nil {
		log.Fatalf("failed to update config: %s", err)
	}
}

func TestMpUserTestSuite(t *testing.T) {
	suite.Run(t, new(MpUserTestSuite))
}
