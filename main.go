package main

import (
	"context"

	cli "github.com/lashajini/mind-palace/cli/cmd"
	"github.com/lashajini/mind-palace/config"
	rpcclient "github.com/lashajini/mind-palace/rpc/client"
	pb "github.com/lashajini/mind-palace/rpc/client/gen/proto"
)

func main() {
	config := config.NewConfig()

	cli.Execute()
	c := rpcclient.NewClient(config.GRPC_SERVER_PORT)
	m := pb.Memory{File: cli.FILE}
	c.Add(context.Background(), &m)
}
