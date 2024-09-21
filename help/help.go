package help

import (
	"fmt"
	"os"
	"strings"
)

const HELP = `
tclip - Terminal clipboard utility

tclip prints the contents of the user's clipboard.
Any input piped into tclip will be copied to the user's clipboard, then printed.

Usage: tclip

Examples:
  tclip > clipboard.txt
  echo pipedinput | tclip
`

type Help interface {
	PrintHelp()
}

type THelp struct{}

func (THelp) PrintHelp() {
	fmt.Print(strings.TrimSpace(HELP))
	os.Exit(0)
}
