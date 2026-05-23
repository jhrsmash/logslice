// Package grep provides substring and regular-expression matching
// against the raw text of a parsed log line.
package grep

import (
	"fmt"
	"regexp"

	"github.com/user/logslice/internal/parser"
)

// Matcher decides whether a log line's raw text satisfies a search
// expression. A nil Matcher always matches.
type Matcher struct {
	re      *regexp.Regexp
	invert  bool
}

// New compiles pattern into a Matcher. If pattern is empty the returned
// Matcher is a no-op (matches everything). When invert is true the
// Matcher accepts lines that do NOT match the pattern.
func New(pattern string, invert bool) (*Matcher, error) {
	if pattern == "" {
		return &Matcher{}, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("grep: invalid pattern %q: %w", pattern, err)
	}
	return &Matcher{re: re, invert: invert}, nil
}

// Match reports whether line satisfies the matcher.
// A nil line or a Matcher with no pattern always returns true.
func (m *Matcher) Match(line *parser.LogLine) bool {
	if line == nil {
		return true
	}
	if m == nil || m.re == nil {
		return true
	}
	hit := m.re.MatchString(line.Raw)
	if m.invert {
		return !hit
	}
	return hit
}

// String returns a human-readable description of the matcher.
func (m *Matcher) String() string {
	if m == nil || m.re == nil {
		return "<no pattern>"
	}
	if m.invert {
		return fmt.Sprintf("NOT /%s/", m.re.String())
	}
	return fmt.Sprintf("/%s/", m.re.String())
}
