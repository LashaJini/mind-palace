package addonrpc

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
	*rpcclient.Client[pb.AddonClient, *log.Client]
}

func NewGrpcClient(cfg *common.Config) *Client {
	client := rpcclient.NewGrpcClient(
		cfg.PALACE_GRPC_SERVER_PORT,
		"Palace(Addon)",
		RETRY_COUNT,
		pb.NewAddonClient,
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

func (c *Client) Add(ctx context.Context, file string, steps []string) (<-chan *pb.AddonResult, error) {
	addonResultC := make(chan *pb.AddonResult)
	go func() {
		resource := &pb.Resource{
			File:  file,
			Steps: steps,
		}
		joinedAddons, _ := c.Service.JoinAddons(ctx, resource)

		select {
		case <-ctx.Done():
			close(addonResultC)
			return
		default:
			if joinedAddons != nil {
				for _, joinedAddon := range joinedAddons.Addons {
					tmp := &pb.JoinedAddons{
						File:   file,
						Addons: joinedAddon,
					}
					// server may decide that it's more efficient to join multiple addons together
					addonResult, _ := c.Service.ApplyAddon(ctx, tmp)

					addonResultC <- addonResult
				}
			}
		}

		close(addonResultC)
	}()

	return addonResultC, ctx.Err()
}
