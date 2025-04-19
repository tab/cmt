package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"cmt/internal/config"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
	PanicLevel = "panic"
	TraceLevel = "trace"

	ConsoleFormat = "console"
	JSONFormat    = "json"

	TimeFormat = "02.01.2006 15:04:05"
)

// Logger interface for application logging
type Logger interface {
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
}

type Event interface {
	Msg(msg string)
	Msgf(format string, v ...interface{})
	Str(key, value string) Event
	Int(key string, value int) Event
	Dur(key string, value time.Duration) Event
	Err(err error) Event
}

// AppLogger represents a logger implementation using zerolog
type AppLogger struct {
	log zerolog.Logger
}

// NewLogger creates a new logger instance
func NewLogger(cfg *config.Config) Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339

	var level zerolog.Level
	var output io.Writer

	if cfg.Logging.Level == "" {
		cfg.Logging.Level = InfoLevel
	}

	if cfg.Logging.Format == "" {
		cfg.Logging.Format = ConsoleFormat
	}

	level = getLogLevel(cfg.Logging.Level)

	switch cfg.Logging.Format {
	case JSONFormat:
		output = os.Stdout
	case ConsoleFormat:
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: TimeFormat,
		}
	default:
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: TimeFormat,
		}
	}

	logger := zerolog.
		New(output).
		Level(level).
		With().
		Timestamp().
		Str("version", config.Version).
		Logger()

	return &AppLogger{log: logger}
}

func (l *AppLogger) Debug() *zerolog.Event {
	return l.log.Debug()
}

func (l *AppLogger) Info() *zerolog.Event {
	return l.log.Info()
}

func (l *AppLogger) Warn() *zerolog.Event {
	return l.log.Warn()
}

func (l *AppLogger) Error() *zerolog.Event {
	return l.log.Error()
}

// getLogLevel converts string level to zerolog.Level
func getLogLevel(level string) zerolog.Level {
	switch level {
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	case FatalLevel:
		return zerolog.FatalLevel
	case PanicLevel:
		return zerolog.PanicLevel
	case TraceLevel:
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
}
