package rpcclient

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	pb "github.com/lashajini/mind-palace/pkg/rpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Pinger interface {
	Ping(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.Empty, error)
}

type Logger interface {
	Debug(ctx context.Context, format string, v ...interface{}) (*pb.Empty, error)
	Warn(ctx context.Context, format string, v ...interface{}) (*pb.Empty, error)
	Error(ctx context.Context, err error, format string, v ...interface{}) (*pb.Empty, error)
	Fatal(ctx context.Context, err error, format string, v ...interface{}) (*pb.Empty, error)
	Info(ctx context.Context, format string, v ...interface{}) (*pb.Empty, error)
	DBInfo(ctx context.Context, id uuid.UUID, format string, v ...interface{}) (*pb.Empty, error)
	TXInfo(ctx context.Context, id uuid.UUID, format string, v ...interface{}) (*pb.Empty, error)
}

type FnPbClient[T any] func(cc grpc.ClientConnInterface) T
type Client[T Pinger, L Logger] struct {
	Service T
	Logger  L

	name       string
	retryCount int
}

func NewGrpcClient[T Pinger, L Logger](port int, name string, retryCount int, fnPbClient FnPbClient[T], logger L) *Client[T, L] {
	addr := fmt.Sprintf("localhost:%d", port)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		logger.Fatal(context.Background(), err, "")
		panic(err)
	}
	service := fnPbClient(conn)

	c := &Client[T, L]{Service: service, name: name, retryCount: retryCount, Logger: logger}

	return c
}

func (c *Client[T, L]) Ping(ctx context.Context) error {
	var err error
	for i := 1; i <= c.retryCount; i++ {
		_, err = c.Service.Ping(ctx, &pb.Empty{})
		if err != nil {
			c.Logger.Warn(ctx, "%s grpc server ping '%d' failed (retrying in 1 sec), reason: %v", c.name, i, err)
		} else {
			c.Logger.Info(ctx, "%s grpc server ping '%d' successful", c.name, i)
			err = nil
			break
		}
		time.Sleep(1 * time.Second)
	}

	return err
}
