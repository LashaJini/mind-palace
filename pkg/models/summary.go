package models

import (
	"time"

	"github.com/google/uuid"
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
	now := time.Now().UTC().Unix()
	createdAt := now
	updatedAt := now

	joinedColumns, numColumns := joinColumns(summaryColumns)
	placeholders := placeholdersString(1, numColumns)
	q := insertF(database.Table.Summary, joinedColumns, placeholders, "")

	return tx.Exec(q, summaryID, memoryID, summary, createdAt, updatedAt)
}
