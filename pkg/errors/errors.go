package errors

import (
	"fmt"
	"os"

	"github.com/lashajini/mind-palace/pkg/common"
)

type _err struct {
	error
}

func On(err error) *_err { return &_err{err} }

func (e *_err) Warn() {
	if e.error != nil {
		common.Log.Warn().Stack().Err(e.error).Send()
	}
}

func (e *_err) ExitWithMsgf(format string, a ...any) {
	if e.error != nil {
		common.Log.Error().Stack().Err(e.error).Msgf(format, a...)
		os.Exit(1)
	}
}

func ExitWithMsgf(format string, a ...any) {
	common.Log.Error().Stack().Err(nil).Msgf(format, a...)
	os.Exit(1)
}

func (e *_err) ExitWithMsg(msg string) {
	if e.error != nil {
		common.Log.Error().Stack().Err(e.error).Msg(msg)
		os.Exit(1)
	}
}

func ExitWithMsg(msg string) {
	common.Log.Error().Stack().Err(nil).Msg(msg)
	os.Exit(1)
}

func (e *_err) Exit() {
	if e.error != nil {
		common.Log.Error().Stack().Err(e.error).Send()
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
