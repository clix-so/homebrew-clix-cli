package utils

import (
	"fmt"
	"time"
)

type Spinner struct {
	frames []rune
	delay  time.Duration
	stop   chan struct{}
}

func NewSpinner() *Spinner {
	return &Spinner{
		frames: []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'},
		delay:  100 * time.Millisecond,
		stop:   make(chan struct{}),
	}
}

// Start spinner in goroutine with a message
func (s *Spinner) Start(msg string) {
	go func() {
		i := 0
		for {
			select {
			case <-s.stop:
				return
			default:
				frame := s.frames[i%len(s.frames)]
				fmt.Printf("\r\033[90;1m%c %s\033[0m", frame, msg)
				time.Sleep(s.delay)
				i++
			}
		}
	}()
}

// Stop spinner and clear line
func (s *Spinner) Stop() {
	s.stop <- struct{}{}
	fmt.Printf("\n")
}
