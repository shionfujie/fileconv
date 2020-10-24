package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func main() {
	fileNames := os.Args[1:]
	c := make(chan string, len(fileNames))

	// Start the first spinner
	d := 100 * time.Millisecond
	spinner := NewSpinner(d)

	for _, file := range fileNames {
		go func(file string) {
			err := convert(file)
			c <- formatMassage(file, err)
		}(file)
	}

	length := len(fileNames)
	width := length/10 + 1
	for i := 1; i <= length; i++ {
		s := <-c
		// Stop the spinner temporarily
		spinner.Stop()
		fmt.Printf("[%*d/%d] %s\n", width, i, length, s)
		// Restart a spinner
		spinner = NewSpinner(d)
	}
	spinner.Stop()
}

func formatMassage(file string, err error) string {
	if err != nil {
		return fmt.Sprintf("Fail\tconvert %s: %v", file, err)
	}
	return fmt.Sprintf("ok  \tconvert %s: DJVU file converted successfully to PDF file", file)
}

func convert(file string) error {
	if _, err := os.Stat(file); err != nil {
		return err.(*os.PathError).Err
	}
	if _, err := exec.Command("djvu2pdf", file).Output(); err != nil {
		return fmt.Errorf("failed to convert DJVU file to PDF file: %v", err)
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
