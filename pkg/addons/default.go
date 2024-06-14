package addons

import (
	"context"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/errors"
	"github.com/lashajini/mind-palace/pkg/models"
	rpcclient "github.com/lashajini/mind-palace/pkg/rpc/client"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/lashajini/mind-palace/pkg/types"
)

type DefaultAddon struct {
	Addon
}

func (d *DefaultAddon) Action(db *database.MindPalaceDB, memoryIDC chan uuid.UUID, args ...any) (err error) {
	rpcClient := args[0].(*rpcclient.Client)
	maxBufSize := args[1].(int)
	resourceID := args[2].(uuid.UUID)
	resourcePath := args[3].(string)

	ctx := context.Background()
	tx := database.NewMultiInstruction(ctx, db.DB())
	defer revert(tx)

	memory := models.NewMemory()
	err = tx.Begin()
	errors.On(err).Panic()

	memoryID, err := models.InsertMemoryTx(tx, memory)
	errors.On(err).Panic()

	resource := models.NewResource(resourceID, memoryID, resourcePath)

	err = models.InsertResourceTx(tx, resource)
	errors.On(err).Panic()

	err = tx.Commit()
	errors.On(err).Panic()

	for range maxBufSize {
		memoryIDC <- memoryID
	}

	err = rpcClient.VDBInsert(ctx, memoryID, d.Output.([]string)[0])
	errors.On(err).Panic()

	return nil
}

func revert(tx *database.MultiInstruction) {
	if r := recover(); r != nil {
		err := tx.Rollback()
		errors.On(err).PanicWithMsg("failed to rollback")

		panic(r)
	}
}

var DefaultAddonInstance = DefaultAddon{
	Addon: Addon{
		Name:        types.AddonDefault,
		Description: "Default",
		InputTypes:  []types.IOType{types.Text},
		OutputTypes: []types.IOType{types.Text},
	},
}
