package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/lashajini/mind-palace/pkg/common"
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

func InsertManyKeywordsTx(tx *database.MultiInstruction, keywords []string) (map[string]int, error) {
	if len(keywords) == 0 {
		return nil, errors.New("reason: empty keywords")
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
	common.Log.DBInfo(tx.ID, q)

	rows, err := tx.Query(q)
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
