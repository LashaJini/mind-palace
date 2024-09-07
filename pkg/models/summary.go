package models

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/mperrors"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

type Summary struct {
	ID        uuid.UUID
	MemoryID  uuid.UUID
	Text      string
	CreatedAt int64
	UpdatedAt int64
}

func (s Summary) String() string {
	return fmt.Sprintf("Summary{\n\tID: %s,\n\tMemoryID: %s,\n\tText: %s,\n\tCreatedAt: %d,\n\tUpdatedAt: %d\n}",
		s.ID,
		s.MemoryID,
		s.Text,
		s.CreatedAt,
		s.UpdatedAt,
	)
}

var summaryColumns = []string{
	"id",
	"memory_id",
	"text",
	"created_at",
	"updated_at",
}

func InsertSummaryTx(ctx context.Context, tx *database.MultiInstruction, memoryID, summaryID uuid.UUID, summary string) error {
	if len(summary) == 0 {
		return mperrors.Onf("empty summary")
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
	loggers.Log.DBInfo(ctx, tx.ID, q)

	return tx.Exec(ctx, q)
}
