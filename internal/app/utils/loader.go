package utils

import (
	"fmt"
	"sync/atomic"
	"time"
)

const (
	step = 100 * time.Millisecond
)

// Loader defines the operations a loader can perform
type Loader interface {
	Start()
	Stop()
}

// loader represents a simple loading animation
type loader struct {
	running int32
	done    chan struct{}
}

// NewLoader creates a new Loader instance
func NewLoader() Loader {
	return &loader{
		done: make(chan struct{}),
	}
}

// Start starts the loading animation
func (l *loader) Start() {
	if !atomic.CompareAndSwapInt32(&l.running, 0, 1) {
		return
	}

	chars := []string{".  ", ".. ", "...", " ..", "  .", "   "}
	idx := 0

	go func() {
		for {
			select {
			case <-l.done:
				return
			default:
				fmt.Printf("\rLoading %s ", chars[idx])
				idx = (idx + 1) % len(chars)
				time.Sleep(step)
			}
		}
	}()
}

// Stop stops the loading animation
func (l *loader) Stop() {
	if !atomic.CompareAndSwapInt32(&l.running, 1, 0) {
		return
	}

	close(l.done)
	l.clear()
}

// clear clears the loading animation from the console
func (l *loader) clear() {
	fmt.Printf("\r\033[K")
}
