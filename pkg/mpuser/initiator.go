package mpuser

import (
	"context"
	"encoding/json"
	"os"

	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/mperrors"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
)

func CreateMindPalace(user string) (*Config, error) {
	ctx := context.Background()
	_, exists, err := UserExists(user)
	if err != nil {
		return nil, mperrors.On(err).Wrap("failed to check if user exists")
	}

	if !exists {
		memoryHierarchy := common.MemoryPath(user, true)
		resourceHierarchy := common.OriginalResourcePath(user, true)

		if err := os.MkdirAll(memoryHierarchy, os.ModePerm); err != nil {
			return nil, mperrors.On(err).Wrap("could not create memory hierarchy")
		}

		if err := os.MkdirAll(resourceHierarchy, os.ModePerm); err != nil {
			return nil, mperrors.On(err).Wrap("could not create resource hierarchy")
		}

		userConfig := NewUserConfig(user)
		d, err := json.Marshal(userConfig)
		if err != nil {
			return nil, mperrors.On(err).Wrap("could not encode user config")
		}

		if err := os.WriteFile(common.UserConfigPath(user, true), d, 0777); err != nil {
			return nil, mperrors.On(err).Wrap("could not create user config")
		}

		return userConfig, nil
	}

	loggers.Log.Warn(ctx, "user %s already exists. No actions taken", user)
	return ReadConfig(user)
}
