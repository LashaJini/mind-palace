package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

type Memory struct {
	ID        uuid.UUID
	CreatedAt int64
	UpdatedAt int64
}

func NewMemory() *Memory {
	now := time.Now().UTC().Unix()

	return &Memory{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

var memoryColumns = []string{
	"id",
	"created_at",
	"updated_at",
}

func InsertMemoryTx(tx *database.MultiInstruction, memory *Memory) (uuid.UUID, error) {
	ctx := context.Background()
	createdAt := memory.CreatedAt
	updatedAt := memory.UpdatedAt

	joinedColumns, _ := joinColumns(memoryColumns, "id")

	var valueTuples [][]any
	valueTuple := []any{createdAt, updatedAt}
	valueTuples = append(valueTuples, valueTuple)

	values := valuesString(valueTuples)

	q := insertF(tx.CurrentSchema(), database.Table.Memory, joinedColumns, values, "RETURNING id")
	loggers.Log.DBInfo(ctx, tx.ID, q)

	var id string
	err := tx.QueryRow(q).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(id)
}
