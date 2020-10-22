package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

type result struct {
	file string
	err  error
}

func main() {
	fileNames := os.Args[1:]
	c := make(chan result, len(fileNames))

	d := 100 * time.Millisecond
	spinner := NewSpinner(d)
	for _, file := range fileNames {
		go func(file string) {
			var r result
			r.file, r.err = file, convert(file)
			c <- r
		}(file)
	}

	for range fileNames {
		r := <-c
		// Stop the spinner temporarily
		spinner.Stop()
		if r.err != nil {
			fmt.Printf("convert %s: failed to convert DJVU file to PDF file: %s\n", r.file, r.err.Error())
		} else {
			fmt.Printf("convert %s: DJVU file converted successfully to PDF file\n", r.file)
		}
		// Restart a spinner
		spinner = NewSpinner(d)
	}
	spinner.Stop()
}

func convert(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return err
	}
	if _, err := exec.Command("djvu2pdf", file).Output(); err != nil {
		return err
	}
	return nil
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
