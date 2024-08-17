package rpcclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/mpuser"
	pb "github.com/lashajini/mind-palace/pkg/rpc/client/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
