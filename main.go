package main

import (
	cli "github.com/lashajini/mind-palace/cli/cmd"
	"github.com/lashajini/mind-palace/config"
	"github.com/lashajini/mind-palace/storage/vdatabase"
)

func main() {
	cli.Execute()

	cfg := config.NewConfig()
	vdatabase.InitVDB(cfg)

	// _ = database.InitDB(cfg)
	//
	// c := rpcclient.NewClient(cfg)
	// m := pb.Memory{File: cli.FILE}
	// c.Add(context.Background(), &m)
}
