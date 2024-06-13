package models

import (
	"fmt"

	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/stretchr/testify/assert"
)

func (suite *ModelsTestSuite) Test_InsertManyKeywordsTx_success() {
	t := suite.T()
	tx := database.NewMultiInstruction(suite.ctx, suite.db.DB())

	t.Cleanup(suite.KeywordCleanup)

	err := tx.Begin()
	assert.NoError(t, err)

	keywords := []string{"keyword1", "keyword2", "keyword3"}
	keywordIDs, _ := InsertManyKeywordsTx(tx, keywords)
	expectedIDs := map[string]int{
		"keyword1": 1,
		"keyword2": 2,
		"keyword3": 3,
	}

	err = tx.Commit()
	assert.NoError(t, err)

	if len(keywordIDs) != len(expectedIDs) {
		t.Fatalf("Unexpected number of keyword IDs. Expected %d, got %d", len(expectedIDs), len(keywordIDs))
	}
	for keyword, id := range expectedIDs {
		if retrievedID, ok := keywordIDs[keyword]; !ok || retrievedID != id {
			t.Fatalf("Unexpected keyword ID for keyword %s. Expected %d, got %d", keyword, id, retrievedID)
		}
	}

	rows, err := suite.db.DB().Query(fmt.Sprintf("select id, name from %s", database.Table.Keyword))
	assert.NoError(t, err)
	defer rows.Close()

	for rows.Next() {
		var (
			id   int
			name string
		)

		err := rows.Scan(&id, &name)
		assert.NoError(t, err)

		assert.Equal(t, expectedIDs[name], id)
	}
}

func (suite *ModelsTestSuite) KeywordCleanup() {
	_, err := suite.db.DB().Exec(fmt.Sprintf("delete from %s", database.Table.Keyword))
	if err != nil {
		fmt.Println(err)
	}
}
