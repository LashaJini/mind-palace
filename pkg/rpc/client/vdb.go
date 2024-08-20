package rpcclient

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/mpuser"
	pb "github.com/lashajini/mind-palace/pkg/rpc/client/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type VDBGrpcClient struct {
	client pb.VDBClient

	userCfg *mpuser.Config
}

func NewVDBGrpcClient(cfg *common.Config, userCfg *mpuser.Config) *VDBGrpcClient {
	addr := fmt.Sprintf("localhost:%d", cfg.VDB_GRPC_SERVER_PORT)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, _ := grpc.NewClient(addr, opts...)
	client := pb.NewVDBClient(conn)

	return &VDBGrpcClient{client, userCfg}
}

func (c *VDBGrpcClient) Insert(ctx context.Context, ids []uuid.UUID, outputs []string) error {
	var _ids []string
	for _, id := range ids {
		_ids = append(_ids, id.String())
	}

	vdbRows := &pb.VDBRows{
		User:   c.userCfg.Config.User,
		Ids:    _ids,
		Inputs: outputs,
	}

	_, err := c.client.Insert(ctx, vdbRows)

	return err
}

func (c *VDBGrpcClient) Ping(ctx context.Context) error {
	var err error
	for i := 1; i <= RETRY_COUNT; i++ {
		_, err := c.client.Ping(ctx, &pb.Empty{})
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

func (c *VDBGrpcClient) Drop(ctx context.Context) error {
	_, err := c.client.Drop(ctx, &pb.Empty{})

	return err
}
