package utils

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewLoader(t *testing.T) {
	l := NewLoader()
	assert.NotNil(t, l)

	instance, ok := l.(*loader)
	assert.True(t, ok)
	assert.NotNil(t, instance.done)
	assert.Equal(t, int32(0), instance.running)
}

func Test_Start(t *testing.T) {
	tests := []struct {
		name          string
		startTimes    int
		sleepDuration time.Duration
		wantOutput    bool
	}{
		{
			name:          "No start",
			startTimes:    0,
			sleepDuration: 100 * time.Millisecond,
			wantOutput:    false,
		},
		{
			name:          "Single start",
			startTimes:    1,
			sleepDuration: 250 * time.Millisecond,
			wantOutput:    true,
		},
		{
			name:          "Double start",
			startTimes:    2,
			sleepDuration: 250 * time.Millisecond,
			wantOutput:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, err := os.Pipe()
			assert.NoError(t, err)

			os.Stdout = w
			defer func() {
				w.Close()
				r.Close()
				os.Stdout = old
			}()

			l := NewLoader()
			for i := 0; i < tt.startTimes; i++ {
				l.Start()
			}
			time.Sleep(tt.sleepDuration)
			l.Stop()

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			assert.NoError(t, err)
			output := buf.String()

			if tt.wantOutput {
				assert.Contains(t, output, "Loading")
			} else {
				assert.NotContains(t, output, "Loading")
			}
		})
	}
}

func Test_Stop(t *testing.T) {
	tests := []struct {
		name          string
		startFirst    bool
		stopTimes     int
		sleepDuration time.Duration
		wantClear     int
	}{
		{
			name:       "Stop without start",
			startFirst: false,
			stopTimes:  1,
			wantClear:  0,
		},
		{
			name:          "Start then stop",
			startFirst:    true,
			stopTimes:     1,
			sleepDuration: 100 * time.Millisecond,
			wantClear:     1,
		},
		{
			name:          "Start then double stop",
			startFirst:    true,
			stopTimes:     2,
			sleepDuration: 100 * time.Millisecond,
			wantClear:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, err := os.Pipe()
			assert.NoError(t, err)
			os.Stdout = w
			defer func() {
				w.Close()
				r.Close()
				os.Stdout = old
			}()

			l := NewLoader()
			if tt.startFirst {
				l.Start()
				time.Sleep(tt.sleepDuration)
			}
			for i := 0; i < tt.stopTimes; i++ {
				l.Stop()
			}

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			assert.NoError(t, err)
			output := buf.String()

			count := bytes.Count([]byte(output), []byte("\033[K"))
			assert.Equal(t, tt.wantClear, count)
		})
	}
}
