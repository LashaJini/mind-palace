package models

import (
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

	q := insertF(database.Table.Keyword, joinedColumns, values, "RETURNING id")
	common.Log.DBInfo(tx.ID, q)

	rows, err := tx.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	i := 0
	keywordIDs := make(map[string]int)
	for rows.Next() {
		defer rows.Close()
		var id int
		_ = rows.Scan(&id)
		keywordIDs[keywords[i]] = id
		i++
	}

	return keywordIDs, nil
}
