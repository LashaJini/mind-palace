package models

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/stretchr/testify/assert"
)

func (suite *ModelsTestSuite) Test_InsertSummaryTx_success() {
	t := suite.T()
	tx := database.NewMultiInstruction(suite.ctx, suite.db)

	t.Cleanup(suite.SummaryCleanup)

	err := tx.Begin()
	assert.NoError(t, err)

	memory := NewMemory()
	memoryID, err := InsertMemoryTx(tx, memory)
	memory.ID = memoryID
	assert.NoError(t, err)

	summaryID := uuid.New()
	summary := "this is a summary"
	err = InsertSummaryTx(tx, memoryID, summaryID, summary)
	assert.NoError(t, err)

	err = tx.Commit()
	assert.NoError(t, err)

	q := fmt.Sprintf("select memory_id, text from %s.%s where id = '%s'", suite.currentSchema, database.Table.Summary, summaryID.String())
	rows, err := suite.db.DB().Query(q)
	assert.NoError(t, err)
	defer rows.Close()

	totalRows := 0
	for rows.Next() {
		var (
			memoryID uuid.UUID
			text     string
		)

		err := rows.Scan(&memoryID, &text)
		assert.NoError(t, err)

		assert.Equal(t, memory.ID.String(), memoryID.String())
		assert.Equal(t, summary, text)

		totalRows++
	}

	assert.Equal(t, totalRows, 1)
}

func (suite *ModelsTestSuite) SummaryCleanup() {
	_, err := suite.db.DB().Exec(fmt.Sprintf("delete from %s.%s", suite.currentSchema, database.Table.Summary))
	if err != nil {
		log.Fatal(err)
	}
}
