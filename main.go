package main

import (
	"context"

	cli "github.com/lashajini/mind-palace/cli/cmd"
	"github.com/lashajini/mind-palace/config"
	rpcclient "github.com/lashajini/mind-palace/rpc/client"
	pb "github.com/lashajini/mind-palace/rpc/client/gen/proto"
	"github.com/lashajini/mind-palace/storage/database"
)

func main() {
	cfg := config.NewConfig()
	database.InitDB(cfg)

	cli.Execute()
	c := rpcclient.NewClient(cfg)
	m := pb.Memory{File: cli.FILE}
	c.Add(context.Background(), &m)
}
