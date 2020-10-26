package ioutil

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Spinner struct {
	w    io.Writer
	d    time.Duration
	done chan struct{}
}

func NewSpinner(d time.Duration) *Spinner {
	return NewFspinner(os.Stdout, d)
}

func NewFspinner(w io.Writer, d time.Duration) *Spinner {
	s := &Spinner{
		w:    w,
		d:    d,
		done: make(chan struct{}, 1),
	}
	startSpinner(s)
	return s
}

func (s *Spinner) Stop() {
	fmt.Fprintf(s.w, "\033[F")
	s.done <- struct{}{}
}

func (s *Spinner) Printf(format string, a ...interface{}) {
	// Stop the ticker temporarily
	s.Stop()
	fmt.Fprintf(s.w, format, a...)
	// Restart a ticker
	startSpinner(s)
}

func startSpinner(s *Spinner) {
	t := time.NewTicker(s.d)
	fmt.Fprintln(s.w)
	go func() {
	loop:
		for i := 0; ; i = (i + 1) % 4 {
			select {
			case <-t.C:
				fmt.Fprintf(s.w, "\033[F%c\n", `-\|/`[i])
			case <-s.done:
				t.Stop()
				break loop
			}
		}
	}()
}
