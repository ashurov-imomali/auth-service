package pkg

import (
	"github.com/rs/zerolog"
	"io"
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
	file, err := os.OpenFile("./logs/auth-service.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	multiWriter := io.MultiWriter(os.Stdout, file)
	logger := zerolog.New(multiWriter).With().Timestamp().Caller().Logger()
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
