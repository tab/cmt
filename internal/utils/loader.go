package utils

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Loader represents a loader
type Loader struct {
	running int32
	done    chan struct{}
}

// NewLoader creates a new loader instance
func NewLoader() *Loader {
	return &Loader{
		done: make(chan struct{}),
	}
}

// Start starts the loader
func (l *Loader) Start() {
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
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the loader
func (l *Loader) Stop() {
	if !atomic.CompareAndSwapInt32(&l.running, 1, 0) {
		return
	}

	close(l.done)
	l.clear()
}

func (l *Loader) clear() {
	fmt.Printf("\r\033[K")
}
