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

type Chunk struct {
	ID        uuid.UUID
	MemoryID  uuid.UUID
	Sequence  int
	Chunk     string
	CreatedAt int64
	UpdatedAt int64
}

func (c Chunk) String() string {
	return fmt.Sprintf("Chunk{\n\tID: %s,\n\tMemoryID: %s,\n\tSequence: %d,\n\tChunk: %s,\n\tCreatedAt: %d,\n\tUpdatedAt: %d\n}",
		c.ID,
		c.MemoryID,
		c.Sequence,
		c.Chunk,
		c.CreatedAt,
		c.UpdatedAt,
	)
}

var chunkColumns = []string{
	"id",
	"memory_id",
	"sequence",
	"chunk",
	"created_at",
	"updated_at",
}

func InsertManyChunksTx(ctx context.Context, tx *database.MultiInstruction, memoryID uuid.UUID, chunks []string) ([]uuid.UUID, error) {
	if len(chunks) == 0 {
		return nil, mperrors.Onf("empty chunks")
	}

	joinedColumns, _ := joinColumns(chunkColumns, "id")

	now := time.Now().UTC().Unix()
	createdAt := now
	updatedAt := now

	var valueTuples [][]any
	for sequence, chunk := range chunks {
		tuple := []any{memoryID, sequence, chunk, createdAt, updatedAt}
		valueTuples = append(
			valueTuples,
			tuple,
		)
	}

	values := valuesString(valueTuples)

	q := insertF(tx.CurrentSchema(), database.Table.Chunk, joinedColumns, values, "RETURNING id")
	loggers.Log.DBInfo(ctx, tx.ID, q)

	rows, err := tx.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunkIDs []uuid.UUID
	for rows.Next() {
		var id string
		_ = rows.Scan(&id)

		parsedID, err := uuid.Parse(id)
		if err != nil {
			return chunkIDs, err
		}

		chunkIDs = append(chunkIDs, parsedID)
	}

	return chunkIDs, nil
}
