package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type Logger struct {
	zerolog.Logger
}

func New(level string, pretty bool) (*Logger, error) {
	var writer io.Writer = os.Stdout
	
	if pretty {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	baseLogger := zerolog.New(writer).
		Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{baseLogger}, nil
}

func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{l.With().Str("component", component).Logger()}
}

func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{l.With().Str("request_id", requestID).Logger()}
}

var globalLogger *Logger

func InitGlobalLogger(level string, pretty bool) error {
	logger, err := New(level, pretty)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

func Global() *Logger {
	if globalLogger == nil {
		logger, _ := New("info", true)
		globalLogger = logger
	}
	return globalLogger
}
