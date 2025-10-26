package logger

import (
	"sync"
	"time"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Level          string
	Message        string
	FormattedLine  string
	Timestamp      time.Time
}

// LogBuffer implements a ring buffer for storing logs in memory
type LogBuffer struct {
	entries []LogEntry
	maxSize int
	mu      sync.RWMutex
}

// NewLogBuffer creates a new log buffer with a fixed size
func NewLogBuffer(maxSize int) *LogBuffer {
	if maxSize <= 0 {
		maxSize = 1000
	}
	return &LogBuffer{
		entries: make([]LogEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Add appends a log entry to the ring buffer
func (b *LogBuffer) Add(level, message string) {
	b.AddFormatted(level, message, "")
}

// AddFormatted appends a log entry with formatted line to the ring buffer
func (b *LogBuffer) AddFormatted(level, message, formattedLine string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	entry := LogEntry{
		Level:         level,
		Message:       message,
		FormattedLine: formattedLine,
		Timestamp:     time.Now(),
	}

	if len(b.entries) >= b.maxSize {
		b.entries = b.entries[1:]
	}

	b.entries = append(b.entries, entry)
}

// Entries returns a copy of all log entries
func (b *LogBuffer) Entries() []LogEntry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make([]LogEntry, len(b.entries))
	copy(result, b.entries)
	return result
}

// Clear removes all log entries
func (b *LogBuffer) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.entries = b.entries[:0]
}
