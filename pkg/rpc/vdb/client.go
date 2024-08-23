package vdbrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/common"
	pb "github.com/lashajini/mind-palace/pkg/rpc/gen"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const RETRY_COUNT = 15

type VDBGrpcClient struct {
	client pb.VDBClient

	user string
}

func NewVDBGrpcClient(cfg *common.Config, user string) *VDBGrpcClient {
	addr := fmt.Sprintf("localhost:%d", cfg.VDB_GRPC_SERVER_PORT)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, _ := grpc.NewClient(addr, opts...)
	client := pb.NewVDBClient(conn)

	return &VDBGrpcClient{client, user}
}

func (c *VDBGrpcClient) Insert(ctx context.Context, ids []uuid.UUID, outputs []string) error {
	var _ids []string
	for _, id := range ids {
		_ids = append(_ids, id.String())
	}

	vdbRows := &pb.VDBRows{
		User:   c.user,
		Ids:    _ids,
		Inputs: outputs,
	}

	_, err := c.client.Insert(ctx, vdbRows)

	return err
}

func (c *VDBGrpcClient) Ping(ctx context.Context) error {
	var err error
	for i := 1; i <= RETRY_COUNT; i++ {
		_, err := c.client.Ping(ctx, &pb.Empty{})
		if err != nil {
			loggers.Log.Warn(ctx, "vector database ping '%d' failed (retrying in 1 sec), reason: %v", i, err)
		} else {
			loggers.Log.Info(ctx, "vector database ping '%d' successful", i)
			err = nil
			break
		}
		time.Sleep(1 * time.Second)
	}

	return err
}

func (c *VDBGrpcClient) Drop(ctx context.Context) error {
	_, err := c.client.Drop(ctx, &pb.Empty{})

	return err
}
