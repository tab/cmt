package logger

import (
	"fmt"
	"sync"
	"time"
)

const maxLogEntries = 1000

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
	Fields    map[string]interface{}
}

// LogBuffer is a thread-safe ring buffer for storing log entries
type LogBuffer struct {
	mu      sync.RWMutex
	entries []LogEntry
	maxSize int
	index   int
}

// NewLogBuffer creates a new log buffer
func NewLogBuffer() *LogBuffer {
	return &LogBuffer{
		entries: make([]LogEntry, 0, maxLogEntries),
		maxSize: maxLogEntries,
		index:   0,
	}
}

// Add adds a log entry to the buffer
func (b *LogBuffer) Add(entry LogEntry) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.entries) < b.maxSize {
		b.entries = append(b.entries, entry)
	} else {
		b.entries[b.index] = entry
		b.index = (b.index + 1) % b.maxSize
	}
}

// GetEntries returns a copy of all log entries in chronological order
func (b *LogBuffer) GetEntries() []LogEntry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if len(b.entries) < b.maxSize {
		result := make([]LogEntry, len(b.entries))
		copy(result, b.entries)
		return result
	}

	result := make([]LogEntry, b.maxSize)
	copy(result, b.entries[b.index:])
	copy(result[b.maxSize-b.index:], b.entries[:b.index])
	return result
}

// Clear clears all log entries from the buffer
func (b *LogBuffer) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.entries = make([]LogEntry, 0, b.maxSize)
	b.index = 0
}

// Format formats a log entry for display
func (e *LogEntry) Format() string {
	timestamp := e.Timestamp.Format("15:04:05")
	levelLabel := fmt.Sprintf("%-5s", e.Level)

	if len(e.Fields) == 0 {
		return fmt.Sprintf("%s %s %s", timestamp, levelLabel, e.Message)
	}

	fieldsStr := ""
	for k, v := range e.Fields {
		fieldsStr += fmt.Sprintf(" %s=%v", k, v)
	}

	return fmt.Sprintf("%s %s %s%s", timestamp, levelLabel, e.Message, fieldsStr)
}
