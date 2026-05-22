package filter

import (
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// Options holds the filtering criteria for log lines.
type Options struct {
	From     time.Time
	To       time.Time
	MinLevel parser.Severity
}

// Filter decides whether a parsed log line should be included in output.
type Filter struct {
	opts Options
}

// New creates a new Filter with the given options.
func New(opts Options) *Filter {
	return &Filter{opts: opts}
}

// Match returns true if the log line falls within the time range and
// meets the minimum severity level.
func (f *Filter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}

	if !f.opts.From.IsZero() && line.Timestamp.Before(f.opts.From) {
		return false
	}

	if !f.opts.To.IsZero() && line.Timestamp.After(f.opts.To) {
		return false
	}

	if line.Severity < f.opts.MinLevel {
		return false
	}

	return true
}

// MatchAll returns only the log lines from the provided slice that satisfy
// the filter criteria. It returns an empty (non-nil) slice if none match.
func (f *Filter) MatchAll(lines []*parser.LogLine) []*parser.LogLine {
	result := make([]*parser.LogLine, 0, len(lines))
	for _, line := range lines {
		if f.Match(line) {
			result = append(result, line)
		}
	}
	return result
}
