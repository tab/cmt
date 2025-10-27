package spinner

import (
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	Frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	FPS    = 80 * time.Millisecond
)

// Model represents a spinner model interface for Bubble Tea
type Model interface {
	Update(tea.Msg) (Model, tea.Cmd)
	View() string
	Tick() tea.Msg
}

// bubblesSpinner wraps the bubbles spinner to implement our Model interface
type bubblesSpinner struct {
	spinner.Model
}

// Update updates the spinner state
func (s bubblesSpinner) Update(msg tea.Msg) (Model, tea.Cmd) {
	updated, cmd := s.Model.Update(msg)
	return bubblesSpinner{updated}, cmd
}

// Tick returns the tick message
func (s bubblesSpinner) Tick() tea.Msg {
	return s.Model.Tick()
}

// NewSpinner creates a new Bubble Tea spinner model with app-wide configuration
func NewSpinner() Model {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: Frames,
		FPS:    FPS,
	}
	return bubblesSpinner{s}
}

// Spinner represents a simple text spinner for CLI operations
type Spinner struct {
	mu      sync.RWMutex
	message string
	active  bool
	running bool
	done    chan bool
	wg      sync.WaitGroup
}

// New creates a new CLI spinner with the given message
func New(message string) *Spinner {
	return &Spinner{
		message: message,
		active:  false,
		running: false,
		done:    make(chan bool, 1),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.active = true
	s.wg.Add(1)
	s.mu.Unlock()

	go func() {
		defer s.wg.Done()
		i := 0
		for {
			select {
			case <-s.done:
				return
			default:
				s.mu.RLock()
				msg := s.message
				s.mu.RUnlock()
				fmt.Printf("\r%s %s", Frames[i%len(Frames)], msg)
				i++
				time.Sleep(FPS)
			}
		}
	}()
}

// SetMessage updates the spinner message
func (s *Spinner) SetMessage(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
}

// Stop stops the spinner and clears the line
func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	wasRunning := s.running
	if wasRunning {
		s.running = false
	}
	s.mu.Unlock()

	if wasRunning {
		select {
		case s.done <- true:
		default:
		}
		s.wg.Wait()
	}

	fmt.Print("\r\033[K")
}
