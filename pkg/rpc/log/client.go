package log

import (
	"context"
	"fmt"
	"runtime"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/common"
	rpcclient "github.com/lashajini/mind-palace/pkg/rpc"
	pb "github.com/lashajini/mind-palace/pkg/rpc/gen"
)

const RETRY_COUNT = 15

type Client struct {
	*rpcclient.Client[pb.LogClient, *Client]
	serviceName string
}

func NewGrpcClient(cfg *common.Config, serviceName string) *Client {
	c := &Client{serviceName: serviceName}

	client := rpcclient.NewGrpcClient(
		cfg.LOG_GRPC_SERVER_PORT,
		"Log",
		RETRY_COUNT,
		pb.NewLogClient,
		c,
	)

	c.Client = client

	ctx := context.Background()
	if err := c.Ping(ctx); err != nil {
		c.Logger.Fatal(ctx, err, "")
		panic(err)
	}

	return c
}

func (c *Client) request(logType string, id uuid.UUID, caller_incr int, format string, v ...interface{}) pb.LogRequest {
	// caller [2 + caller_incr] -> log type (Info, Debug, ...) [1] -> request [0]
	_, filename, line, _ := runtime.Caller(2 + caller_incr)

	msg := fmt.Sprintf(format, v...)
	return pb.LogRequest{
		Msg:         msg,
		Filename:    filename,
		Line:        int32(line),
		ServiceName: c.serviceName,
		Type:        logType,
		Id:          id.String(),
	}
}

func (c *Client) Debug(ctx context.Context, format string, v ...interface{}) (*pb.Empty, error) {
	request := c.request("debug", uuid.Nil, 0, format, v...)

	return c.Service.Message(ctx, &request)
}

func (c *Client) Warn(ctx context.Context, format string, v ...interface{}) (*pb.Empty, error) {
	request := c.request("warning", uuid.Nil, 0, format, v...)

	return c.Service.Message(ctx, &request)
}

func (c *Client) Error(ctx context.Context, err error, format string, v ...interface{}) (*pb.Empty, error) {
	_format := format
	if err != nil {
		_format = fmt.Sprintf("%s: %s", err.Error(), _format)
	}
	request := c.request("error", uuid.Nil, 0, _format, v...)

	return c.Service.Message(ctx, &request)
}

func (c *Client) Fatal(ctx context.Context, err error, format string, v ...interface{}) (*pb.Empty, error) {
	_format := format
	if err != nil {
		_format = fmt.Sprintf("%s: %s", err.Error(), _format)
	}
	request := c.request("fatal", uuid.Nil, 0, _format, v...)

	return c.Service.Message(ctx, &request)
}

func (c *Client) Info(ctx context.Context, format string, v ...interface{}) (*pb.Empty, error) {
	request := c.request("info", uuid.Nil, 0, format, v...)

	return c.Service.Message(ctx, &request)
}

func (c *Client) DBInfo(ctx context.Context, id uuid.UUID, format string, v ...interface{}) (*pb.Empty, error) {
	request := c.request("db_info", id, 0, format, v...)

	return c.Service.Message(ctx, &request)
}

// this is only called from inside MultiInstruction's methods,
// that's why the request's caller_inc is set to 1
func (c *Client) TXInfo(ctx context.Context, id uuid.UUID, format string, v ...interface{}) (*pb.Empty, error) {
	request := c.request("tx_info", id, 1, format, v...)

	return c.Service.Message(ctx, &request)
}
