package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	cli "github.com/lashajini/mind-palace/cli/cmd"
	"github.com/lashajini/mind-palace/pkg/mperrors"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	go func() {
		<-ctx.Done()
		fmt.Println("received graceful shutdown signal")
	}()

	if err := cli.Execute(ctx); err != nil {
		mperrors.On(err).Exit()
	}
}
