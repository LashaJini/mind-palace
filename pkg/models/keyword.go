package models

import (
	"context"
	"fmt"
	"time"

	"github.com/lashajini/mind-palace/pkg/mperrors"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

type Keyword struct {
	ID        int
	Name      string
	CreatedAt int64
	UpdatedAt int64
}

var keywordColumns = []string{
	"id",
	"name",
	"created_at",
	"updated_at",
}

func InsertManyKeywordsTx(ctx context.Context, tx *database.MultiInstruction, keywords []string) (map[string]int, error) {
	if len(keywords) == 0 {
		return nil, mperrors.Onf("empty keywords")
	}

	joinedColumns, _ := joinColumns(keywordColumns, "id")

	now := time.Now().UTC().Unix()
	createdAt := now
	updatedAt := now

	var valueTuples [][]any
	for _, keyword := range keywords {
		tuple := []any{keyword, createdAt, updatedAt}
		valueTuples = append(
			valueTuples,
			tuple,
		)
	}

	values := valuesString(valueTuples)

	q := insertF(
		tx.CurrentSchema(),
		database.Table.Keyword,
		joinedColumns,
		values,
		fmt.Sprintf(`ON CONFLICT (name) DO UPDATE SET updated_at = %d RETURNING name, id`, now),
	)
	loggers.Log.DBInfo(ctx, tx.ID, q)

	rows, err := tx.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keywordIDs := make(map[string]int)
	for rows.Next() {
		var id int
		var name string
		_ = rows.Scan(&name, &id)
		keywordIDs[name] = id
	}

	return keywordIDs, nil
}
