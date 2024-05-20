package main

import (
	cli "github.com/lashajini/mind-palace/cli/cmd"
)

func main() {
	cli.Execute()

	// c := rpcclient.NewClient(cfg)
	// m := pb.Memory{File: cli.FILE}
	// c.Add(context.Background(), &m)
}
