package log

import (
	"context"
	"fmt"
	"runtime"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/common"
	pb "github.com/lashajini/mind-palace/pkg/rpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LogClient struct {
	client      pb.LogClient
	serviceName string
}

func NewLogClient(cfg *common.Config, serviceName string) *LogClient {
	addr := fmt.Sprintf("localhost:%d", cfg.LOG_GRPC_SERVER_PORT)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, _ := grpc.NewClient(addr, opts...)
	client := pb.NewLogClient(conn)

	return &LogClient{client, serviceName}
}

func (c *LogClient) request(logType string, id uuid.UUID, caller_incr int, format string, v ...interface{}) pb.LogRequest {
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

func (c *LogClient) Debug(ctx context.Context, format string, v ...interface{}) (*pb.Empty, error) {
	request := c.request("debug", uuid.Nil, 0, format, v...)

	return c.client.Message(ctx, &request)
}

func (c *LogClient) Warn(ctx context.Context, format string, v ...interface{}) (*pb.Empty, error) {
	request := c.request("warning", uuid.Nil, 0, format, v...)

	return c.client.Message(ctx, &request)
}

func (c *LogClient) Error(ctx context.Context, err error, format string, v ...interface{}) (*pb.Empty, error) {
	_format := format
	if err != nil {
		_format = fmt.Sprintf("%s: %s", err.Error(), _format)
	}
	request := c.request("error", uuid.Nil, 0, _format, v...)

	return c.client.Message(ctx, &request)
}

func (c *LogClient) Info(ctx context.Context, format string, v ...interface{}) (*pb.Empty, error) {
	request := c.request("info", uuid.Nil, 0, format, v...)

	return c.client.Message(ctx, &request)
}

func (c *LogClient) DBInfo(ctx context.Context, id uuid.UUID, format string, v ...interface{}) (*pb.Empty, error) {
	request := c.request("db_info", id, 0, format, v...)

	return c.client.Message(ctx, &request)
}

// this is only called from inside MultiInstruction's methods,
// that's why the request's caller_inc is set to 1
func (c *LogClient) TXInfo(ctx context.Context, id uuid.UUID, format string, v ...interface{}) (*pb.Empty, error) {
	request := c.request("tx_info", id, 1, format, v...)

	return c.client.Message(ctx, &request)
}
