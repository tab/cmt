package logger

import (
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

	TimeFormat = "02.01.2006 15:04:05"

	BufferSize = 1000
)

// Logger interface for application logging
type Logger interface {
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
	GetBuffer() *LogBuffer
}

// AppLogger represents a logger implementation using zerolog
type AppLogger struct {
	buffer *LogBuffer
	log    zerolog.Logger
}

// NewLogger creates a new logger instance with buffer (no console output)
func NewLogger(cfg *config.Config) Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339

	var level zerolog.Level

	if cfg.Logging.Level == "" {
		cfg.Logging.Level = InfoLevel
	}

	level = getLogLevel(cfg.Logging.Level)

	buffer := NewLogBuffer(BufferSize)

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: TimeFormat,
	}

	output := &bufferedWriter{
		buffer:        buffer,
		consoleWriter: consoleWriter,
	}

	logger := zerolog.
		New(output).
		Level(level).
		With().
		Timestamp().
		Str("version", config.Version).
		Logger()

	return &AppLogger{
		buffer: buffer,
		log:    logger,
	}
}

// GetBuffer returns the log buffer if available
func (l *AppLogger) GetBuffer() *LogBuffer {
	return l.buffer
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
