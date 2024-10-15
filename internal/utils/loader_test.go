package utils

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoader_Start(t *testing.T) {
	tests := []struct {
		name          string
		waitTime      time.Duration
		expectRunning bool
	}{
		{
			name:          "Loader starts and runs",
			waitTime:      200 * time.Millisecond,
			expectRunning: true,
		},
		{
			name:          "Loader runs for a short time",
			waitTime:      100 * time.Millisecond,
			expectRunning: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoader()

			l.Start()
			time.Sleep(tt.waitTime)

			running := atomic.LoadInt32(&l.running) == 1
			assert.Equal(t, tt.expectRunning, running)

			l.Stop()
		})
	}
}

func TestLoader_Stop(t *testing.T) {
	tests := []struct {
		name          string
		waitTime      time.Duration
		expectRunning bool
	}{
		{
			name:          "Loader stops after start",
			waitTime:      100 * time.Millisecond,
			expectRunning: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoader()

			l.Start()
			l.Stop()
			time.Sleep(tt.waitTime)

			running := atomic.LoadInt32(&l.running) == 1
			assert.Equal(t, tt.expectRunning, running)
		})
	}
}
