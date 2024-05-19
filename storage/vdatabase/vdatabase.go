package vdatabase

import (
	"context"
	"fmt"
	"log"

	"github.com/lashajini/mind-palace/config"
	"github.com/lashajini/mind-palace/constants"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

var LLM_DIM = "384"
var NLIST = 1536

type MindPalaceVDB struct {
	vdb *client.Client
}

func InitVDB(cfg *config.Config) *MindPalaceVDB {
	collectionName := constants.VDB_ORIGINAL_RESOURCE_COLLECTION_NAME
	ctx := context.Background()
	cli, err := client.NewGrpcClient(ctx, cfg.VDBAddr())
	if err != nil {
		log.Fatal(err)
	}

	var dbExists bool
	vdbStrs, _ := cli.ListDatabases(ctx)
	for _, vdbStr := range vdbStrs {
		if vdbStr.Name == cfg.VDB_NAME {
			dbExists = true
			break
		}
	}

	if !dbExists {
		fmt.Printf("Vector database '%s' does not exist. Creating...\n", cfg.VDB_NAME)
		if err := cli.CreateDatabase(ctx, cfg.VDB_NAME); err != nil {
			log.Fatal(err)
		}
	}

	cli.UsingDatabase(ctx, cfg.VDB_NAME)

	originalCollectionExists, err := cli.HasCollection(ctx, collectionName)
	if err != nil {
		log.Fatal(err)
	}

	if !originalCollectionExists {
		fmt.Printf("Collection '%s' does not exist. Creating...\n", collectionName)
		schema := &entity.Schema{
			CollectionName: collectionName,
			Description:    "Original resource collection",
			Fields: []*entity.Field{
				{
					Name:       "id",
					DataType:   entity.FieldTypeVarChar,
					PrimaryKey: true,
					AutoID:     false,
					TypeParams: map[string]string{
						"max_length": "64",
					},
				},
				{
					Name:     "vector",
					DataType: entity.FieldTypeFloatVector,
					TypeParams: map[string]string{
						"dim": LLM_DIM,
					},
				},
			},
		}

		if err := cli.CreateCollection(ctx, schema, 2); err != nil {
			log.Fatal(err)
		}

		vectorIdx, err := entity.NewIndexIvfFlat(entity.L2, NLIST)
		if err != nil {
			cli.DropCollection(ctx, collectionName)
			log.Fatal(err)
		}

		err = cli.CreateIndex(ctx, collectionName, "vector", vectorIdx, false)
		if err != nil {
			cli.DropCollection(ctx, collectionName)
			log.Fatal(err)
		}
	}

	return &MindPalaceVDB{vdb: &cli}
}
