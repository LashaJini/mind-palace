package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

type Chunk struct {
	ID        uuid.UUID
	Sequence  int
	Chunk     string
	CreatedAt int64
	UpdatedAt int64
}

var chunkColumns = []string{
	"id",
	"memory_id",
	"sequence",
	"chunk",
	"created_at",
	"updated_at",
}

func InsertManyChunksTx(tx *database.MultiInstruction, memoryID uuid.UUID, chunks []string) ([]uuid.UUID, error) {
	if len(chunks) == 0 {
		return nil, errors.New("reason: empty chunks")
	}
	ctx := context.Background()

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

	rows, err := tx.Query(q)
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
