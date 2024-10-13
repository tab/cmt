package loader

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Start(t *testing.T) {
	tests := []struct {
		name     string
		waitTime time.Duration
		expect   bool
	}{
		{
			name:     "Running",
			waitTime: 50 * time.Millisecond,
			expect:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := New()
			loader.Start()
			defer loader.Stop()

			time.Sleep(tt.waitTime)
			assert.Equal(t, tt.expect, loader.running)
		})
	}
}

func Test_Stop(t *testing.T) {
	tests := []struct {
		name     string
		waitTime time.Duration
		expect   bool
	}{
		{
			name:     "Stopped",
			waitTime: 50 * time.Millisecond,
			expect:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := New()
			loader.Start()
			loader.Stop()

			time.Sleep(tt.waitTime)
			assert.Equal(t, tt.expect, loader.running)
		})
	}
}
