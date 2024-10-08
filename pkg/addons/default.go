package addons

import (
	"context"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/models"
	"github.com/lashajini/mind-palace/pkg/mperrors"
	vdbrpc "github.com/lashajini/mind-palace/pkg/rpc/vdb"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/lashajini/mind-palace/pkg/types"
)

type DefaultAddon struct {
	Addon
}

func (d *DefaultAddon) Action(ctx context.Context, db *database.MindPalaceDB, memoryIDC chan uuid.UUID, args ...any) (err error) {
	vdbGrpcClient := args[0].(*vdbrpc.Client)
	maxBufSize := args[1].(int)
	resourceID := args[2].(uuid.UUID)
	resourcePath := args[3].(string)
	cancel := args[4].(context.CancelFunc)

	tx := database.NewMultiInstruction(db)
	defer func() {
		if r := recover(); r != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = mperrors.On(rollbackErr).Wrap("default rollback failed")
			} else {
				err = mperrors.Onf("(recovered) panic: %v", r)
			}

			cancel()
		}
	}()

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = mperrors.On(rollbackErr).Wrap("default rollback failed")
			}

			cancel()
		}
	}()

	memory := models.NewMemory()
	err = tx.BeginTx(ctx)
	if err != nil {
		return mperrors.On(err).Wrap("default transaction begin failed")
	}

	memoryID, err := models.InsertMemoryTx(ctx, tx, memory)
	if err != nil {
		return mperrors.On(err).Wrap("default insert memory failed")
	}

	resource := models.NewResource(resourceID, memoryID, resourcePath)

	err = models.InsertResourceTx(ctx, tx, resource)
	if err != nil {
		return mperrors.On(err).Wrap("default insert resource failed")
	}

	defaultResponse := d.Response.GetDefaultResponse().Default
	if defaultResponse == "" {
		return mperrors.Onf("server didn't send default addon response")
	}

	err = vdbGrpcClient.Insert(ctx, []uuid.UUID{memoryID}, []string{defaultResponse}, []string{vdbrpc.ROW_TYPE_WHOLE})
	if err != nil {
		return mperrors.On(err).Wrap("default vdb insert failed")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return mperrors.On(err).Wrap("default transaction commit failed")
	}

	for range maxBufSize {
		memoryIDC <- memoryID
	}

	return nil
}

var DefaultAddonInstance = DefaultAddon{
	Addon: Addon{
		Name:        types.AddonDefault,
		Description: "Default",
		InputTypes:  []types.IOType{types.Text},
		OutputTypes: []types.IOType{types.Text},
	},
}
