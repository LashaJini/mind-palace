package addons

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/models"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/lashajini/mind-palace/pkg/types"
)

type KeywordsAddon struct {
	Addon
}

func (k *KeywordsAddon) Action(db *database.MindPalaceDB, memoryID uuid.UUID, args ...any) (err error) {
	keywords := k.Output.([]string)

	ctx := context.Background()
	tx := database.NewMultiInstruction(ctx, db.DB())

	defer func() {
		if err != nil {
			fmt.Println(err)

			if err := tx.Rollback(); err != nil {
				panic(fmt.Errorf("Failed to rollback: %w", err))
			}
		}
	}()

	err = tx.Begin()
	if err != nil {
		return fmt.Errorf("Failed to start transaction: %w", err)
	}

	keywordIDs, err := models.InsertManyKeywordsTx(tx, keywords)
	if err != nil {
		return fmt.Errorf("Failed to insert keywords: %w", err)
	}

	err = models.InsertManyMemoryKeywordsTx(tx, keywordIDs, memoryID)
	if err != nil {
		return fmt.Errorf("Failed to insert memory keyword pairs: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
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
