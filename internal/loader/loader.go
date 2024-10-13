package loader

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Loader struct {
	running int32
}

func New() *Loader {
	return &Loader{}
}

func (l *Loader) Start() {
	atomic.StoreInt32(&l.running, 1)

	chars := []rune{'|', '/', '-', '\\'}
	idx := 0

	go func() {
		for atomic.LoadInt32(&l.running) == 1 {
			fmt.Printf("\r%c Loading...", chars[idx])
			idx = (idx + 1) % len(chars)
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (l *Loader) Stop() {
	atomic.StoreInt32(&l.running, 0)
	fmt.Printf("\r\033[K")
}
