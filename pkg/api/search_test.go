package api

import (
	"context"

	"github.com/google/uuid"
	rpcclient "github.com/lashajini/mind-palace/pkg/rpc"
	pb "github.com/lashajini/mind-palace/pkg/rpc/gen"
	vdbrpc "github.com/lashajini/mind-palace/pkg/rpc/vdb"
	"github.com/lashajini/mind-palace/pkg/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *APITestSuite) TestSearch_Success() {
	t := suite.T()

	mockGrpcClient := new(rpcclient.MockGrpcClient)
	mockDB := new(database.MockDB)
	ctx := context.Background()

	ids := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}
	searchResponse := &pb.SearchResponse{
		Rows: []*pb.SearchResponse_VDBRow{
			{
				Id:   ids[0],
				Type: vdbrpc.ROW_TYPE_WHOLE,
			},
			{
				Id:   ids[1],
				Type: vdbrpc.ROW_TYPE_WHOLE,
			},
			{
				Id:   ids[2],
				Type: vdbrpc.ROW_TYPE_CHUNK,
			},
		},
	}
	mockGrpcClient.On("Search", ctx, "sample text").Return(searchResponse, nil)

	mockRows := &database.MockRows{
		Values: [][]interface{}{
			{"memory_id_1"},
		},
	}
	mockDB.On("Query", mock.Anything, mock.Anything).Return(mockRows, nil)

	result, err := Search(ctx, "sample text", mockDB, mockGrpcClient)

	assert.NoError(t, err)
	assert.Len(t, result.Response, len(searchResponse.Rows))
	assert.Contains(t, ids, result.Response[0].MemoryID)
}
