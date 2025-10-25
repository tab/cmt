package logger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_LogBuffer_Add(t *testing.T) {
	buffer := NewLogBuffer()
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "debug",
		Message:   "test message",
		Fields:    map[string]interface{}{"key": "value"},
	}

	buffer.Add(entry)
	entries := buffer.GetEntries()

	assert.Len(t, entries, 1)
	assert.Equal(t, "test message", entries[0].Message)
	assert.Equal(t, "debug", entries[0].Level)
}

func Test_LogBuffer_RingBuffer(t *testing.T) {
	buffer := NewLogBuffer()

	for i := 0; i < maxLogEntries+10; i++ {
		buffer.Add(LogEntry{
			Timestamp: time.Now(),
			Level:     "info",
			Message:   "test",
			Fields:    make(map[string]interface{}),
		})
	}

	entries := buffer.GetEntries()
	assert.Len(t, entries, maxLogEntries)
}

func Test_LogBuffer_Clear(t *testing.T) {
	buffer := NewLogBuffer()
	buffer.Add(LogEntry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "test",
		Fields:    make(map[string]interface{}),
	})

	entries := buffer.GetEntries()
	assert.Len(t, entries, 1)

	buffer.Clear()
	entries = buffer.GetEntries()
	assert.Len(t, entries, 0)
}

func Test_LogEntry_Format(t *testing.T) {
	tests := []struct {
		name     string
		entry    LogEntry
		contains []string
	}{
		{
			name: "Message only",
			entry: LogEntry{
				Timestamp: time.Date(2024, 1, 1, 12, 34, 56, 0, time.UTC),
				Level:     "INF",
				Message:   "test message",
				Fields:    make(map[string]interface{}),
			},
			contains: []string{"12:34:56", "INF", "test message"},
		},
		{
			name: "Message with fields",
			entry: LogEntry{
				Timestamp: time.Date(2024, 1, 1, 12, 34, 56, 0, time.UTC),
				Level:     "DBG",
				Message:   "test message",
				Fields:    map[string]interface{}{"key": "value", "count": 42},
			},
			contains: []string{"12:34:56", "DBG", "test message", "key=value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := tt.entry.Format()
			for _, substr := range tt.contains {
				assert.Contains(t, formatted, substr)
			}
		})
	}
}
