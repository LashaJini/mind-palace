package models

import (
	"time"

	"github.com/google/uuid"
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
	createdAt := memory.CreatedAt
	updatedAt := memory.UpdatedAt

	joinedColumns, numColumns := joinColumns(memoryColumns, "id")
	placeholders := placeholdersString(1, numColumns)
	q := insertF(database.Table.Memory, joinedColumns, placeholders, "RETURNING id")

	var id string
	err := tx.QueryRow(q, createdAt, updatedAt).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(id)
}
