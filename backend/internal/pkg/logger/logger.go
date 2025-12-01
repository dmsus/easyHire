package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type Logger struct {
	*zerolog.Logger
}

func New(level string, pretty bool) (*Logger, error) {
	var writer io.Writer = os.Stdout
	
	if pretty {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			FormatLevel: func(i interface{}) string {
				var l string
				if ll, ok := i.(string); ok {
					switch ll {
					case "trace":
						l = colorize("TRC", 36) // Cyan
					case "debug":
						l = colorize("DBG", 32) // Green
					case "info":
						l = colorize("INF", 34) // Blue
					case "warn":
						l = colorize("WRN", 33) // Yellow
					case "error":
						l = colorize("ERR", 31) // Red
					case "fatal":
						l = colorize("FTL", 35) // Magenta
					case "panic":
						l = colorize("PNC", 35) // Magenta
					default:
						l = colorize("???", 37) // White
					}
				} else {
					l = colorize("???", 37)
				}
				return l
			},
			FormatMessage: func(i interface{}) string {
				if msg, ok := i.(string); ok {
					return msg
				}
				return ""
			},
		}
	}

	// Parse log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	logger := zerolog.New(writer).
		Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{&logger}, nil
}

func (l *Logger) WithComponent(component string) *Logger {
	logger := l.With().Str("component", component).Logger()
	return &Logger{&logger}
}

func (l *Logger) WithRequestID(requestID string) *Logger {
	logger := l.With().Str("request_id", requestID).Logger()
	return &Logger{&logger}
}

func (l *Logger) WithUserID(userID string) *Logger {
	logger := l.With().Str("user_id", userID).Logger()
	return &Logger{&logger}
}

func (l *Logger) WithAssessmentID(assessmentID string) *Logger {
	logger := l.With().Str("assessment_id", assessmentID).Logger()
	return &Logger{&logger}
}

func (l *Logger) WithCandidateID(candidateID string) *Logger {
	logger := l.With().Str("candidate_id", candidateID).Logger()
	return &Logger{&logger}
}

// Helper functions for different log levels with context
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.Debug().Fields(fieldsToMap(fields...)).Msg(msg)
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.Info().Fields(fieldsToMap(fields...)).Msg(msg)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.Warn().Fields(fieldsToMap(fields...)).Msg(msg)
}

func (l *Logger) Error(msg string, err error, fields ...interface{}) {
	if err != nil {
		allFields := append(fields, "error", err.Error())
		l.Error().Fields(fieldsToMap(allFields...)).Msg(msg)
	} else {
		l.Error().Fields(fieldsToMap(fields...)).Msg(msg)
	}
}

func (l *Logger) Fatal(msg string, err error, fields ...interface{}) {
	if err != nil {
		allFields := append(fields, "error", err.Error())
		l.Fatal().Fields(fieldsToMap(allFields...)).Msg(msg)
	} else {
		l.Fatal().Fields(fieldsToMap(fields...)).Msg(msg)
	}
}

func (l *Logger) Panic(msg string, err error, fields ...interface{}) {
	if err != nil {
		allFields := append(fields, "error", err.Error())
		l.Panic().Fields(fieldsToMap(allFields...)).Msg(msg)
	} else {
		l.Panic().Fields(fieldsToMap(fields...)).Msg(msg)
	}
}

func fieldsToMap(fields ...interface{}) map[string]interface{} {
	if len(fields)%2 != 0 {
		// If odd number of fields, ignore the last one
		fields = fields[:len(fields)-1]
	}

	result := make(map[string]interface{})
	for i := 0; i < len(fields); i += 2 {
		if key, ok := fields[i].(string); ok {
			result[key] = fields[i+1]
		}
	}
	return result
}

func colorize(s interface{}, c int) string {
	return "\x1b[" + string(rune(c)) + "m" + s.(string) + "\x1b[0m"
}

// Global logger instance
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
		// Initialize with defaults if not initialized
		logger, _ := New("info", true)
		globalLogger = logger
	}
	return globalLogger
}
