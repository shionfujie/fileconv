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
	done chan struct{}
}

func NewSpinner(d time.Duration) *Spinner {
	return NewFspinner(os.Stdout, d)
}

func NewFspinner(w io.Writer, d time.Duration) *Spinner {
	s := &Spinner{
		done: make(chan struct{}, 1),
	}
	go startSpinner(s, d, w)
	return s
}

func (s *Spinner) Stop() {
	s.done <- struct{}{}
}

func startSpinner(s *Spinner, d time.Duration, w io.Writer) {
	ticker := time.NewTicker(d)
	fmt.Fprintln(w)
loop:
	for {
		for _, r := range `-\|/` {
			select {
			case <-ticker.C:
				fmt.Fprintf(w, "\033[F%c\n", r)
			case <-s.done:
				break loop
			}
		}
	}
	ticker.Stop()
	fmt.Fprintf(w, "\033[F")
}
