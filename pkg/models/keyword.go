package models

import (
	"time"

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
	joinedColumns, numColumns := joinColumns(keywordColumns, "id")

	now := time.Now().UTC().Unix()
	createdAt := now
	updatedAt := now

	placeholders := placeholdersString(len(keywords), numColumns)

	values := []any{}
	for _, keyword := range keywords {
		values = append(
			values,
			keyword,
			createdAt,
			updatedAt,
		)
	}

	q := insertF(database.Table.Keyword, joinedColumns, placeholders, "RETURNING id")

	rows, err := tx.Query(q, values...)
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
