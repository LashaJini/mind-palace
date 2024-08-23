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

func (c *LogClient) request(msg, filename, logType string, id uuid.UUID, line int) pb.LogRequest {
	return pb.LogRequest{
		Msg:         msg,
		Filename:    filename,
		Line:        int32(line),
		ServiceName: c.serviceName,
		Type:        logType,
		Id:          id.String(),
	}
}

func (c *LogClient) Debug(ctx context.Context, format string, v ...interface{}) error {
	_, filename, line, _ := runtime.Caller(1)
	msg := fmt.Sprintf(format, v...)
	request := c.request(msg, filename, "debug", uuid.Nil, line)

	_, err := c.client.Message(ctx, &request)

	return err
}

func (c *LogClient) Warn(ctx context.Context, format string, v ...interface{}) error {
	_, filename, line, _ := runtime.Caller(1)
	msg := fmt.Sprintf(format, v...)
	request := c.request(msg, filename, "warning", uuid.Nil, line)

	_, err := c.client.Message(ctx, &request)

	return err
}

func (c *LogClient) Error(ctx context.Context, err error, format string, v ...interface{}) error {
	_, filename, line, _ := runtime.Caller(1)
	msg := fmt.Sprintf(format, v...)
	if err != nil {
		msg = fmt.Sprintf("%s: %s", err.Error(), msg)
	}
	request := c.request(msg, filename, "error", uuid.Nil, line)

	_, _err := c.client.Message(ctx, &request)

	return _err
}

func (c *LogClient) Info(ctx context.Context, format string, v ...interface{}) error {
	_, filename, line, _ := runtime.Caller(1)
	msg := fmt.Sprintf(format, v...)
	request := c.request(msg, filename, "info", uuid.Nil, line)

	_, err := c.client.Message(ctx, &request)

	return err
}

func (c *LogClient) DBInfo(ctx context.Context, id uuid.UUID, format string, v ...interface{}) error {
	_, filename, line, _ := runtime.Caller(1)
	msg := fmt.Sprintf(format, v...)
	request := c.request(msg, filename, "db_info", id, line)

	_, err := c.client.Message(ctx, &request)

	return err
}

// this should be called only from inside MultiInstruction's methods, because the caller stack is set to 2
func (c *LogClient) TXInfo(ctx context.Context, id uuid.UUID, format string, v ...interface{}) error {
	_, filename, line, _ := runtime.Caller(2)
	msg := fmt.Sprintf(format, v...)
	request := c.request(msg, filename, "tx_info", id, line)

	_, err := c.client.Message(ctx, &request)

	return err
}
