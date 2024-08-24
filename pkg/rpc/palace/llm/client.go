package llmrpc

import (
	"context"

	"github.com/lashajini/mind-palace/pkg/common"
	rpcclient "github.com/lashajini/mind-palace/pkg/rpc"
	pb "github.com/lashajini/mind-palace/pkg/rpc/gen"
	"github.com/lashajini/mind-palace/pkg/rpc/log"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
)

const RETRY_COUNT = 15

type Client struct {
	*rpcclient.Client[pb.LLMClient, *log.Client]
}

func NewGrpcClient(cfg *common.Config) *Client {
	client := rpcclient.NewGrpcClient(
		cfg.PALACE_GRPC_SERVER_PORT,
		"Palace(LLM)",
		RETRY_COUNT,
		pb.NewLLMClient,
		loggers.Log,
	)
	c := &Client{client}

	ctx := context.Background()
	if err := c.Ping(ctx); err != nil {
		c.Logger.Fatal(ctx, err, "")
		panic(err)
	}

	return c
}

func (c *Client) SetConfig(ctx context.Context, cfg map[string]string) error {
	m := make(map[string]string)
	for k, v := range cfg {
		m[k] = v
	}

	pbCfg := &pb.LLMConfig{Map: m}

	_, err := c.Service.SetConfig(ctx, pbCfg)
	return err
}
