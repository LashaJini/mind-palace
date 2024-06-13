package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/addons"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/config"
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
	currentUser := getCurrentUser()
	validateFile(file)

	resourceID := uuid.New()
	fileExtension := filepath.Ext(file)
	fileName := resourceID.String() + fileExtension
	originalResourceFullPath := config.OriginalResourceFullPath(currentUser)
	dst := filepath.Join(originalResourceFullPath, fileName)
	resourcePath := filepath.Join(config.OriginalResourceRelativePath(currentUser), fileName)

	copyFile(file, dst)

	cfg := config.NewConfig()
	rpcClient := rpcclient.NewClient(cfg)
	db := database.InitDB(cfg)

	// TODO: ctx
	ctx := context.Background()
	defer revert(dst)

	maxBufSize := len(addons.List) - 1 // all addons - default
	memoryIDC := make(chan uuid.UUID, maxBufSize)

	var wg sync.WaitGroup

	userCfg := userConfig(currentUser)
	addonResultC, _ := rpcClient.Add(ctx, dst, userCfg)
	for addonResult := range addonResultC {
		addons, err := addons.ToAddons(addonResult)
		errors.Handle(err)

		for _, addon := range addons {
			wg.Add(1)

			go func() {
				defer wg.Done()
				addon.Action(db, memoryIDC, maxBufSize, resourceID, resourcePath)
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
		fmt.Println("Error:", r)

		fmt.Println("Reverting..")

		err := os.Remove(dst)
		errors.Handle(err)

		err = os.Remove(dst)
		errors.Handle(err)
		fmt.Println("File removed", dst)
	}
}

func copyFile(src, dst string) {
	err := common.CopyFile(src, dst)
	errors.Handle(err)
}

func userConfig(currentUser string) *mpuser.Config {
	userCfg, err := mpuser.ReadConfig(currentUser)
	errors.Handle(err)
	return userCfg
}

func validateFile(file string) {
	exists, err := common.FileExists(file)
	errors.Handle(err)

	if !exists {
		fmt.Printf("Error: File %s does not exist\n", file)
		os.Exit(1)
	}

	isText, err := common.IsTextFile(file)
	errors.Handle(err)

	if !isText {
		fmt.Printf("Error: File %s is not a text file\n", file)
		os.Exit(1)
	}
}

func getCurrentUser() string {
	currentUser, err := config.CurrentUser()
	errors.Handle(err)

	if currentUser == "" {
		fmt.Println("Error: There are no users available.")
		fmt.Printf("\nCreate one by using: mind-palace user --new <name>\n\n")
		os.Exit(1)
	}

	return currentUser
}
