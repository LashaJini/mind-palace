package models

import (
	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

type ChunkKeyword struct {
	KeywordID int
	ChunkID   uuid.UUID
}

var chunkKeywordColumns = []string{
	"keyword_id",
	"chunk_id",
}

func InsertManyChunksKeywordsTx(tx *database.MultiInstruction, chunkIDKeywordIDsMap map[uuid.UUID][]int) error {
	joinedColumns, _ := joinColumns(chunkKeywordColumns)

	var valueTuples [][]any
	for chunkID, keywordIDs := range chunkIDKeywordIDsMap {
		for _, keywordID := range keywordIDs {
			tuple := []any{keywordID, chunkID}
			valueTuples = append(
				valueTuples,
				tuple,
			)
		}
	}

	values := valuesString(valueTuples)

	q := insertF(tx.CurrentSchema(), database.Table.ChunkKeyword, joinedColumns, values, "")
	common.Log.DBInfo(tx.ID, q)

	err := tx.Exec(q)
	if err != nil {
		return err
	}

	return nil
}
