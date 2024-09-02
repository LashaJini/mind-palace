package api

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
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
	addonrpc "github.com/lashajini/mind-palace/pkg/rpc/palace/addon"
	vdbrpc "github.com/lashajini/mind-palace/pkg/rpc/vdb"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

func Add(file string) error {
	cfg := common.NewConfig()
	currentUser := getCurrentUser()
	validateFile(file)

	resourceID := uuid.New()
	fileExtension := filepath.Ext(file)
	fileName := resourceID.String() + fileExtension
	originalResourceFullPath := common.OriginalResourceFullPath(currentUser)
	dst := filepath.Join(originalResourceFullPath, fileName)
	resourcePath := filepath.Join(common.OriginalResourceRelativePath(currentUser), fileName)

	err := common.CopyFile(file, dst)
	if err != nil {
		return err
	}

	userCfg, err := mpuser.ReadConfig(currentUser)
	if err != nil {
		return err
	}

	addonClient := addonrpc.NewGrpcClient(cfg)
	vdbGrpcClient := vdbrpc.NewGrpcClient(cfg, userCfg.Config.User)
	db := database.InitDB(cfg)
	db.SetSchema(db.ConstructSchema(currentUser))

	ctx, cancel := context.WithCancel(context.Background())
	defer revertAdd(dst, cancel)

	maxBufSize := len(addons.List) - 1 // all addons - default
	memoryIDC := make(chan uuid.UUID, maxBufSize)

	var wg sync.WaitGroup

	addonResultC, _ := addonClient.Add(ctx, dst, userCfg.Steps())
	for addonResult := range addonResultC {
		addons, err := addons.ToAddons(addonResult)
		if err != nil {
			return err
		}

		for _, addon := range addons {
			wg.Add(1)

			go func() {
				defer wg.Done()
				err := addon.Action(ctx, db, memoryIDC, vdbGrpcClient, maxBufSize, resourceID, resourcePath, cancel)
				errors.On(err).Warn()
			}()
		}
	}

	wg.Wait()

	// clear channel
	for range len(memoryIDC) {
		<-memoryIDC
	}

	return nil
}

func validateFile(file string) {
	exists, err := common.FileExists(file)
	errors.On(err).Exit()

	if !exists {
		errors.ExitWithMsgf("file '%s' does not exist", file)
	}

	isText, err := common.IsTextFile(file)
	errors.On(err).Exit()

	if !isText {
		errors.ExitWithMsgf("file '%s' is not a text file\n", file)
	}
}

func getCurrentUser() string {
	ctx := context.Background()
	currentUser, err := common.CurrentUser()
	errors.On(err).Exit()

	if currentUser == "" {
		msg := "there are no users available. Create one by using: mind-palace user --new <name>"
		errors.ExitWithMsg(msg)
	}

	loggers.Log.Info(ctx, "current user '%s'", currentUser)
	return currentUser
}

func revertAdd(dst string, cancel context.CancelFunc) {
	if r := recover(); r != nil {
		ctx := context.Background()
		loggers.Log.Info(ctx, "Reverting...")

		err := os.Remove(dst)
		errors.On(err).Exit()

		err = os.Remove(dst)
		errors.On(err).Exit()

		loggers.Log.Info(ctx, "File removed %s", dst)

		cancel()
	}
}
