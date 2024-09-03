package cli

import (
	"fmt"

	"github.com/lashajini/mind-palace/pkg/addons"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/mperrors"
	"github.com/lashajini/mind-palace/pkg/mpuser"
	"github.com/spf13/cobra"
)

var (
	LIST               bool
	ENABLE_ADDON_NAME  string
	DISABLE_ADDON_NAME string
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
	addonsCmd.Flags().StringVarP(&ENABLE_ADDON_NAME, "enable", "e", "", "enable addon <name>")
	addonsCmd.Flags().StringVarP(&DISABLE_ADDON_NAME, "disable", "d", "", "disable addon <name>")
}

func Addons(cmd *cobra.Command, args []string) {
	if LIST {
		listAddons()
		return
	}

	if ENABLE_ADDON_NAME != "" {
		handleAddon(ENABLE_ADDON_NAME, true)
		return
	}

	if DISABLE_ADDON_NAME != "" {
		handleAddon(DISABLE_ADDON_NAME, false)
		return
	}
}

func listAddons() {
	for _, addon := range addons.List {
		fmt.Println(addon)
	}
}

func handleAddon(addonName string, enable bool) {
	addon, err := addons.Find(addonName)
	mperrors.On(err).Exit()

	user, err := common.CurrentUser()
	mperrors.On(err).Exit()

	userCfg, err := mpuser.ReadConfig(user)
	mperrors.On(err).Exit()

	if enable {
		err = userCfg.EnableAddon(addon)
	} else {
		err = userCfg.DisableAddon(addon)
	}

	mperrors.On(err).Exit()
}
