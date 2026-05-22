package highlight

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// ANSI color codes.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
)

// Highlighter applies ANSI color codes to log lines based on severity.
type Highlighter struct {
	enabled bool
}

// New returns a new Highlighter. When enabled is false, Format returns
// the raw line unchanged — useful when output is not a TTY.
func New(enabled bool) *Highlighter {
	return &Highlighter{enabled: enabled}
}

// Format returns a colorized representation of the log line.
// If the highlighter is disabled the original raw string is returned.
func (h *Highlighter) Format(raw string, sev parser.Severity) string {
	if !h.enabled {
		return raw
	}
	return fmt.Sprintf("%s%s%s", severityColor(sev), raw, colorReset)
}

// FormatSeverity returns a colorized severity label string.
func (h *Highlighter) FormatSeverity(sev parser.Severity) string {
	if !h.enabled {
		return sev.String()
	}
	return fmt.Sprintf("%s%s%s%s", colorBold, severityColor(sev), sev.String(), colorReset)
}

func severityColor(sev parser.Severity) string {
	switch strings.ToUpper(sev.String()) {
	case "ERROR", "FATAL", "CRITICAL":
		return colorRed
	case "WARN", "WARNING":
		return colorYellow
	case "INFO":
		return colorCyan
	default:
		return colorWhite
	}
}
