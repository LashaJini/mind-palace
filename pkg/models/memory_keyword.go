package models

import (
	"github.com/google/uuid"
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
	joinedColumns, numColumns := joinColumns(memoryKeywordColumns)
	placeholders := placeholdersString(len(keywords), numColumns)

	values := []any{}
	for _, keywordID := range keywords {
		values = append(
			values,
			keywordID,
			memoryID,
		)
	}

	q := insertF("memory_keyword", joinedColumns, placeholders, "")

	err := tx.Exec(q, values...)
	if err != nil {
		return err
	}

	return nil
}
