package help

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestPrintHelp(t *testing.T) {
	h := THelp{}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Use a channel to signal when the goroutine is done
	done := make(chan bool)

	// Run PrintHelp in a goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Check if the panic was due to os.Exit(0)
				if exitCode, ok := r.(int); ok && exitCode == 0 {
					done <- true
					return
				}
				// If it was a different panic, re-panic
				panic(r)
			}
		}()

		h.PrintHelp()
		// If we reach here, it means os.Exit(0) wasn't called
		t.Error("Expected os.Exit(0) to be called")
		done <- true
	}()

	// Wait for the goroutine to finish
	<-done

	// Close the write end of the pipe and restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read the captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expectedHelp := strings.TrimSpace(HELP)
	if output != expectedHelp {
		t.Errorf("Expected help message:\n%s\n\nGot:\n%s", expectedHelp, output)
	}
}

func TestHelpInterface(t *testing.T) {
	// Verify that THelp implements the Help interface
	var _ Help = (*THelp)(nil)
}
