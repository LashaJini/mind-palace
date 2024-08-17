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

// returns comma separated column names
func joinColumns(columns []string, excluded ...string) (string, int) {
	joinedColumns := ""
	remainingColumns := excludeColumns(columns, excluded...)
	for _, column := range remainingColumns {
		joinedColumns += fmt.Sprintf("%s,", column)
	}
	joinedColumns = strings.TrimSuffix(joinedColumns, ",")
	return joinedColumns, len(remainingColumns)
}

func placeholdersString(numRows int, numColumns int) string {
	placeholders := ""
	i := 1
	for range numRows {
		placeholders += "("
		for j := i; j < i+numColumns; j++ {
			placeholders += fmt.Sprintf("$%d, ", j)
		}
		placeholders = strings.TrimSuffix(placeholders, ", ")
		placeholders += "),"

		i += numColumns
	}
	placeholders = strings.TrimSuffix(placeholders, ",")

	return placeholders
}

// not sanitized :)
func valuesString(tuples [][]any) string {
	values := ""
	for _, value := range tuples {
		values += "("
		for _, v := range value {
			s := v
			if _, ok := v.(string); ok {
				s = strings.ReplaceAll(s.(string), "'", "''")
			}
			values += fmt.Sprintf("'%v',", s)
		}
		values = strings.TrimSuffix(values, ",")
		values += "),"
	}
	values = strings.TrimSuffix(values, ",")

	return values
}

func insertF(schema, table, columns, values, additional string) string {
	q := strings.TrimSpace(
		fmt.Sprintf(`INSERT INTO %s.%s (%s) VALUES %s %s`,
			schema,
			table,
			columns,
			values,
			additional,
		),
	)

	return q
}
