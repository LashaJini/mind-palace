package cli

import (
	"fmt"

	"github.com/lashajini/mind-palace/pkg/addons"
	"github.com/lashajini/mind-palace/pkg/config"
	"github.com/lashajini/mind-palace/pkg/mpuser"
	"github.com/spf13/cobra"
)

var (
	LIST    bool
	ENABLE  string
	DISABLE string
)

var addonsCmd = &cobra.Command{
	Use:   "addons",
	Short: "addons <short description>",
	Long:  "addons <long description>",
	Run:   Addons,
}

func init() {
	rootCmd.AddCommand(addonsCmd)
	addonsCmd.Flags().BoolVarP(&LIST, "list", "l", false, "list available addons")
	addonsCmd.Flags().StringVarP(&ENABLE, "enable", "e", "", "enable addon")
	addonsCmd.Flags().StringVarP(&DISABLE, "disable", "d", "", "disable addon")
}

func Addons(cmd *cobra.Command, args []string) {
	if LIST {
		for _, addon := range addons.List {
			fmt.Println(addon)
			fmt.Println()
		}

		return
	}

	if ENABLE != "" {
		addon := addons.Find(ENABLE)
		user, _ := config.CurrentUser()
		userCfg, _ := mpuser.ReadUserConfig(user)
		_ = userCfg.EnableAddon(addon)

		return
	}

	if DISABLE != "" {
		addon := addons.Find(DISABLE)
		user, _ := config.CurrentUser()
		userCfg, _ := mpuser.ReadUserConfig(user)
		_ = userCfg.DisableAddon(addon)

		return
	}
}
