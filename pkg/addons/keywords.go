package addons

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/errors"
	"github.com/lashajini/mind-palace/pkg/models"
	rpcclient "github.com/lashajini/mind-palace/pkg/rpc/client"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/lashajini/mind-palace/pkg/types"
)

type KeywordsAddon struct {
	Addon
}

func (k *KeywordsAddon) Action(db *database.MindPalaceDB, memoryIDC chan uuid.UUID, args ...any) (err error) {
	vdbGrpcClient := args[0].(*rpcclient.VDBGrpcClient)
	keywordsChunks := k.Response.GetKeywordsResponse().List

	ctx := context.Background()
	tx := database.NewMultiInstruction(ctx, db)
	defer revert(tx)

	err = tx.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// chunk keywords may overlap with each other,
	// so we are going to deduplicate them
	var chunks []string
	uniqueKeywords := make(map[string]any)
	for _, keywordChunk := range keywordsChunks {
		for _, keyword := range keywordChunk.Keywords {
			uniqueKeywords[keyword] = struct{}{}
		}

		chunks = append(chunks, keywordChunk.Chunk)
	}

	var keywords []string
	for keyword := range uniqueKeywords {
		keywords = append(keywords, keyword)
	}

	keywordIDsMap, err := models.InsertManyKeywordsTx(tx, keywords)
	if err != nil {
		return fmt.Errorf("failed to insert keywords: %w", err)
	}

	memoryID := <-memoryIDC
	chunkIDs, err := models.InsertManyChunksTx(tx, memoryID, chunks)
	if err != nil {
		return fmt.Errorf("failed to insert chunks: %w", err)
	}

	chunkIDKeywordIDsMap := make(map[uuid.UUID][]int)
	for i, chunkID := range chunkIDs {
		var keywordIDs []int
		for _, keyword := range keywordsChunks[i].Keywords {
			keywordIDs = append(keywordIDs, keywordIDsMap[keyword])
		}

		chunkIDKeywordIDsMap[chunkID] = keywordIDs
	}

	err = models.InsertManyChunksKeywordsTx(tx, chunkIDKeywordIDsMap)
	if err != nil {
		return fmt.Errorf("failed to insert chunk keyword pairs: %w", err)
	}

	err = vdbGrpcClient.Insert(ctx, chunkIDs, chunks)
	errors.On(err).Panic()

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

var KeywordsAddonInstance = KeywordsAddon{
	Addon: Addon{
		Name:        types.AddonResourceKeywords,
		Description: "Extracts keywords from a resource",
		InputTypes:  []types.IOType{types.Text},
		OutputTypes: []types.IOType{types.Text},
	},
}
