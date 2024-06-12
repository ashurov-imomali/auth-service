package pkg

import (
	"github.com/rs/zerolog"
	"os"
)

type Logger struct {
	Log *zerolog.Logger
}
type Log interface {
	Info(string)
	Error(error, string)
	Warn(string)
	Debug(string)
}

func GetLogger() (Log, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().CallerWithSkipFrameCount(3).Logger()
	return &Logger{Log: &logger}, nil
}

func (l *Logger) Info(msg string) {
	l.Log.Info().Msg(msg)
}

func (l *Logger) Error(err error, msg string) {
	l.Log.Error().Err(err).Msg(msg)
}

func (l *Logger) Warn(msg string) {
	l.Log.Warn().Msg(msg)
}

func (l *Logger) Debug(msg string) {
	l.Log.Debug().Msg(msg)
}
