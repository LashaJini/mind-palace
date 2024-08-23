package loggers

import (
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/rpc/log"
)

var cfg = common.NewConfig()
var Log = log.NewLogClient(cfg, "CLI")
