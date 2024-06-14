package cli

import (
	"github.com/lashajini/mind-palace/pkg/errors"
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
	errors.On(err).Exit()
}
