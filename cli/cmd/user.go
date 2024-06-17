package cli

import (
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/errors"
	"github.com/lashajini/mind-palace/pkg/mpuser"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/spf13/cobra"
)

var (
	NEW    string
	SWITCH string
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

	userCmd.MarkFlagsOneRequired("new", "switch")
	userCmd.MarkFlagsMutuallyExclusive("new", "switch")
}

func User(cmd *cobra.Command, args []string) {
	newUser, _ := cmd.Flags().GetString("new")
	switchUser, _ := cmd.Flags().GetString("switch")

	currentUser := ""
	if cmd.Flags().Changed("new") {
		if len(newUser) < 1 {
			errors.ExitWithMsg("new user cannot be empty")
		}

		currentUser = newUser

		err := mpuser.CreateMindPalace(newUser)
		errors.On(err).Exit()

		cfg := common.NewConfig()
		db := database.InitDB(cfg)
		defer db.DB().Close()

		err = db.CreateSchema(newUser)
		errors.On(err).Exit()
	} else if cmd.Flags().Changed("switch") {
		mindPalaceUserPath := common.UserPath(switchUser, true)

		exists, err := common.DirExists(mindPalaceUserPath)
		errors.On(err).Exit()

		if !exists {
			errors.ExitWithMsgf("user '%s' does not exist\n", switchUser)
		}

		currentUser = switchUser
	}

	common.UpdateMindPalaceInfo(common.MindPalaceInfo{CurrentUser: currentUser})
	user(args...)
}

func user(args ...string) {}
