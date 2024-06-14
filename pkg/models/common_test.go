package models

import (
	"reflect"
	"strings"
	"testing"
)

func Test_excludeColumns(t *testing.T) {
	tests := []struct {
		columns  []string
		names    []string
		expected []string
	}{
		{
			columns:  []string{"id", "name", "age"},
			names:    []string{"id", "age"},
			expected: []string{"name"},
		},
		{
			columns:  []string{"username", "email", "password"},
			names:    []string{"password"},
			expected: []string{"username", "email"},
		},
	}

	for _, test := range tests {
		result := excludeColumns(test.columns, test.names...)

		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Test failed. Columns: %v, Names: %v, Expected: %v, Got: %v", test.columns, test.names, test.expected, result)
		}
	}
}

func Test_joinColumns(t *testing.T) {
	tests := []struct {
		numRows    int
		numColumns int
		expected   string
	}{
		{
			numRows:    3,
			numColumns: 2,
			expected:   "($1, $2),($3, $4),($5, $6)",
		},
		{
			numRows:    2,
			numColumns: 3,
			expected:   "($1, $2, $3),($4, $5, $6)",
		},
	}

	for _, tt := range tests {
		result := placeholdersString(tt.numRows, tt.numColumns)

		if result != tt.expected {
			t.Errorf("Failed: numRows=%d, numColumns=%d.\nExpected: %s\ngot: %s", tt.numRows, tt.numColumns, tt.expected, result)
		}
	}
}

func Test_insertF(t *testing.T) {
	tests := []struct {
		name       string
		table      string
		columns    string
		values     string
		additional string
		expected   string
	}{
		{
			name:       "With values",
			table:      "users",
			columns:    "name, age",
			values:     "'Alice', 30",
			additional: "",
			expected:   strings.TrimSpace(`INSERT INTO users (name, age) VALUES 'Alice', 30`),
		},
		{
			table:      "products",
			columns:    "name, price",
			values:     "'Apple', 1.5",
			additional: "RETURNING id",
			expected:   strings.TrimSpace(`INSERT INTO products (name, price) VALUES 'Apple', 1.5 RETURNING id`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := insertF(tt.table, tt.columns, tt.values, tt.additional)

			if result != tt.expected {
				t.Errorf("Failed: %s.\nExpected: %s\ngot: %s", tt.name, tt.expected, result)
			}
		})
	}
}
