package highlight

import (
	"os"

	"golang.org/x/term"
)

// IsTTY reports whether the given file descriptor refers to a terminal.
// It is used by callers to decide whether to enable ANSI highlighting.
func IsTTY(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}
