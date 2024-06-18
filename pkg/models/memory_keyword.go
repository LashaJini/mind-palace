package models

import (
	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

type MemoryKeyword struct {
	KeywordID int
	MemoryID  uuid.UUID
}

var memoryKeywordColumns = []string{
	"keyword_id",
	"memory_id",
}

func InsertManyMemoryKeywordsTx(tx *database.MultiInstruction, keywords map[string]int, memoryID uuid.UUID) error {
	joinedColumns, _ := joinColumns(memoryKeywordColumns)

	var valueTuples [][]any
	for _, keywordID := range keywords {
		tuple := []any{keywordID, memoryID}
		valueTuples = append(
			valueTuples,
			tuple,
		)
	}

	values := valuesString(valueTuples)

	q := insertF(tx.CurrentSchema(), database.Table.MemoryKeyword, joinedColumns, values, "")
	common.Log.DBInfo(tx.ID, q)

	err := tx.Exec(q)
	if err != nil {
		return err
	}

	return nil
}
