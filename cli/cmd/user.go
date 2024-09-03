package cli

import (
	"bytes"
	"path/filepath"

	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/mperrors"
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
			mperrors.ExitWithMsg("new user cannot be empty")
		}

		_, err := mpuser.CreateMindPalace(newUser)
		mperrors.On(err).Exit()

		cfg := common.NewConfig()
		db := database.InitDB(cfg)
		defer closeDB(db)

		_, err = db.CreateSchema(newUser)
		mperrors.On(err).Exit()

		err = migrateNewUser(cfg, db)
		mperrors.On(err).Exit()

		currentUser = newUser
	} else if cmd.Flags().Changed("switch") {
		mindPalaceUserPath := common.UserPath(switchUser, true)

		exists, err := common.DirExists(mindPalaceUserPath)
		mperrors.On(err).Exit()

		if !exists {
			mperrors.ExitWithMsgf("user '%s' does not exist\n", switchUser)
		}

		currentUser = switchUser
	}

	err := common.UpdateMindPalaceInfo(common.MindPalaceInfo{CurrentUser: currentUser})
	mperrors.On(err).Exit()
}

func migrateNewUser(cfg *common.Config, db *database.MindPalaceDB) error {
	pattern := filepath.Join(cfg.MIGRATIONS_DIR, "*.up.sql")
	sqlUpFiles, err := filepath.Glob(pattern)
	if err != nil {
		return mperrors.On(err).Wrap("failed to find '*.up.sql' files")
	}

	sqlTemplate := common.SQLTemplate{Namespace: db.CurrentSchema}
	for _, sqlUpFile := range sqlUpFiles {
		var sqlBuffer bytes.Buffer
		err = sqlTemplate.Inject(&sqlBuffer, sqlUpFile)
		if err != nil {
			return mperrors.On(err).Wrap("failed to inject sql")
		}

		_, err = db.DB().Exec(sqlBuffer.String())
		if err != nil {
			return mperrors.On(err).Wrap("failed to exec sql")
		}
	}

	return nil
}
