package main

import (
	"context"

	cli "github.com/lashajini/mind-palace/cli/cmd"
	rpcclient "github.com/lashajini/mind-palace/rpc/client"
	pb "github.com/lashajini/mind-palace/rpc/client/gen/proto"
)

func main() {
	cli.Execute()

	c := rpcclient.NewClient(50052)
	m := pb.Memory{File: cli.FILE}
	c.Add(context.Background(), &m)
}
