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
			return fmt.Errorf("Error: Could not create memory hierarchy.")
		}

		if err := os.MkdirAll(resourceHierarchy, os.ModePerm); err != nil {
			return fmt.Errorf("Error: Could not create resource hierarchy.")
		}

		userConfig := NewUserConfig(user)
		d, err := json.Marshal(userConfig)
		if err != nil {
			return fmt.Errorf("Error: Could not encode user config.")
		}

		if err := os.WriteFile(common.UserConfigPath(user, true), d, 0777); err != nil {
			fmt.Println(err)
			return fmt.Errorf("Error: Could not create user config.")
		}

		return nil
	}

	fmt.Println("User already exists. No actions taken.")
	return nil
}
