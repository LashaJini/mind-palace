package models

import (
	"context"
	"database/sql"
	"fmt"
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

func InsertMemory(ctx context.Context, db *sql.DB, memory *Memory) (uuid.UUID, error) {
	tx := database.NewMultiInstruction(ctx, db)
	if err := tx.Begin(); err != nil {
		return uuid.Nil, err
	}

	id, err := InsertMemoryTx(tx, memory)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func InsertMemoryTx(tx *database.MultiInstruction, memory *Memory) (uuid.UUID, error) {
	q := fmt.Sprintf(`
INSERT INTO memory (
	created_at,
	updated_at
)
VALUES (
	$1, $2
) 
RETURNING id`)

	var id string
	err := tx.QueryRow(q, memory.CreatedAt, memory.UpdatedAt).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(id)
}
