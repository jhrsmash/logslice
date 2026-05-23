// Package fieldextract provides utilities for extracting named fields
// from structured log message payloads (e.g. key=value or JSON fragments).
package fieldextract

import (
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Extractor holds compiled field extraction configuration.
type Extractor struct {
	fields []string // field names to extract; empty means all
}

// New returns an Extractor that will pull the given field names from a log
// line's message. Pass no names to extract every key=value pair found.
func New(fields ...string) *Extractor {
	return &Extractor{fields: fields}
}

// Extract parses key=value pairs from line.Message and returns a map
// containing only the requested fields (or all fields when none were
// specified). Quoted values ("val") are unquoted automatically.
// Returns nil when line is nil.
func (e *Extractor) Extract(line *parser.LogLine) map[string]string {
	if line == nil {
		return nil
	}

	pairs := parseKV(line.Message)

	if len(e.fields) == 0 {
		return pairs
	}

	out := make(map[string]string, len(e.fields))
	for _, f := range e.fields {
		if v, ok := pairs[f]; ok {
			out[f] = v
		}
	}
	return out
}

// parseKV scans a string for key=value tokens separated by whitespace.
// Values may be bare words or double-quoted strings.
func parseKV(s string) map[string]string {
	out := make(map[string]string)
	for _, token := range strings.Fields(s) {
		idx := strings.IndexByte(token, '=')
		if idx <= 0 || idx == len(token)-1 {
			continue
		}
		key := token[:idx]
		val := token[idx+1:]
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}
		out[key] = val
	}
	return out
}
