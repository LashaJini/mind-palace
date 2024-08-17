package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

var summaryColumns = []string{
	"id",
	"memory_id",
	"text",
	"created_at",
	"updated_at",
}

func InsertSummaryTx(tx *database.MultiInstruction, memoryID, summaryID uuid.UUID, summary string) error {
	if len(summary) == 0 {
		return errors.New("reason: empty summary")
	}

	now := time.Now().UTC().Unix()
	createdAt := now
	updatedAt := now

	var valueTuples [][]any
	valueTuple := []any{summaryID, memoryID, summary, createdAt, updatedAt}
	valueTuples = append(valueTuples, valueTuple)

	joinedColumns, _ := joinColumns(summaryColumns)
	values := valuesString(valueTuples)
	q := insertF(tx.CurrentSchema(), database.Table.Summary, joinedColumns, values, "")
	common.Log.DBInfo(tx.ID, q)

	return tx.Exec(q)
}
