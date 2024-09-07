package database

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Query(query string, args ...interface{}) (SQLRows, error) {
	_args := m.Called(query, args)
	return _args.Get(0).(SQLRows), _args.Error(1)
}

func (m *MockDB) CurrentSchema() string {
	return "MOCKED_DB_SCHEMA"
}

func (m *MockDB) Close() error {
	return nil
}

func (m *MockDB) SetSchema(schema string) {
	m.Called(schema)
}

// Mocking the database rows
type MockRows struct {
	mock.Mock
	index  int
	Values [][]interface{}
}

func (rows *MockRows) Next() bool {
	// Increment the index and check if there are more rows
	if rows.index < len(rows.Values) {
		rows.index++
		return true
	}
	return false
}

func (rows *MockRows) Scan(dest ...interface{}) error {
	// Assign values to dest based on the current index
	for i, v := range rows.Values[rows.index-1] {
		switch d := dest[i].(type) {
		case *string:
			*d = v.(string)
		}
	}
	return nil
}

func (rows *MockRows) Close() error                            { return nil }
func (rows *MockRows) ColumnTypes() ([]*sql.ColumnType, error) { return nil, nil }
func (rows *MockRows) Columns() ([]string, error)              { return nil, nil }
func (rows *MockRows) Err() error                              { return nil }
func (rows *MockRows) NextResultSet() bool                     { return false }
