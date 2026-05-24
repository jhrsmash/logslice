// Package normalize provides log line field normalization,
// standardizing field names and values across different log formats.
package normalize

import (
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule describes a single normalization: rename a field key and optionally
// transform its value via a replacer function.
type Rule struct {
	// From is the source field name to match (case-insensitive).
	From string
	// To is the canonical field name to write.
	To string
	// Transform is an optional value transformer. If nil the value is kept as-is.
	Transform func(string) string
}

// Normalizer applies a set of Rules to each log line's extra fields.
type Normalizer struct {
	rules []Rule
	// index maps lowercase source key to rule index for O(1) lookup.
	index map[string]int
}

// New creates a Normalizer from the given rules.
// Duplicate From keys are last-write-wins.
func New(rules []Rule) *Normalizer {
	idx := make(map[string]int, len(rules))
	for i, r := range rules {
		idx[strings.ToLower(r.From)] = i
	}
	return &Normalizer{rules: rules, index: idx}
}

// Apply rewrites the Fields map of line in-place and returns the same pointer.
// If line is nil, Apply is a no-op and returns nil.
func (n *Normalizer) Apply(line *parser.LogLine) *parser.LogLine {
	if line == nil {
		return nil
	}
	if len(line.Fields) == 0 {
		return line
	}

	newFields := make(map[string]string, len(line.Fields))
	for k, v := range line.Fields {
		if i, ok := n.index[strings.ToLower(k)]; ok {
			r := n.rules[i]
			if r.Transform != nil {
				v = r.Transform(v)
			}
			newFields[r.To] = v
		} else {
			newFields[k] = v
		}
	}
	line.Fields = newFields
	return line
}
