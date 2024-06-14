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
}

func NewClient(cfg *common.Config) *Client {
	addr := fmt.Sprintf("localhost:%d", cfg.GRPC_SERVER_PORT)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, _ := grpc.NewClient(addr, opts...)
	client := pb.NewPalaceClient(conn)

	return &Client{client}
}

func (c *Client) Add(ctx context.Context, file string, userCfg *mpuser.Config) (<-chan *pb.AddonResult, error) {
	addonResultC := make(chan *pb.AddonResult)
	go func() {
		resource := &pb.Resource{
			File:  file,
			Steps: userCfg.Steps(),
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

func (c *Client) VDBInsert(ctx context.Context, memoryID uuid.UUID, output string) error {
	vdbRow := &pb.VDBRow{
		Id:    memoryID.String(),
		Input: output,
	}

	_, err := c.client.VDBInsert(ctx, vdbRow)

	return err
}
