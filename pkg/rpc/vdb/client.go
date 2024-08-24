package vdbrpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/common"
	rpcclient "github.com/lashajini/mind-palace/pkg/rpc"
	pb "github.com/lashajini/mind-palace/pkg/rpc/gen"
	"github.com/lashajini/mind-palace/pkg/rpc/log"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
)

const RETRY_COUNT = 15

type Client struct {
	*rpcclient.Client[pb.VDBClient, *log.Client]

	user string
}

func NewGrpcClient(cfg *common.Config, user string) *Client {
	client := rpcclient.NewGrpcClient(
		cfg.VDB_GRPC_SERVER_PORT,
		"VDB",
		RETRY_COUNT,
		pb.NewVDBClient,
		loggers.Log,
	)
	c := &Client{client, user}

	ctx := context.Background()
	if err := c.Ping(ctx); err != nil {
		c.Logger.Fatal(ctx, err, "")
		panic(err)
	}

	return c
}

func (c *Client) Insert(ctx context.Context, ids []uuid.UUID, outputs []string) error {
	var _ids []string
	for _, id := range ids {
		_ids = append(_ids, id.String())
	}

	vdbRows := &pb.VDBRows{
		User:   c.user,
		Ids:    _ids,
		Inputs: outputs,
	}

	_, err := c.Service.Insert(ctx, vdbRows)

	return err
}

func (c *Client) Drop(ctx context.Context) error {
	_, err := c.Service.Drop(ctx, &pb.Empty{})

	return err
}
