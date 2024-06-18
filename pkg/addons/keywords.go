package addons

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/errors"
	"github.com/lashajini/mind-palace/pkg/models"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/lashajini/mind-palace/pkg/types"
)

type KeywordsAddon struct {
	Addon
}

func (k *KeywordsAddon) Action(db *database.MindPalaceDB, memoryIDC chan uuid.UUID, args ...any) (err error) {
	keywords := k.Output.([]string)

	ctx := context.Background()
	tx := database.NewMultiInstruction(ctx, db)

	defer func() {
		if err != nil {
			err := tx.Rollback()
			errors.On(err).PanicWithMsg("failed to rollback")
		}
	}()

	err = tx.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	keywordIDs, err := models.InsertManyKeywordsTx(tx, keywords)
	if err != nil {
		return fmt.Errorf("failed to insert keywords: %w", err)
	}

	memoryID := <-memoryIDC
	err = models.InsertManyMemoryKeywordsTx(tx, keywordIDs, memoryID)
	if err != nil {
		return fmt.Errorf("failed to insert memory keyword pairs: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (k *KeywordsAddon) SetOutput(output any) {
	k.Output = output
}

var KeywordsAddonInstance = KeywordsAddon{
	Addon: Addon{
		Name:        types.AddonResourceKeywords,
		Description: "Extracts keywords from a resource",
		InputTypes:  []types.IOType{types.Text},
		OutputTypes: []types.IOType{types.Text},
	},
}
