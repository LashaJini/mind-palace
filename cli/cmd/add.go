package cli

import (
	"github.com/lashajini/mind-palace/pkg/api"
	"github.com/lashajini/mind-palace/pkg/mperrors"
	"github.com/spf13/cobra"
)

var (
	ADD_FILE    string
	ADD_PREVIEW bool
)

var addCmd = &cobra.Command{
	Use:   "add -f [FILE]",
	Short: "add <short description>",
	Long:  "add <long description>",
	Run:   Add,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&ADD_FILE, "file", "f", "", "file location")

	_ = addCmd.MarkFlagRequired("file")

	addCmd.Flags().BoolVarP(&ADD_PREVIEW, "preview", "p", false, "preview result")
}

func Add(cmd *cobra.Command, args []string) {
	file, _ := cmd.Flags().GetString("file")

	err := api.Add(cmd.Context(), file)
	mperrors.On(err).Exit()
}
