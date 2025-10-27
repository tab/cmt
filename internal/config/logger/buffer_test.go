package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewLogBuffer(t *testing.T) {
	tests := []struct {
		name            string
		maxSize         int
		expectedMaxSize int
	}{
		{
			name:            "Success with specified size",
			maxSize:         100,
			expectedMaxSize: 100,
		},
		{
			name:            "Success with zero size",
			maxSize:         0,
			expectedMaxSize: 1000,
		},
		{
			name:            "Success with negative size",
			maxSize:         -5,
			expectedMaxSize: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := NewLogBuffer(tt.maxSize)

			assert.NotNil(t, buffer)
			assert.Equal(t, tt.expectedMaxSize, buffer.maxSize)
			assert.NotNil(t, buffer.entries)
			assert.Equal(t, 0, len(buffer.entries))
		})
	}
}

func Test_LogBuffer_Add(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		message string
	}{
		{
			name:    "Success with info log entry",
			level:   "info",
			message: "test message",
		},
		{
			name:    "Success with error log entry",
			level:   "error",
			message: "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := NewLogBuffer(10)
			buffer.Add(tt.level, tt.message)

			entries := buffer.Entries()
			assert.Equal(t, 1, len(entries))
			assert.Equal(t, tt.level, entries[0].Level)
			assert.Equal(t, tt.message, entries[0].Message)
			assert.Equal(t, "", entries[0].FormattedLine)
		})
	}
}

func Test_LogBuffer_AddFormatted(t *testing.T) {
	tests := []struct {
		name          string
		level         string
		message       string
		formattedLine string
	}{
		{
			name:          "Success with formatted log entry",
			level:         "debug",
			message:       "debug message",
			formattedLine: "[DEBUG] debug message",
		},
		{
			name:          "Success with empty formatted line",
			level:         "warn",
			message:       "warning message",
			formattedLine: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := NewLogBuffer(10)
			buffer.AddFormatted(tt.level, tt.message, tt.formattedLine)

			entries := buffer.Entries()
			assert.Equal(t, 1, len(entries))
			assert.Equal(t, tt.level, entries[0].Level)
			assert.Equal(t, tt.message, entries[0].Message)
			assert.Equal(t, tt.formattedLine, entries[0].FormattedLine)
			assert.False(t, entries[0].Timestamp.IsZero())
		})
	}
}

func Test_LogBuffer_RingBuffer(t *testing.T) {
	tests := []struct {
		name      string
		maxSize   int
		addCount  int
		expectLen int
	}{
		{
			name:      "Success when exceeding max size",
			maxSize:   3,
			addCount:  5,
			expectLen: 3,
		},
		{
			name:      "Success when keeping latest entries",
			maxSize:   2,
			addCount:  4,
			expectLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := NewLogBuffer(tt.maxSize)

			for i := 0; i < tt.addCount; i++ {
				buffer.Add("info", "message")
			}

			entries := buffer.Entries()
			assert.Equal(t, tt.expectLen, len(entries))
		})
	}
}

func Test_LogBuffer_Entries(t *testing.T) {
	tests := []struct {
		name     string
		addCount int
	}{
		{
			name:     "Success with empty buffer",
			addCount: 0,
		},
		{
			name:     "Success with entries copy",
			addCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := NewLogBuffer(10)

			for i := 0; i < tt.addCount; i++ {
				buffer.Add("info", "message")
			}

			entries := buffer.Entries()
			assert.Equal(t, tt.addCount, len(entries))

			entries1 := buffer.Entries()
			assert.Equal(t, len(entries), len(entries1))

			if len(entries) > 0 {
				entries[0].Message = "modified"
				entries2 := buffer.Entries()
				assert.NotEqual(t, "modified", entries2[0].Message)
			}
		})
	}
}

func Test_LogBuffer_Clear(t *testing.T) {
	tests := []struct {
		name     string
		addCount int
	}{
		{
			name:     "Success when clearing with entries",
			addCount: 5,
		},
		{
			name:     "Success when clearing empty buffer",
			addCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := NewLogBuffer(10)

			for i := 0; i < tt.addCount; i++ {
				buffer.Add("info", "message")
			}

			assert.Equal(t, tt.addCount, len(buffer.Entries()))

			buffer.Clear()

			entries := buffer.Entries()
			assert.Equal(t, 0, len(entries))
		})
	}
}
