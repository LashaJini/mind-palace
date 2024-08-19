package rpcclient

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/mpuser"
	pb "github.com/lashajini/mind-palace/pkg/rpc/client/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const RETRY_COUNT = 15

type Client struct {
	client pb.PalaceClient

	userCfg *mpuser.Config
}

func NewClient(cfg *common.Config, userCfg *mpuser.Config) *Client {
	addr := fmt.Sprintf("localhost:%d", cfg.GRPC_SERVER_PORT)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, _ := grpc.NewClient(addr, opts...)
	client := pb.NewPalaceClient(conn)

	return &Client{client, userCfg}
}

func (c *Client) Add(ctx context.Context, file string) (<-chan *pb.AddonResult, error) {
	addonResultC := make(chan *pb.AddonResult)
	go func() {
		resource := &pb.Resource{
			File:  file,
			Steps: c.userCfg.Steps(),
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

func (c *Client) VDBInsert(ctx context.Context, ids []uuid.UUID, outputs []string) error {
	var _ids []string
	for _, id := range ids {
		_ids = append(_ids, id.String())
	}

	vdbRows := &pb.VDBRows{
		User:   c.userCfg.Config.User,
		Ids:    _ids,
		Inputs: outputs,
	}

	_, err := c.client.VDBInsert(ctx, vdbRows)

	return err
}

func (c *Client) VDBPing(ctx context.Context) error {
	var err error
	for i := 1; i <= RETRY_COUNT; i++ {
		_, err := c.client.VDBPing(ctx, &pb.Empty{})
		if err != nil {
			common.Log.Warn().Msgf("vector database ping '%d' failed (retrying in 1 sec), reason: %v", i, err)
		} else {
			common.Log.Info().Msgf("vector database ping '%d' successful", i)
			err = nil
			break
		}
		time.Sleep(1 * time.Second)
	}

	return err
}

func (c *Client) Ping(ctx context.Context) error {
	var err error
	for i := 1; i <= RETRY_COUNT; i++ {
		_, err := c.client.Ping(ctx, &pb.Empty{})
		if err != nil {
			common.Log.Warn().Msgf("grpc server ping '%d' failed (retrying in 1 sec), reason: %v", i, err)
		} else {
			common.Log.Info().Msgf("grpc server ping '%d' successful", i)
			err = nil
			break
		}
		time.Sleep(1 * time.Second)
	}

	return err
}

func (c *Client) VDBDrop(ctx context.Context) error {
	_, err := c.client.VDBDrop(ctx, &pb.Empty{})

	return err
}
