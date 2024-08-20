package rpcclient

import (
	"context"
	"fmt"

	"github.com/lashajini/mind-palace/pkg/common"
	pb "github.com/lashajini/mind-palace/pkg/rpc/client/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LLMClient struct {
	client pb.LLMClient
}

func NewLLMClient(cfg *common.Config) *LLMClient {
	addr := fmt.Sprintf("localhost:%d", cfg.PALACE_GRPC_SERVER_PORT)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, _ := grpc.NewClient(addr, opts...)
	client := pb.NewLLMClient(conn)

	return &LLMClient{client}
}

func (c *LLMClient) SetConfig(ctx context.Context, cfg map[string]string) error {
	m := make(map[string]string)
	for k, v := range cfg {
		m[k] = v
	}

	pbCfg := &pb.LLMConfig{Map: m}

	_, err := c.client.SetConfig(ctx, pbCfg)
	return err
}
