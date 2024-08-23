package errors

import (
	"context"
	"fmt"
	"os"

	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
)

type _err struct {
	error
}

func On(err error) *_err { return &_err{err} }

func (e *_err) Warn() {
	if e.error != nil {
		ctx := context.Background()
		loggers.Log.Warn(ctx, e.Error())
	}
}

func (e *_err) ExitWithMsgf(format string, a ...any) {
	if e.error != nil {
		ctx := context.Background()
		loggers.Log.Error(ctx, e.error, format, a...)
		os.Exit(1)
	}
}

func ExitWithMsgf(format string, a ...any) {
	ctx := context.Background()
	loggers.Log.Error(ctx, nil, format, a...)
	os.Exit(1)
}

func (e *_err) ExitWithMsg(msg string) {
	if e.error != nil {
		ctx := context.Background()
		loggers.Log.Error(ctx, e.error, msg)
		os.Exit(1)
	}
}

func ExitWithMsg(msg string) {
	ctx := context.Background()
	loggers.Log.Error(ctx, nil, msg)
	os.Exit(1)
}

func (e *_err) Exit() {
	if e.error != nil {
		ctx := context.Background()
		loggers.Log.Error(ctx, e.error, "")
		os.Exit(1)
	}
}

func (e *_err) PanicWithMsg(msg string) {
	if e.error != nil {
		panic(fmt.Errorf("%s %w", msg, e.error))
	}
}

func (e *_err) Panic() {
	if e.error != nil {
		panic(e.error)
	}
}
