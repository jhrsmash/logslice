// Package mask provides redaction of sensitive fields from log lines.
// It scans the raw log message for patterns (e.g. passwords, tokens, credit
// card numbers) and replaces matched values with a fixed placeholder so that
// sensitive data is never written to output.
package mask

import (
	"regexp"
	"sync"

	"github.com/yourorg/logslice/internal/parser"
)

const redacted = "[REDACTED]"

// Masker applies a set of compiled regular expressions to each log line and
// replaces any captured group (group 1) with the redacted placeholder.
type Masker struct {
	mu       sync.RWMutex
	patterns []*regexp.Regexp
}

// New returns a Masker pre-loaded with the supplied regex patterns.
// Each pattern must contain exactly one capturing group that delimits the
// sensitive value; the surrounding text is preserved.
func New(patterns []string) (*Masker, error) {
	m := &Masker{}
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		m.patterns = append(m.patterns, re)
	}
	return m, nil
}

// AddPattern compiles and registers an additional pattern at runtime.
func (m *Masker) AddPattern(pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	m.mu.Lock()
	m.patterns = append(m.patterns, re)
	m.mu.Unlock()
	return nil
}

// Mask returns a copy of line with all sensitive values replaced.
// If line is nil it is returned unchanged.
func (m *Masker) Mask(line *parser.LogLine) *parser.LogLine {
	if line == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	raw := line.Raw
	for _, re := range m.patterns {
		raw = re.ReplaceAllStringFunc(raw, func(match string) string {
			// Replace only the first submatch (the sensitive portion).
			return re.ReplaceAllString(match, "${1}"+redacted)
		})
	}
	if raw == line.Raw {
		return line
	}
	copy := *line
	copy.Raw = raw
	return &copy
}
