package rpcclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/config"
	pb "github.com/lashajini/mind-palace/pkg/rpc/client/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client pb.PalaceClient
}

func NewClient(cfg *config.Config) *Client {
	addr := fmt.Sprintf("localhost:%d", cfg.GRPC_SERVER_PORT)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, _ := grpc.NewClient(addr, opts...)
	client := pb.NewPalaceClient(conn)

	return &Client{client}
}

func (c *Client) Add(ctx context.Context, file string, id uuid.UUID, userCfg *config.UserConfig) (*pb.AddonResult, error) {
	for _, step := range userCfg.Steps() {
		addonResult, _ := c.client.ApplyAddon(ctx, &pb.Resource{
			File: file,
			Id:   id.String(),
			Step: step,
		})

		fmt.Println("step:", step)
		if addonResult != nil {
			fmt.Println("result:", addonResult.Data)
		}
	}

	return nil, nil
}
