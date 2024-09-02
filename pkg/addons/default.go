package addons

import (
	"context"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/errors"
	"github.com/lashajini/mind-palace/pkg/models"
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
			rollback(tx)

			cancel()
		}
	}()

	memory := models.NewMemory()
	err = tx.Begin()
	errors.On(err).Panic()

	memoryID, err := models.InsertMemoryTx(tx, memory)
	errors.On(err).Panic()

	resource := models.NewResource(resourceID, memoryID, resourcePath)

	err = models.InsertResourceTx(tx, resource)
	errors.On(err).Panic()

	defaultResponse := d.Response.GetDefaultResponse().Default
	if defaultResponse == "" {
		errors.ExitWithMsg("reason: server didn't send default addon response")
	}

	err = vdbGrpcClient.Insert(ctx, []uuid.UUID{memoryID}, []string{defaultResponse})
	errors.On(err).Panic()

	err = tx.Commit()
	errors.On(err).Panic()

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
