package models

import (
	"fmt"

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

func InsertManyMemoryKeywordsTx(tx *database.MultiInstruction, keywords map[string]int, memoryID uuid.UUID) (map[string]int, error) {
	fmt.Println(keywords)
	return nil, nil
}
