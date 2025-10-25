package logger

import (
	"time"

	"github.com/rs/zerolog"
)

// BufferHook implements zerolog.Hook to capture log entries into a buffer
type BufferHook struct {
	buffer *LogBuffer
}

// NewBufferHook creates a new buffer hook
func NewBufferHook(buffer *LogBuffer) *BufferHook {
	return &BufferHook{
		buffer: buffer,
	}
}

// Run is called for every log event
func (h *BufferHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if h.buffer == nil {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level.String(),
		Message:   msg,
		Fields:    make(map[string]interface{}),
	}

	h.buffer.Add(entry)
}
