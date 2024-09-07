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
	ADDONS_LIST               bool
	ADDONS_ENABLE_ADDON_NAME  string
	ADDONS_DISABLE_ADDON_NAME string
)

var addonsCmd = &cobra.Command{
	Use:   "addons",
	Short: "addons <short description>",
	Long:  "addons <long description>",
	Run:   Addons,
}

func init() {
	rootCmd.AddCommand(addonsCmd)
	addonsCmd.Flags().BoolVarP(&ADDONS_LIST, "list", "l", false, "list available addons")
	addonsCmd.Flags().StringVarP(&ADDONS_ENABLE_ADDON_NAME, "enable", "e", "", "enable addon <name>")
	addonsCmd.Flags().StringVarP(&ADDONS_DISABLE_ADDON_NAME, "disable", "d", "", "disable addon <name>")
}

func Addons(cmd *cobra.Command, args []string) {
	if ADDONS_LIST {
		listAddons()
		return
	}

	if ADDONS_ENABLE_ADDON_NAME != "" {
		handleAddon(ADDONS_ENABLE_ADDON_NAME, true)
		return
	}

	if ADDONS_DISABLE_ADDON_NAME != "" {
		handleAddon(ADDONS_DISABLE_ADDON_NAME, false)
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
