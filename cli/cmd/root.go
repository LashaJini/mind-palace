package cli

import (
	"github.com/lashajini/mind-palace/pkg/mperrors"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mind-palace",
	Short: "Mind Palace <short description>",
	Long:  "Mind Palace <long description>",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	err := rootCmd.Execute()
	mperrors.On(err).Exit()
}

func closeDB(db *database.MindPalaceDB) {
	if err := db.DB().Close(); err != nil {
		mperrors.On(err).Exit()
	}
}
