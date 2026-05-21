package parser

import (
	"regexp"
	"strings"
	"time"
)

// Common log timestamp formats to attempt when parsing.
var timestampFormats = []string{
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
}

// logLineRegexp matches: <timestamp> <SEVERITY> <message>
// e.g. "2024-01-15T10:23:45Z INFO  server started"
var logLineRegexp = regexp.MustCompile(
	`^(\S+(?:\s\S+)?)\s+(DEBUG|INFO|WARN|WARNING|ERROR|FATAL)\s+(.*)$`,
)

// Parser parses raw log lines into LogLine structs.
type Parser struct{}

// New returns a new Parser.
func New() *Parser {
	return &Parser{}
}

// Parse attempts to parse a raw log line string into a LogLine.
// Returns ErrUnparseable if the line does not match the expected format.
func (p *Parser) Parse(raw string) (*LogLine, error) {
	raw = strings.TrimRight(raw, "\r\n")
	matches := logLineRegexp.FindStringSubmatch(raw)
	if matches == nil {
		return nil, &ErrUnparseable{Line: raw}
	}

	timestampStr := matches[1]
	severityStr := matches[2]
	message := matches[3]

	ts, err := parseTimestamp(timestampStr)
	if err != nil {
		return nil, &ErrUnparseable{Line: raw}
	}

	return &LogLine{
		Timestamp: ts,
		Severity:  ParseSeverity(severityStr),
		Message:   message,
		Raw:       raw,
	}, nil
}

// parseTimestamp tries each known format until one succeeds.
func parseTimestamp(s string) (time.Time, error) {
	for _, fmt := range timestampFormats {
		if t, err := time.Parse(fmt, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, &ErrUnparseable{Line: s}
}
