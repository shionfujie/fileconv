package ioutil_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	. "github.com/shionfujie/fileconv/io/ioutil"
)

func TestSpinner(t *testing.T) {
	count := 10
	delta := 100 * time.Millisecond
	var rb bytes.Buffer
	spinner := NewFspinner(&rb, delta)

	expected := "\n"
	if rb.String() != expected {
		t.Errorf("got %q, expected immediate newline", rb.String())
	}
	for i := 0; i < count; i++ {
		// Let the spinner pass first.
		time.Sleep(delta * 101 / 100)
		// Since the terminal is tty, we use ansi escaping to animate
		// a spinner.
		r := `-\|/`[i%4]
		expected = fmt.Sprintf("%s\033[F%c\n", expected, r)
		if rb.String() != expected {
			t.Errorf("after %d ticks: got %q, expected %q", i+1, rb.String(), expected)
		}
	}
	spinner.Stop()
	// Test tear down escaping.
	expected = fmt.Sprintf("%s\033[F", expected)
	if rb.String() != expected {
		t.Errorf("teardown: got %q, expected %q", rb.String(), expected)
	}
	// Now test that the spinner stopped.
	for i := 0; i < count; i++ {
		time.Sleep(delta)
		if rb.String() != expected {
			t.Error("Spinner did not shut down")
			break
		}
	}
}
