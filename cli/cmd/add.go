package cli

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/addons"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/errors"
	"github.com/lashajini/mind-palace/pkg/mpuser"
	rpcclient "github.com/lashajini/mind-palace/pkg/rpc/client"
	"github.com/lashajini/mind-palace/pkg/storage/database"
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
	addCmd.MarkFlagRequired("file")
	addCmd.Flags().BoolVarP(&PREVIEW, "preview", "p", false, "preview result")
}

func Add(cmd *cobra.Command, args []string) {
	file, _ := cmd.Flags().GetString("file")

	add(file)
}

func add(file string) {
	cfg := common.NewConfig()
	currentUser := getCurrentUser()
	validateFile(file)

	resourceID := uuid.New()
	fileExtension := filepath.Ext(file)
	fileName := resourceID.String() + fileExtension
	originalResourceFullPath := common.OriginalResourceFullPath(currentUser)
	dst := filepath.Join(originalResourceFullPath, fileName)
	resourcePath := filepath.Join(common.OriginalResourceRelativePath(currentUser), fileName)

	copyFile(file, dst)

	userCfg := userConfig(currentUser)
	rpcClient := rpcclient.NewClient(cfg, userCfg)
	db := database.InitDB(cfg)
	db.SetSchema(db.ConstructSchema(currentUser))

	// TODO: ctx
	ctx := context.Background()
	defer revert(dst)

	maxBufSize := len(addons.List) - 1 // all addons - default
	memoryIDC := make(chan uuid.UUID, maxBufSize)

	var wg sync.WaitGroup

	addonResultC, _ := rpcClient.Add(ctx, dst)
	for addonResult := range addonResultC {
		addons, err := addons.ToAddons(addonResult)
		errors.On(err).Exit()

		for _, addon := range addons {
			wg.Add(1)

			go func() {
				defer wg.Done()
				err := addon.Action(db, memoryIDC, rpcClient, maxBufSize, resourceID, resourcePath)
				errors.On(err).Warn()
			}()
		}
	}

	wg.Wait()

	// clear channel
	for range len(memoryIDC) {
		<-memoryIDC
	}
}

func revert(dst string) {
	if r := recover(); r != nil {
		common.Log.Info().Msg("Reverting...")

		err := os.Remove(dst)
		errors.On(err).Exit()

		err = os.Remove(dst)
		errors.On(err).Exit()

		common.Log.Info().Msgf("File removed %s", dst)
	}
}

func copyFile(src, dst string) {
	err := common.CopyFile(src, dst)
	errors.On(err).Exit()
}

func userConfig(currentUser string) *mpuser.Config {
	userCfg, err := mpuser.ReadConfig(currentUser)
	errors.On(err).Exit()
	return userCfg
}

func validateFile(file string) {
	exists, err := common.FileExists(file)
	errors.On(err).Exit()

	if !exists {
		errors.ExitWithMsgf("File %s does not exist", file)
	}

	isText, err := common.IsTextFile(file)
	errors.On(err).Exit()

	if !isText {
		errors.ExitWithMsgf("File %s is not a text file\n", file)
	}
}

func getCurrentUser() string {
	currentUser, err := common.CurrentUser()
	errors.On(err).Exit()

	if currentUser == "" {
		msg := "there are no users available. Create one by using: mind-palace user --new <name>"
		errors.ExitWithMsg(msg)
	}

	common.Log.Info().Msgf("current user %s", currentUser)
	return currentUser
}
