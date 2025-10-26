package spinner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Module(t *testing.T) {
	assert.NotNil(t, Module)
}

func Test_New(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "creates spinner with message",
			message: "Loading...",
		},
		{
			name:    "creates spinner with empty message",
			message: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			name: "starts spinner animation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New("test")
			s.Start()

			time.Sleep(10 * time.Millisecond)

			s.mu.RLock()
			assert.True(t, s.running)
			assert.True(t, s.active)
			s.mu.RUnlock()

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
			name:  "multiple start calls do not leak goroutines",
			calls: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New("test")

			for i := 0; i < tt.calls; i++ {
				s.Start()
				time.Sleep(5 * time.Millisecond)
			}

			s.mu.RLock()
			assert.True(t, s.running)
			s.mu.RUnlock()

			s.Stop()

			time.Sleep(10 * time.Millisecond)

			s.mu.RLock()
			assert.False(t, s.running)
			assert.False(t, s.active)
			s.mu.RUnlock()
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
			name:       "updates message",
			initial:    "Loading...",
			newMessage: "Processing...",
		},
		{
			name:       "sets empty message",
			initial:    "Loading...",
			newMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			name: "stops running spinner",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New("test")
			s.Start()

			time.Sleep(10 * time.Millisecond)

			s.Stop()

			time.Sleep(10 * time.Millisecond)

			s.mu.RLock()
			assert.False(t, s.active)
			assert.False(t, s.running)
			s.mu.RUnlock()
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
			name: "creates new bubble tea spinner model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewSpinner()
			assert.NotNil(t, model)

			view := model.View()
			assert.NotEmpty(t, view)
		})
	}
}

func Test_bubblesSpinner_Update(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "updates spinner state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewSpinner()
			assert.NotNil(t, model)

			tickMsg := model.Tick()
			updatedModel, cmd := model.Update(tickMsg)

			assert.NotNil(t, updatedModel)
			assert.NotNil(t, cmd)
		})
	}
}

func Test_bubblesSpinner_Tick(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "returns tick message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewSpinner()
			assert.NotNil(t, model)

			msg := model.Tick()
			assert.NotNil(t, msg)
		})
	}
}
