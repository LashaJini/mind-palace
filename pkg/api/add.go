package api

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/addons"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/mperrors"
	"github.com/lashajini/mind-palace/pkg/mpuser"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
	addonrpc "github.com/lashajini/mind-palace/pkg/rpc/palace/addon"
	vdbrpc "github.com/lashajini/mind-palace/pkg/rpc/vdb"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

func Add(file string) error {
	cfg := common.NewConfig()
	currentUser, err := getCurrentUser()
	if err != nil {
		return mperrors.On(err).Wrap("failed to get current user")
	}

	if err := validateFile(file); err != nil {
		return mperrors.On(err).Wrap("failed to validate file")
	}

	resourceID := uuid.New()
	fileExtension := filepath.Ext(file)
	fileName := resourceID.String() + fileExtension
	originalResourceFullPath := common.OriginalResourceFullPath(currentUser)
	dst := filepath.Join(originalResourceFullPath, fileName)
	resourcePath := filepath.Join(common.OriginalResourceRelativePath(currentUser), fileName)

	err = common.CopyFile(file, dst)
	if err != nil {
		return mperrors.On(err).Wrap("failed to copy file")
	}

	userCfg, err := mpuser.ReadConfig(currentUser)
	if err != nil {
		return mperrors.On(err).Wrap("failed to read user config")
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

	addonResultC, err := addonClient.Add(ctx, dst, userCfg.Steps())
	if err != nil {
		return mperrors.On(err).Wrap("failed to add addon")
	}

	for addonResult := range addonResultC {
		addons, err := addons.ToAddons(addonResult)
		if err != nil {
			return mperrors.On(err).Wrap("failed to convert addons")
		}

		for _, addon := range addons {
			wg.Add(1)

			go func() {
				defer wg.Done()
				err := addon.Action(ctx, db, memoryIDC, vdbGrpcClient, maxBufSize, resourceID, resourcePath, cancel)
				mperrors.On(err).Warn()
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

func validateFile(file string) error {
	exists, err := common.FileExists(file)
	if err != nil {
		return mperrors.On(err).Wrap("failed to check if file exists")
	}

	if !exists {
		return mperrors.Onf("file '%s' does not exist", file)
	}

	isText, err := common.IsTextFile(file)
	if err != nil {
		return mperrors.On(err).Wrap("failed to check if file is a text file")
	}

	if !isText {
		return mperrors.Onf("file '%s' is not a text file\n", file)
	}

	return nil
}

func getCurrentUser() (string, error) {
	ctx := context.Background()
	currentUser, err := common.CurrentUser()
	if err != nil {
		return "", mperrors.On(err).Wrap("failed to get current user")
	}

	if currentUser == "" {
		msg := "there are no users available. Create one by using: mind-palace user --new <name>"
		return "", mperrors.Onf(msg)
	}

	loggers.Log.Info(ctx, "current user '%s'", currentUser)
	return currentUser, nil
}

func revertAdd(dst string, cancel context.CancelFunc) {
	if r := recover(); r != nil {
		ctx := context.Background()
		loggers.Log.Info(ctx, "Reverting...")

		err := os.Remove(dst)
		mperrors.On(err).Exit()

		loggers.Log.Info(ctx, "File removed %s", dst)

		cancel()
	}
}
