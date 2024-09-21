package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"tclip/help"

	"github.com/atotto/clipboard"
)

type Clipboard interface {
	WriteAll(string)
	ReadAll() (string, error)
}

type TClipboard struct{}

func (TClipboard) ReadAll() (string, error) {
	return clipboard.ReadAll()
}

func (TClipboard) WriteAll(content string) {
	clipboard.WriteAll(content)
}

func run(cb Clipboard, h help.Help) {
	if len(os.Args) > 1 {
		h.PrintHelp()
	}

	writeStdinToClipboard(cb)
	printClipboardContents(cb)
}

func writeStdinToClipboard(cb Clipboard) {
	input := readPipedStdin()
	if input != "" {
		cb.WriteAll(input)
	}
}

func printClipboardContents(cb Clipboard) {
	if content, err := cb.ReadAll(); err == nil {
		fmt.Print(content)
	}
}

func readPipedStdin() string {
	// Check if stdin is being piped
	info, err := os.Stdin.Stat()
	if err != nil {
		return ""
	}

	// Check if the file is a character device (terminal)
	if (info.Mode() & os.ModeCharDevice) != 0 {
		return ""
	}

	// Read from stdin
	reader := bufio.NewReader(os.Stdin)
	var output strings.Builder

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				output.WriteString(input)
				break
			}
			return ""
		}
		output.WriteString(input)
	}

	return output.String()
}

func main() {
	run(TClipboard{}, help.THelp{})
}
