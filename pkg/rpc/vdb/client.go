package vdbrpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/mperrors"
	rpcclient "github.com/lashajini/mind-palace/pkg/rpc"
	pb "github.com/lashajini/mind-palace/pkg/rpc/gen"
	"github.com/lashajini/mind-palace/pkg/rpc/log"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
)

const RETRY_COUNT = 15
const ROW_TYPE_CHUNK = common.ROW_TYPE_CHUNK
const ROW_TYPE_WHOLE = common.ROW_TYPE_WHOLE

type VDBClient interface {
	Search(ctx context.Context, text string) (*pb.SearchResponse, error)
	Insert(ctx context.Context, ids []uuid.UUID, outputs []string, types []string) error
	Drop(ctx context.Context) error
}

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

func (c *Client) Insert(ctx context.Context, ids []uuid.UUID, outputs []string, types []string) error {
	if len(ids) != len(outputs) && len(ids) != len(types) {
		return mperrors.Onf("length of ids, outputs and types must be same")
	}

	var vdbRows []*pb.InsertRequest_VDBRow
	for i, id := range ids {
		vdbRow := &pb.InsertRequest_VDBRow{
			Id:    id.String(),
			Input: outputs[i],
			Type:  types[i],
		}

		vdbRows = append(vdbRows, vdbRow)
	}

	insertRequest := &pb.InsertRequest{
		User: c.user,
		Rows: vdbRows,
	}

	_, err := c.Service.Insert(ctx, insertRequest)

	return err
}

func (c *Client) Search(ctx context.Context, text string) (*pb.SearchResponse, error) {
	vdbSearchRequest := &pb.SearchRequest{
		Text: text, User: c.user,
	}

	return c.Service.Search(ctx, vdbSearchRequest)
}

func (c *Client) Drop(ctx context.Context) error {
	_, err := c.Service.Drop(ctx, &pb.Empty{})

	return err
}
