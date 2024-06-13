package cli

import (
	"fmt"
	"os"

	"github.com/lashajini/mind-palace/pkg/common"
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

		if err := common.CreateMindPalace(newUser); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if switchUser != "" {
		mindPalaceUserPath := config.UserPath(switchUser, true)

		exists, err := common.DirExists(mindPalaceUserPath)
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
