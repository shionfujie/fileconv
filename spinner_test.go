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

	spinAt := func(i int) string {
		r := `-\|/`[i%4]
		return fmt.Sprintf("%c\n", r)
	}

	var rb bytes.Buffer
	spinner := NewFspinner(&rb, delta)
	ticker := time.NewTicker(delta)

	if rb.String() != "\n" {
		t.Errorf("got %q, expected immediate newline", rb.String())
	}
	for i := 0; i < count; i++ {
		<-ticker.C
		if rb.String() != spinAt(i) {
			t.Errorf("got %q after %d ticks, expected %q", rb.String(), i, spinAt(i))
		}
	}
	spinner.Stop()
	if rb.String() != "" {
		t.Errorf("got %q, expected empty", rb.String())
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
