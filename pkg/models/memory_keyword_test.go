package models

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/stretchr/testify/assert"
)

func (suite *ModelsTestSuite) Test_InsertManyMemoryKeywordsTx_success() {
	t := suite.T()
	tx := database.NewMultiInstruction(suite.db)

	err := tx.Begin()
	assert.NoError(t, err)
	t.Cleanup(suite.MemoryKeywordCleanup)

	memory := NewMemory()
	memoryID, err := InsertMemoryTx(tx, memory)
	memory.ID = memoryID
	assert.NoError(t, err)

	keywords := []string{
		"keyword1",
		"keyword2",
		"keyword3",
	}
	keywordIDs, err := InsertManyKeywordsTx(tx, keywords)
	assert.NoError(t, err)

	err = InsertManyMemoryKeywordsTx(tx, keywordIDs, memory.ID)
	assert.NoError(t, err)

	err = tx.Commit()
	assert.NoError(t, err)

	q := fmt.Sprintf(fmt.Sprintf("select keyword_id, memory_id from %s.%s order by keyword_id asc", suite.currentSchema, database.Table.MemoryKeyword))
	rows, err := suite.db.DB().Query(q)
	assert.NoError(t, err)
	defer rows.Close()

	totalRows := 0
	for rows.Next() {
		var (
			keywordID int
			memoryID  uuid.UUID
		)
		_ = rows.Scan(&keywordID, &memoryID)

		nthKeyword := keywords[totalRows]
		nthKeywordID := keywordIDs[nthKeyword]
		assert.Equal(t, nthKeywordID, keywordID)

		assert.Equal(t, memoryID, memory.ID)
		totalRows++
	}

	assert.Equal(t, len(keywords), totalRows)
}

func (suite *ModelsTestSuite) MemoryKeywordCleanup() {
	// cascades and deletes memory_keyword pairs as well
	suite.KeywordCleanup()

	_, err := suite.db.DB().Exec(fmt.Sprintf("delete from %s.%s", suite.currentSchema, database.Table.Memory))
	if err != nil {
		log.Fatal(err)
	}
}
