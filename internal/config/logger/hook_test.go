package logger

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_BufferHook_Run(t *testing.T) {
	buffer := NewLogBuffer()
	hook := NewBufferHook(buffer)

	hook.Run(nil, zerolog.InfoLevel, "test message")

	entries := buffer.GetEntries()
	assert.Len(t, entries, 1)
	assert.Equal(t, "test message", entries[0].Message)
	assert.Equal(t, "info", entries[0].Level)
}

func Test_BufferHook_NilBuffer(t *testing.T) {
	hook := NewBufferHook(nil)

	assert.NotPanics(t, func() {
		hook.Run(nil, zerolog.InfoLevel, "test message")
	})
}
