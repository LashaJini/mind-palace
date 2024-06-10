package addons

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/models"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

type KeywordsAddon struct {
	Addon
}

func (k *KeywordsAddon) Action(db *database.MindPalaceDB, memoryID uuid.UUID, args ...any) (bool, error) {
	keywords := k.Output.([]string)

	ctx := context.Background()
	tx := database.NewMultiInstruction(ctx, db.DB())
	_ = tx.Begin()
	keywordIDs, _ := models.InsertManyKeywordsTx(tx, keywords)
	memoryKeywordID, _ := models.InsertManyMemoryKeywordsTx(tx, keywordIDs, memoryID)
	fmt.Println(memoryKeywordID)
	_ = tx.Commit()

	// TODO: rollback
	return true, nil
}

func (k *KeywordsAddon) SetOutput(output any) {
	k.Output = output
}

var KeywordsAddonInstance = KeywordsAddon{
	Addon: Addon{
		Name:        ResourceKeywords,
		Description: "Extracts keywords from a resource",
		InputTypes:  []Type{Text},
		OutputTypes: []Type{Text},
	},
}
