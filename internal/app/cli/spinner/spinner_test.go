package spinner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// waitForState polls the spinner state until the expected condition is met or timeout occurs.
func waitForState(s *Spinner, checkFunc func() bool, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		s.mu.RLock()
		result := checkFunc()
		s.mu.RUnlock()
		if result {
			return true
		}
		time.Sleep(time.Millisecond)
	}
	return false
}

// waitForRunning waits for the spinner to reach the expected running state.
func waitForRunning(s *Spinner, expectedRunning bool) bool {
	return waitForState(s, func() bool {
		return s.running == expectedRunning
	}, 100*time.Millisecond)
}

// waitForActive waits for the spinner to reach the expected active state.
func waitForActive(s *Spinner, expectedActive bool) bool {
	return waitForState(s, func() bool {
		return s.active == expectedActive
	}, 100*time.Millisecond)
}

func Test_Module(t *testing.T) {
	assert.NotNil(t, Module)
}

func Test_New(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "Success with message",
			message: "Loading…",
		},
		{
			name:    "Success with empty message",
			message: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := New(tt.message)

			assert.NotNil(t, s)
			assert.Equal(t, tt.message, s.message)
			assert.False(t, s.active)
			assert.False(t, s.running)
			assert.NotNil(t, s.done)
			assert.Equal(t, 1, cap(s.done))
		})
	}
}

func Test_Start(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := New("test")
			s.Start()

			assert.True(t, waitForRunning(s, true), "spinner should be running")
			assert.True(t, waitForActive(s, true), "spinner should be active")

			s.Stop()
		})
	}
}

func Test_Start_Idempotent(t *testing.T) {
	tests := []struct {
		name  string
		calls int
	}{
		{
			name:  "Success with multiple start calls",
			calls: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New("test")

			for i := 0; i < tt.calls; i++ {
				s.Start()
			}

			assert.True(t, waitForRunning(s, true), "spinner should be running after multiple starts")

			s.Stop()

			assert.True(t, waitForRunning(s, false), "spinner should stop running")
			assert.True(t, waitForActive(s, false), "spinner should be inactive")
		})
	}
}

func Test_SetMessage(t *testing.T) {
	tests := []struct {
		name       string
		initial    string
		newMessage string
	}{
		{
			name:       "Success with new message",
			initial:    "Loading…",
			newMessage: "Processing…",
		},
		{
			name:       "Success with empty message",
			initial:    "Loading…",
			newMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := New(tt.initial)
			s.SetMessage(tt.newMessage)

			s.mu.RLock()
			assert.Equal(t, tt.newMessage, s.message)
			s.mu.RUnlock()
		})
	}
}

func Test_Stop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New("test")
			s.Start()

			assert.True(t, waitForRunning(s, true), "spinner should start running")

			s.Stop()

			assert.True(t, waitForActive(s, false), "spinner should be inactive after stop")
			assert.True(t, waitForRunning(s, false), "spinner should stop running")
		})
	}
}

func Test_Stop_NotStarted(t *testing.T) {
	s := New("test")
	s.Stop()

	s.mu.RLock()
	assert.False(t, s.active)
	assert.False(t, s.running)
	s.mu.RUnlock()
}

func Test_NewSpinner(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			model := NewSpinner()
			assert.NotNil(t, model)

			view := model.View()
			assert.NotEmpty(t, view)
		})
	}
}

func Test_BubblesSpinner_Update(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			model := NewSpinner()
			assert.NotNil(t, model)

			tickMsg := model.Tick()
			updatedModel, cmd := model.Update(tickMsg)

			assert.NotNil(t, updatedModel)
			assert.NotNil(t, cmd)
		})
	}
}

func Test_BubblesSpinner_Tick(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			model := NewSpinner()
			assert.NotNil(t, model)

			msg := model.Tick()
			assert.NotNil(t, msg)
		})
	}
}
