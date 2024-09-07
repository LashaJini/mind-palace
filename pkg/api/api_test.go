package api

import (
	"log"
	"os"
	"testing"

	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/mpuser"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	mpuser string
}

func (suite *APITestSuite) SetupSuite() {
	suite.mpuser = common.TEST_USER
	_, err := mpuser.CreateMindPalace(suite.mpuser)
	if err != nil {
		log.Fatal(err)
	}

	err = common.UpdateMindPalaceInfo(common.MindPalaceInfo{CurrentUser: suite.mpuser})
	if err != nil {
		log.Fatal(err)
	}
}

func (suite *APITestSuite) TearDownSuite() {
	err := os.RemoveAll(common.UserPath(suite.mpuser, true))
	if err != nil {
		log.Fatalf("failed to remove user mind palace: %s", err)
	}

	err = common.UpdateMindPalaceInfo(common.MindPalaceInfo{CurrentUser: ""})
	if err != nil {
		log.Fatalf("failed to update config: %s", err)
	}
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
