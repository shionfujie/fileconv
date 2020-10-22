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
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Fatal(err)
	}
	spinner := NewSpinner(100 * time.Millisecond)
	if _, err := exec.Command("djvu2pdf", file).Output(); err != nil {
		log.Fatal(err)
	}
	spinner.Stop()
	fmt.Printf("convert %s: DJVU file converted successfully to PDF file\n", file)
}

type Spinner struct {
	w    io.Writer
	done chan struct{}
}

func NewSpinner(d time.Duration) *Spinner {
	return NewFspinner(os.Stdout, d)
}

func NewFspinner(w io.Writer, d time.Duration) *Spinner {
	s := &Spinner{
		w:    w,
		done: make(chan struct{}, 1),
	}
	startSpinner(s, d)
	return s
}

func (s *Spinner) Stop() {
	fmt.Fprintf(s.w, "\033[F")
	s.done <- struct{}{}
}

func startSpinner(s *Spinner, d time.Duration) {
	t := time.NewTicker(d)
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
