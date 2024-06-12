package models

import (
	"fmt"
	"strings"
)

func excludeColumns(columns []string, names ...string) []string {
	remainingColumns := []string{}

	excludedNamesMap := make(map[string]struct{}, len(names))
	for _, name := range names {
		excludedNamesMap[name] = struct{}{}
	}

	for _, column := range columns {
		if _, found := excludedNamesMap[column]; !found {
			remainingColumns = append(remainingColumns, column)
		}
	}

	return remainingColumns
}

func joinColumns(columns []string, excluded ...string) (string, int) {
	joinedColumns := ""
	remainingColumns := excludeColumns(columns, excluded...)
	for _, column := range remainingColumns {
		joinedColumns += fmt.Sprintf("\t%s,\n", column)
	}
	joinedColumns = strings.TrimSuffix(joinedColumns, ",\n")
	return joinedColumns, len(remainingColumns)
}

func placeholdersString(numRows int, numColumns int) string {
	placeholders := ""
	i := 1
	for range numRows {
		placeholders += "\t("
		for j := i; j < i+numColumns; j++ {
			placeholders += fmt.Sprintf("$%d, ", j)
		}
		placeholders = strings.TrimSuffix(placeholders, ", ")
		placeholders += "),\n"

		i += numColumns
	}
	placeholders = strings.TrimSuffix(placeholders, ",\n")

	return placeholders
}

func insertF(table, columns, values, additional string) string {
	q := fmt.Sprintf(`
INSERT INTO %s (
%s
)
VALUES
%s
%s
`, table, columns, values, additional)

	return q
}
