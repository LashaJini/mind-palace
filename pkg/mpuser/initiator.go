package mpuser

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lashajini/mind-palace/pkg/common"
)

func CreateMindPalace(user string) error {
	mindPalaceUserPath := common.UserPath(user, true)

	exists, err := common.DirExists(mindPalaceUserPath)
	if err != nil {
		return err
	}

	if !exists {
		memoryHierarchy := common.MemoryPath(user, true)
		resourceHierarchy := common.OriginalResourcePath(user, true)

		if err := os.MkdirAll(memoryHierarchy, os.ModePerm); err != nil {
			return fmt.Errorf("could not create memory hierarchy %w", err)
		}

		if err := os.MkdirAll(resourceHierarchy, os.ModePerm); err != nil {
			return fmt.Errorf("could not create resource hierarchy %w", err)
		}

		userConfig := NewUserConfig(user)
		d, err := json.Marshal(userConfig)
		if err != nil {
			return fmt.Errorf("could not encode user config %w", err)
		}

		if err := os.WriteFile(common.UserConfigPath(user, true), d, 0777); err != nil {
			return fmt.Errorf("could not create user config %w", err)
		}

		return nil
	}

	common.Log.Info().Msgf("user %s already exists. No actions taken", user)
	return nil
}
