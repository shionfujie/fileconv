package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestSpinner(t *testing.T) {
	count := 10
	delta := 100 * time.Millisecond

	var rb bytes.Buffer
	spinner := NewFspinner(&rb, delta)
	ticker := time.NewTicker(delta)

	expected := "\n"
	if rb.String() != expected {
		t.Errorf("got %q, expected immediate newline", rb.String())
	}
	for i := 0; i < count; i++ {
		<-ticker.C
		// Since the terminal is tty, we use ansi escaping to animate
		// a spinner.
		r := `-\|/`[i%4]
		expected = fmt.Sprintf("%s\033[F%c\n", expected, r)
		if rb.String() != expected {
			t.Errorf("got %q after %d ticks, expected %q", rb.String(), i, expected)
		}
	}
	spinner.Stop()
	if rb.String() != expected {
		t.Errorf("got %q, expected to be unchanged", rb.String())
	}
	// Now test that the spinner stopped.
	for i := 0; i < count; i++ {
		<-ticker.C
		if rb.String() != "" {
			t.Error("Spinner did not shut down")
			break
		}
	}
	ticker.Stop()
}
