package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/cli/common"
	"github.com/lashajini/mind-palace/cli/errors"
	"github.com/lashajini/mind-palace/pkg/addons"
	"github.com/lashajini/mind-palace/pkg/config"
	"github.com/lashajini/mind-palace/pkg/models"
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

	copyFile(file, dst)

	cfg := config.NewConfig()
	rpcClient := rpcclient.NewClient(cfg)
	db := database.InitDB(cfg)
	memory := models.NewMemory()

	// TODO: ctx
	ctx := context.Background()
	tx := database.NewMultiInstruction(ctx, db.DB())
	defer revert(dst, tx)

	beginTransaction(tx)
	memoryID := insertMemory(tx, memory)
	resourcePath := filepath.Join(config.OriginalResourceRelativePath(currentUser), fileName)
	resource := models.NewResource(resourceID, memoryID, resourcePath)
	insertResource(tx, resource)
	commitTransaction(tx)

	userCfg := userConfig(currentUser)

	addonsResults, _ := rpcClient.Add(ctx, dst, memoryID, userCfg)
	for addonResult := range addonsResults {
		addons, err := addons.ToAddons(addonResult)
		errors.Handle(err)

		for _, addon := range addons {
			addon.Action(db, memoryID)
		}
	}
}

func revert(dst string, tx *database.MultiInstruction) {
	if r := recover(); r != nil {
		fmt.Println("Error:", r)

		fmt.Println("Reverting..")

		err := tx.Rollback()
		errors.Handle(err)
		fmt.Println("Database rollbacked")

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
	userCfg, err := mpuser.ReadUserConfig(currentUser)
	errors.Handle(err)
	return userCfg
}

func commitTransaction(tx *database.MultiInstruction) {
	err := tx.Commit()
	errors.Panic(err)
}

func insertResource(tx *database.MultiInstruction, resource *models.OriginalResource) {
	err := models.InsertResourceTx(tx, resource)
	errors.Panic(err)
}

func insertMemory(tx *database.MultiInstruction, memory *models.Memory) uuid.UUID {
	memoryID, err := models.InsertMemoryTx(tx, memory)
	errors.Panic(err)
	return memoryID
}

func beginTransaction(tx *database.MultiInstruction) *database.MultiInstruction {
	err := tx.Begin()
	errors.Panic(err)

	return tx
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
