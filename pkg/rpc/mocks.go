package rpcclient

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/lashajini/mind-palace/pkg/rpc/gen"
	"github.com/stretchr/testify/mock"
)

type MockGrpcClient struct {
	mock.Mock
}

func (m *MockGrpcClient) Search(ctx context.Context, text string) (*pb.SearchResponse, error) {
	args := m.Called(ctx, text)
	return args.Get(0).(*pb.SearchResponse), args.Error(1)
}

func (m *MockGrpcClient) Insert(ctx context.Context, ids []uuid.UUID, outputs []string, types []string) error {
	return nil
}

func (m *MockGrpcClient) Drop(ctx context.Context) error { return nil }
