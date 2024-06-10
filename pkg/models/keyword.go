package models

import (
	"fmt"
	"strings"
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

func excludedColumns(names ...string) []string {
	remainingColumns := []string{}

	excludedNamesMap := make(map[string]struct{}, len(names))
	for _, name := range names {
		excludedNamesMap[name] = struct{}{}
	}

	for _, column := range keywordColumns {
		if _, found := excludedNamesMap[column]; !found {
			remainingColumns = append(remainingColumns, column)
		}
	}

	return remainingColumns
}

func InsertManyKeywordsTx(tx *database.MultiInstruction, keywords []string) (map[string]int, error) {
	joinedColumns := ""
	columns := excludedColumns("id")
	numColumns := len(columns)
	for _, column := range columns {
		joinedColumns += fmt.Sprintf("\t%s,\n", column)
	}
	joinedColumns = strings.TrimSuffix(joinedColumns, ",\n")

	now := time.Now().UTC().Unix()
	createdAt := now
	updatedAt := now
	placeholders := ""
	values := []any{}
	i := 1
	for _, keyword := range keywords {
		placeholders += "\t("
		for j := i; j < i+numColumns; j++ {
			placeholders += fmt.Sprintf("$%d, ", j)
		}
		placeholders = strings.TrimSuffix(placeholders, ", ")
		placeholders += "),\n"

		values = append(
			values,
			keyword,
			createdAt,
			updatedAt,
		)

		i += numColumns
	}
	placeholders = strings.TrimSuffix(placeholders, ",\n")

	q := fmt.Sprintf(`
INSERT INTO keyword (
%s
)
VALUES
%s
RETURNING id
`, joinedColumns, placeholders)

	rows, err := tx.Query(q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	i = 0
	keywordIDs := make(map[string]int)
	for rows.Next() {
		var id int
		_ = rows.Scan(&id)
		keywordIDs[keywords[i]] = id
		i++
	}

	return keywordIDs, nil
}
