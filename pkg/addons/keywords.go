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

type KeywordsAddon struct {
	Addon
}

func (k *KeywordsAddon) Action(ctx context.Context, db *database.MindPalaceDB, memoryIDC chan uuid.UUID, args ...any) (err error) {
	vdbGrpcClient := args[0].(*vdbrpc.Client)
	keywordsChunks := k.Response.GetKeywordsResponse().List

	tx := database.NewMultiInstruction(db)
	defer func() {
		if r := recover(); r != nil {
			if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
				err = mperrors.On(rollbackErr).Wrap("keywords rollback failed")
			} else {
				err = mperrors.Onf("(recovered) panic: %v", r)
			}
		}
	}()

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
				err = mperrors.On(rollbackErr).Wrap("keywords rollback failed")
			}
		}
	}()

	err = tx.Begin()
	if err != nil {
		return mperrors.On(err).Wrap("keywords transaction begin failed")
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

	keywordIDsMap, err := models.InsertManyKeywordsTx(ctx, tx, keywords)
	if err != nil {
		return mperrors.On(err).Wrap("failed to insert keywords")
	}

	select {
	case memoryID := <-memoryIDC:
		chunkIDs, err := models.InsertManyChunksTx(ctx, tx, memoryID, chunks)
		if err != nil {
			return mperrors.On(err).Wrap("failed to insert chunks")
		}

		chunkIDKeywordIDsMap := make(map[uuid.UUID][]int)
		for i, chunkID := range chunkIDs {
			var keywordIDs []int
			for _, keyword := range keywordsChunks[i].Keywords {
				keywordIDs = append(keywordIDs, keywordIDsMap[keyword])
			}

			chunkIDKeywordIDsMap[chunkID] = keywordIDs
		}

		err = models.InsertManyChunksKeywordsTx(ctx, tx, chunkIDKeywordIDsMap)
		if err != nil {
			return mperrors.On(err).Wrap("failed to insert chunk keyword pairs")
		}

		var types []string
		for range chunks {
			types = append(types, vdbrpc.ROW_TYPE_CHUNK)
		}

		err = vdbGrpcClient.Insert(ctx, chunkIDs, chunks, types)
		if err != nil {
			return mperrors.On(err).Wrap("keywords vdb insert failed")
		}

		err = tx.Commit(ctx)
		if err != nil {
			return mperrors.On(err).Wrap("keywords transaction commit failed")
		}

		return nil
	case <-ctx.Done():
		if ctx.Err() != nil {
			return mperrors.On(ctx.Err()).Wrap("keywords addon action failed")
		}

		return nil
	}
}

var KeywordsAddonInstance = KeywordsAddon{
	Addon: Addon{
		Name:        types.AddonResourceKeywords,
		Description: "Extracts keywords from a resource",
		InputTypes:  []types.IOType{types.Text},
		OutputTypes: []types.IOType{types.Text},
	},
}
