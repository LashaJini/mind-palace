package cli

import (
	"github.com/spf13/cobra"
)

var (
	FILE    string
	PREVIEW bool
)

func init() {
	rootCmd.AddCommand(insertCmd)
	insertCmd.Flags().StringVarP(&FILE, "file", "f", "", "file location")
	insertCmd.MarkFlagRequired("file")
	insertCmd.Flags().BoolVarP(&PREVIEW, "preview", "p", false, "preview result")
}

var insertCmd = &cobra.Command{
	Use:   "add -f [FILE]",
	Short: "add <short description>",
	Long:  "add <long description>",
	Run:   func(cmd *cobra.Command, args []string) {},
}
