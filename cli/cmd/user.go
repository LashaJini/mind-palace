package cli

import (
	"fmt"
	"os"

	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/errors"
	"github.com/lashajini/mind-palace/pkg/mpuser"
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
		fmt.Println("either --new or --switch must be provided")
		cmd.Usage()
		os.Exit(1)
	}

	if newUser != "" && switchUser != "" {
		fmt.Println("only one of --new or --switch can be provided")
		cmd.Usage()
		os.Exit(1)
	}

	if newUser != "" {
		CURRENT_USER = newUser

		if err := mpuser.CreateMindPalace(newUser); err != nil {
			errors.On(err).Exit()
		}
	} else if switchUser != "" {
		mindPalaceUserPath := common.UserPath(switchUser, true)

		exists, err := common.DirExists(mindPalaceUserPath)
		errors.On(err).Exit()

		if !exists {
			fmt.Printf("user '%s' does not exist\n", switchUser)
			cmd.Usage()
			os.Exit(1)
		}

		CURRENT_USER = switchUser
	}

	common.UpdateMindPalaceInfo(common.MindPalaceInfo{CurrentUser: CURRENT_USER})
	user(args...)
}

func user(args ...string) {}
