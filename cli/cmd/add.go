package cli

import (
	"github.com/lashajini/mind-palace/pkg/api"
	"github.com/lashajini/mind-palace/pkg/mperrors"
	"github.com/spf13/cobra"
)

var (
	FILE    string
	PREVIEW bool
)

var addCmd = &cobra.Command{
	Use:   "add -f [FILE]",
	Short: "add <short description>",
	Long:  "add <long description>",
	Run:   Add,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&FILE, "file", "f", "", "file location")

	_ = addCmd.MarkFlagRequired("file")

	addCmd.Flags().BoolVarP(&PREVIEW, "preview", "p", false, "preview result")
}

func Add(cmd *cobra.Command, args []string) {
	file, _ := cmd.Flags().GetString("file")

	err := api.Add(file)
	mperrors.On(err).Exit()
}
