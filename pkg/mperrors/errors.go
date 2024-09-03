package mperrors

import (
	"context"
	"fmt"
	"os"

	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
)

type MPError struct {
	context string
	error
}

func On(err error) *MPError { return &MPError{context: "", error: err} }
func Onf(format string, a ...any) *MPError {
	return &MPError{context: "", error: fmt.Errorf(format, a...)}
}

func (e *MPError) Error() string {
	if len(e.context) != 0 {
		return fmt.Sprintf("%s: %v", e.context, e.error)
	}

	return e.error.Error()
}

func (e *MPError) Wrap(info string) *MPError {
	e.context = info
	return e
}

func (e *MPError) Wrapf(format string, a ...any) *MPError {
	e.context = fmt.Sprintf(format, a...)
	return e
}

func (e *MPError) Warn() {
	if e.error != nil {
		loggers.Log.Caller(1).Warn(context.Background(), e.Error())
	}
}

func (e *MPError) ExitWithMsgf(format string, a ...any) {
	if e.error != nil {
		loggers.Log.Caller(1).Error(context.Background(), e, format, a...)
		os.Exit(1)
	}
}

func ExitWithMsgf(format string, a ...any) {
	loggers.Log.Caller(1).Error(context.Background(), nil, format, a...)
	os.Exit(1)
}

func (e *MPError) ExitWithMsg(msg string) {
	if e.error != nil {
		loggers.Log.Caller(1).Error(context.Background(), e, msg)
		os.Exit(1)
	}
}

func ExitWithMsg(msg string) {
	loggers.Log.Caller(1).Error(context.Background(), nil, msg)
	os.Exit(1)
}

func (e *MPError) Exit() {
	if e.error != nil {
		loggers.Log.Caller(1).Error(context.Background(), e, "")
		os.Exit(1)
	}
}

func (e *MPError) PanicWithMsg(msg string) {
	if e.error != nil {
		panic(fmt.Errorf("%s %w", msg, e))
	}
}

func (e *MPError) Panic() {
	if e.error != nil {
		panic(e)
	}
}
