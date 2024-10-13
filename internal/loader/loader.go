package loader

import (
	"fmt"
	"time"
)

type Loader struct {
	running bool
}

func New() *Loader {
	return &Loader{}
}

func (l *Loader) Start() {
	l.running = true

	chars := []rune{'|', '/', '-', '\\'}
	idx := 0

	go func() {
		for l.running {
			fmt.Printf("\r%c Loading...", chars[idx])
			idx = (idx + 1) % len(chars)
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (l *Loader) Stop() {
	l.running = false
	fmt.Printf("\r\033[K")
}
