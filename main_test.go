package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
)

// MockClipboard implements the Clipboard interface for testing
type MockClipboard struct {
	content string
	readErr error
}

func (m *MockClipboard) ReadAll() (string, error) {
	return m.content, m.readErr
}

func (m *MockClipboard) WriteAll(content string) {
	m.content = content
}

type MockHelp struct {
	called bool
}

func (m *MockHelp) PrintHelp() {
	m.called = true
}

// TestWriteStdinToClipboard tests the writeStdinToClipboard function
func TestWriteStdinToClipboard(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty input", "", ""},
		{"Single line input", "Hello, World!\n", "Hello, World!\n"},
		{"Multi-line input", "Line 1\nLine 2\n", "Line 1\nLine 2\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect stdin
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			r, w, _ := os.Pipe()
			os.Stdin = r

			go func() {
				w.Write([]byte(tt.input))
				w.Close()
			}()

			mockCb := &MockClipboard{}
			writeStdinToClipboard(mockCb)

			if mockCb.content != tt.expected {
				t.Errorf("Expected clipboard content %q, but got %q", tt.expected, mockCb.content)
			}
		})
	}
}

// TestPrintClipboardContents tests the printClipboardContents function
func TestPrintClipboardContents(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		readErr  error
		expected string
	}{
		{"Success case", "Clipboard content", nil, "Clipboard content"},
		{"Empty clipboard", "", nil, ""},
		{"Read error", "", errors.New("read error"), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCb := &MockClipboard{content: tt.content, readErr: tt.readErr}

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			printClipboardContents(mockCb)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if output != tt.expected {
				t.Errorf("Expected output %q, but got %q", tt.expected, output)
			}
		})
	}
}

// TestReadPipedStdin tests the readPipedStdin function
func TestReadPipedStdin(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty input", "", ""},
		{"Single line input", "Hello, World!\n", "Hello, World!\n"},
		{"Multi-line input", "Line 1\nLine 2\n", "Line 1\nLine 2\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect stdin
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			r, w, _ := os.Pipe()
			os.Stdin = r

			go func() {
				w.Write([]byte(tt.input))
				w.Close()
			}()

			output := readPipedStdin()

			if output != tt.expected {
				t.Errorf("Expected output %q, but got %q", tt.expected, output)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		input          string
		initialContent string
		expected       string
		expectHelpCall bool
	}{
		{"No args, no input, empty clipboard", []string{"cmd"}, "", "", "", false},
		{"No args, input provided, empty clipboard", []string{"cmd"}, "New content\n", "", "New content\n", false},
		{"No args, no input, existing clipboard content", []string{"cmd"}, "", "Existing content", "Existing content", false},
		{"No args, input provided, existing clipboard content", []string{"cmd"}, "New content\n", "Existing content", "New content\n", false},
		{"With args, help should be called", []string{"cmd", "-h"}, "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original args and restore after test
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args = tt.args

			// Redirect stdin
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			r, w, _ := os.Pipe()
			os.Stdin = r

			go func() {
				w.Write([]byte(tt.input))
				w.Close()
			}()

			// Capture stdout
			oldStdout := os.Stdout
			stdoutR, stdoutW, _ := os.Pipe()
			os.Stdout = stdoutW

			mockCb := &MockClipboard{content: tt.initialContent}
			mockHelp := &MockHelp{}

			run(mockCb, mockHelp)

			stdoutW.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			io.Copy(&buf, stdoutR)
			output := buf.String()

			if !tt.expectHelpCall {
				if output != tt.expected {
					t.Errorf("Expected output %q, but got %q", tt.expected, output)
				}

				if mockCb.content != tt.expected {
					t.Errorf("Expected clipboard content %q, but got %q", tt.expected, mockCb.content)
				}
			}

			if mockHelp.called != tt.expectHelpCall {
				t.Errorf("Expected help to be called: %v, but was: %v", tt.expectHelpCall, mockHelp.called)
			}
		})
	}
}
