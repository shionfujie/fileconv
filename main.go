package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	file := os.Args[1]
	spinner := NewSpinner(100 * time.Millisecond)
	_, err := exec.Command("djvu2pdf", file).Output()
	if err != nil {
		log.Fatal(err)
	}
	spinner.Stop()
	fmt.Printf("convert %s: DJVU file converted successfully to PDF file\n", file)
}

type Spinner struct {
	w    io.Writer
	t    *time.Ticker
	done chan struct{}
}

func NewSpinner(d time.Duration) *Spinner {
	return NewFspinner(os.Stdout, d)
}

func NewFspinner(w io.Writer, d time.Duration) *Spinner {
	s := &Spinner{
		w:    w,
		t:    time.NewTicker(d),
		done: make(chan struct{}, 1),
	}
	startSpinner(s)
	return s
}

func (s *Spinner) Stop() {
	s.t.Stop()
	fmt.Fprintf(s.w, "\033[F")
	s.done <- struct{}{}
}

func startSpinner(s *Spinner) {
	fmt.Fprintln(s.w)
	go func() {
	loop:
		for i := 0; ; i = (i + 1) % 4 {
			select {
			case <-s.t.C:
				fmt.Fprintf(s.w, "\033[F%c\n", `-\|/`[i])
			case <-s.done:
				break loop
			}
		}
	}()
}
