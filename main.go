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
	t *time.Ticker
	done chan struct{}
}

func NewSpinner(d time.Duration) *Spinner {
	return NewFspinner(os.Stdout, d)
}

func NewFspinner(w io.Writer, d time.Duration) *Spinner {
	s := &Spinner{
		t: time.NewTicker(d),
		done: make(chan struct{}, 1),
	}
	fmt.Fprintln(w)
	go startSpinner(s, w)
	return s
}

func (s *Spinner) Stop() {
	s.done <- struct{}{}
}

func startSpinner(s *Spinner, w io.Writer) {
	t := s.t
loop:
	for {
		for _, r := range `-\|/` {
			select {
			case <-t.C:
				fmt.Fprintf(w, "\033[F%c\n", r)
			case <-s.done:
				break loop
			}
		}
	}
	t.Stop()
	fmt.Fprintf(w, "\033[F")
}
