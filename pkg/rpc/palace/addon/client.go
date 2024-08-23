package addonrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/lashajini/mind-palace/pkg/common"
	pb "github.com/lashajini/mind-palace/pkg/rpc/gen"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const RETRY_COUNT = 15

type Client struct {
	client pb.AddonClient
}

func NewGrpcClient(cfg *common.Config) *Client {
	addr := fmt.Sprintf("localhost:%d", cfg.PALACE_GRPC_SERVER_PORT)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, _ := grpc.NewClient(addr, opts...)
	client := pb.NewAddonClient(conn)

	return &Client{client}
}

func (c *Client) Add(ctx context.Context, file string, steps []string) (<-chan *pb.AddonResult, error) {
	addonResultC := make(chan *pb.AddonResult)
	go func() {
		resource := &pb.Resource{
			File:  file,
			Steps: steps,
		}
		joinedAddons, _ := c.client.JoinAddons(ctx, resource)

		if joinedAddons != nil {
			for _, joinedAddon := range joinedAddons.Addons {
				tmp := &pb.JoinedAddons{
					File:   file,
					Addons: joinedAddon,
				}
				// server may decide that it's more efficient to join multiple addons together
				addonResult, _ := c.client.ApplyAddon(ctx, tmp)

				addonResultC <- addonResult
			}
		}

		close(addonResultC)
	}()

	return addonResultC, nil
}

func (c *Client) Ping(ctx context.Context) error {
	var err error
	for i := 1; i <= RETRY_COUNT; i++ {
		_, err := c.client.Ping(ctx, &pb.Empty{})
		if err != nil {
			loggers.Log.Warn(ctx, "grpc server ping '%d' failed (retrying in 1 sec), reason: %v", i, err)
		} else {
			loggers.Log.Info(ctx, "grpc server ping '%d' successful", i)
			err = nil
			break
		}
		time.Sleep(1 * time.Second)
	}

	return err
}
