package common

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var once sync.Once
var log Logger

type Logger struct {
	zerolog.Logger
}

func (l *Logger) DBInfo(id uuid.UUID, msg string) {
	l.Info().Msgf("%s%s (tx=%s)%s", COLOR_YELLOW, msg, id.String(), COLOR_RESET)
}

func (l *Logger) TXInfo(id uuid.UUID, msg string) {
	l.Info().Msgf("%s%s --- (tx=%s)%s", COLOR_YELLOW, msg, id.String(), COLOR_RESET)
}

func NewLoggger() Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

		logLevel, err := strconv.Atoi(os.Getenv(LOG_LEVEL))
		if err != nil {
			logLevel = int(zerolog.InfoLevel) // default to INFO
		}

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.TimeOnly,
		}

		env := os.Getenv(MP_ENV)
		if env == PROD_ENV {
			fileLogger := lumberjackLogger(env)
			output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
		}

		log = Logger{
			zerolog.New(output).
				Level(zerolog.Level(logLevel)).
				Sample(zerolog.LevelSampler{
					DebugSampler: &zerolog.BasicSampler{N: 10},
				}).
				With().
				Timestamp().
				// CallerWithSkipFrameCount(3). // messes with layout
				Logger(),
		}
	})

	return log
}

func lumberjackLogger(env string) *lumberjack.Logger {
	fileLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("./logs/pp7.%s.log", env),
		MaxSize:    5,
		MaxBackups: 10,
		MaxAge:     7,
		Compress:   true,
	}

	return fileLogger
}
