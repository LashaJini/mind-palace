package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
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

func InsertResource(ctx context.Context, db *sql.DB, resource *OriginalResource) error {
	tx := database.NewMultiInstruction(ctx, db)
	if err := tx.Begin(); err != nil {
		return err
	}

	if err := InsertResourceTx(tx, resource); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func InsertResourceTx(tx *database.MultiInstruction, resource *OriginalResource) error {
	q := fmt.Sprintf(`
INSERT INTO original_resource (
	id,
	memory_id,
	file_path,
	created_at,
	updated_at
)
VALUES (
	$1, $2, $3, $4, $5
)`)

	return tx.Exec(q,
		resource.ID,
		resource.MemoryID,
		resource.FilePath,
		resource.CreatedAt,
		resource.UpdatedAt,
	)
}
