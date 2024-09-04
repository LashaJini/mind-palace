package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

type OriginalResource struct {
	ID        uuid.UUID
	MemoryID  uuid.UUID
	FilePath  string
	CreatedAt int64
	UpdatedAt int64
}

func NewResource(id uuid.UUID, memoryID uuid.UUID, filepath string) *OriginalResource {
	now := time.Now().UTC().Unix()

	return &OriginalResource{
		ID:        id,
		MemoryID:  memoryID,
		FilePath:  filepath,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

var originalResourceColumns = []string{
	"id",
	"memory_id",
	"file_path",
	"created_at",
	"updated_at",
}

func InsertResourceTx(ctx context.Context, tx *database.MultiInstruction, resource *OriginalResource) error {
	joinedColumns, _ := joinColumns(originalResourceColumns)

	var valueTuples [][]any
	valueTuple := []any{
		resource.ID,
		resource.MemoryID,
		resource.FilePath,
		resource.CreatedAt,
		resource.UpdatedAt,
	}
	valueTuples = append(valueTuples, valueTuple)

	values := valuesString(valueTuples)

	q := insertF(tx.CurrentSchema(), database.Table.OriginalResource, joinedColumns, values, "")
	loggers.Log.DBInfo(ctx, tx.ID, q)

	return tx.Exec(ctx, q)
}
