package rpcclient

import (
	"context"
	"fmt"

	pb "github.com/lashajini/mind-palace/rpc/client/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client pb.PalaceClient
}

func NewClient(port int) *Client {
	addr := fmt.Sprintf("localhost:%d", port)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, _ := grpc.NewClient(addr, opts...)
	client := pb.NewPalaceClient(conn)

	return &Client{client}
}

func (c *Client) Add(ctx context.Context, memory *pb.Memory) (*pb.Status, error) {
	return c.client.Add(ctx, memory)
}
