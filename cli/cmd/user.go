package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lashajini/mind-palace/pkg/config"
	"github.com/spf13/cobra"
)

var (
	NEW          string
	SWITCH       string
	CURRENT_USER string
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "user <short description>",
	Long:  "user <long description>",
	Run:   User,
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.Flags().StringVarP(&NEW, "new", "n", "", "new user")
	userCmd.Flags().StringVarP(&SWITCH, "switch", "s", "", "switch user")
}

func User(cmd *cobra.Command, args []string) {
	newUser, _ := cmd.Flags().GetString("new")
	switchUser, _ := cmd.Flags().GetString("switch")

	if newUser == "" && switchUser == "" {
		fmt.Println("Error: Either --new or --switch must be provided")
		cmd.Usage()
		os.Exit(1)
	}

	if newUser != "" && switchUser != "" {
		fmt.Println("Error: Only one of --new or --switch can be provided")
		cmd.Usage()
		os.Exit(1)
	}

	if newUser != "" {
		CURRENT_USER = newUser

		if err := createMindPalace(newUser); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if switchUser != "" {
		mindPalaceUserPath := config.UserPath(switchUser, true)

		exists, err := dirExists(mindPalaceUserPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if !exists {
			fmt.Printf("Error: User '%s' does not exist.\n", switchUser)
			cmd.Usage()
			os.Exit(1)
		}

		CURRENT_USER = switchUser
	}

	config.UpdateMindPalaceInfo(config.MindPalaceInfo{CurrentUser: CURRENT_USER})
	user(args...)
}

func user(args ...string) {}

func createMindPalace(user string) error {
	mindPalaceUserPath := config.UserPath(user, true)

	exists, err := dirExists(mindPalaceUserPath)
	if err != nil {
		return err
	}

	if !exists {
		memoryHierarchy := config.MemoryPath(user, true)
		resourceHierarchy := config.OriginalResourcePath(user, true)

		if err := os.MkdirAll(memoryHierarchy, os.ModePerm); err != nil {
			return fmt.Errorf("Error: Could not create memory hierarchy.")
		}

		if err := os.MkdirAll(resourceHierarchy, os.ModePerm); err != nil {
			return fmt.Errorf("Error: Could not create resource hierarchy.")
		}

		userConfig := config.NewUserConfig(user)
		d, err := json.Marshal(userConfig)
		if err != nil {
			return fmt.Errorf("Error: Could not encode user config.")
		}

		if err := os.WriteFile(config.UserConfigPath(user, true), d, 0777); err != nil {
			fmt.Println(err)
			return fmt.Errorf("Error: Could not create user config.")
		}

		return nil
	}

	fmt.Println("User already exists. No actions taken.")
	return nil
}

func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}
