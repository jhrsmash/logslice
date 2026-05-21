package parser

import (
	"fmt"
	"time"
)

// Severity represents log level severity.
type Severity int

const (
	SeverityUnknown Severity = iota
	SeverityDebug
	SeverityInfo
	SeverityWarn
	SeverityError
	SeverityFatal
)

// severityLabels maps string labels to Severity values.
var severityLabels = map[string]Severity{
	"DEBUG": SeverityDebug,
	"INFO":  SeverityInfo,
	"WARN":  SeverityWarn,
	"WARNING": SeverityWarn,
	"ERROR": SeverityError,
	"FATAL": SeverityFatal,
}

// String returns the string representation of a Severity.
func (s Severity) String() string {
	switch s {
	case SeverityDebug:
		return "DEBUG"
	case SeverityInfo:
		return "INFO"
	case SeverityWarn:
		return "WARN"
	case SeverityError:
		return "ERROR"
	case SeverityFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ParseSeverity parses a severity string, returning SeverityUnknown if unrecognised.
func ParseSeverity(s string) Severity {
	if sev, ok := severityLabels[s]; ok {
		return sev
	}
	return SeverityUnknown
}

// LogLine represents a single parsed log entry.
type LogLine struct {
	Timestamp time.Time
	Severity  Severity
	Message   string
	Raw       string
}

// ErrUnparseable is returned when a line cannot be parsed.
type ErrUnparseable struct {
	Line string
}

func (e *ErrUnparseable) Error() string {
	return fmt.Sprintf("unparseable log line: %q", e.Line)
}
