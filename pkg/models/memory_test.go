package models

import (
	"fmt"
	"log"

	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/stretchr/testify/assert"
)

func (suite *ModelsTestSuite) Test_InsertMemoryTx_success() {
	t := suite.T()
	tx := database.NewMultiInstruction(suite.ctx, suite.db)

	t.Cleanup(suite.MemoryCleanup)

	err := tx.Begin()
	assert.NoError(t, err)

	memory := NewMemory()
	id, err := InsertMemoryTx(tx, memory)
	assert.NoError(t, err)

	err = tx.Commit()
	assert.NoError(t, err)

	q := fmt.Sprintf("select created_at, updated_at from %s.%s where id = '%s'", suite.currentSchema, database.Table.Memory, id.String())
	rows, err := suite.db.DB().Query(q)
	assert.NoError(t, err)
	defer rows.Close()

	totalRows := 0
	for rows.Next() {
		totalRows++
	}

	assert.Equal(t, totalRows, 1)
}

func (suite *ModelsTestSuite) MemoryCleanup() {
	_, err := suite.db.DB().Exec(fmt.Sprintf("delete from %s.%s", suite.currentSchema, database.Table.Memory))
	if err != nil {
		log.Fatal(err)
	}
}
