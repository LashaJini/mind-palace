package addons

import (
	"context"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/errors"
	"github.com/lashajini/mind-palace/pkg/models"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/lashajini/mind-palace/pkg/types"
)

type DefaultAddon struct {
	Addon
}

func (d *DefaultAddon) Action(db *database.MindPalaceDB, memoryIDC chan uuid.UUID, args ...any) (err error) {
	maxBufSize := args[0].(int)
	resourceID := args[1].(uuid.UUID)
	resourcePath := args[2].(string)

	ctx := context.Background()
	tx := database.NewMultiInstruction(ctx, db.DB())

	memory := models.NewMemory()
	err = tx.Begin()
	errors.Panic(err)

	memoryID, err := models.InsertMemoryTx(tx, memory)
	errors.Panic(err)

	resource := models.NewResource(resourceID, memoryID, resourcePath)

	err = models.InsertResourceTx(tx, resource)
	errors.Panic(err)

	err = tx.Commit()
	errors.Panic(err)

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
