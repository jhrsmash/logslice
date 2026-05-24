// Package redact provides a log line redactor that replaces sensitive
// field values with a fixed placeholder before lines are emitted.
//
// Patterns are applied in registration order. Each pattern is a compiled
// regular expression; every non-overlapping match in the raw log line is
// replaced with the configured placeholder (default "[REDACTED]").
package redact

import (
	"fmt"
	"regexp"

	"github.com/user/logslice/internal/parser"
)

const defaultPlaceholder = "[REDACTED]"

// Redactor replaces sensitive content in log lines.
type Redactor struct {
	patterns     []*regexp.Regexp
	placeholder  []byte
}

// New compiles each pattern string and returns a Redactor.
// An error is returned if any pattern fails to compile.
func New(patterns []string, placeholder string) (*Redactor, error) {
	if placeholder == "" {
		placeholder = defaultPlaceholder
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("redact: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return &Redactor{
		patterns:    compiled,
		placeholder: []byte(placeholder),
	}, nil
}

// Redact returns a copy of line with all sensitive patterns replaced.
// If line is nil or no patterns are registered the original pointer is
// returned unchanged.
func (r *Redactor) Redact(line *parser.LogLine) *parser.LogLine {
	if line == nil || len(r.patterns) == 0 {
		return line
	}
	raw := []byte(line.Raw)
	for _, re := range r.patterns {
		raw = re.ReplaceAll(raw, r.placeholder)
	}
	if string(raw) == line.Raw {
		return line
	}
	copy := *line
	copy.Raw = string(raw)
	return &copy
}
