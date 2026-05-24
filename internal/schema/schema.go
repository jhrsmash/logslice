package schema

import (
	"fmt"
	"regexp"
	"strings"
)

// Field represents a named capture group in a log schema pattern.
type Field struct {
	Name     string
	Optional bool
}

// Schema holds a compiled log line pattern with named capture groups.
type Schema struct {
	Name    string
	pattern *regexp.Regexp
	fields  []Field
}

// New compiles a named-capture-group regex pattern into a Schema.
// Pattern must use (?P<name>...) syntax for named groups.
func New(name, pattern string) (*Schema, error) {
	if name == "" {
		return nil, fmt.Errorf("schema name must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid schema pattern: %w", err)
	}
	var fields []Field
	for _, n := range re.SubexpNames() {
		if n != "" {
			fields = append(fields, Field{Name: n})
		}
	}
	if len(fields) == 0 {
		return nil, fmt.Errorf("schema pattern must contain at least one named capture group")
	}
	return &Schema{Name: name, pattern: re, fields: fields}, nil
}

// Match attempts to match raw against the schema pattern.
// Returns a map of field name -> value, or nil if no match.
func (s *Schema) Match(raw string) map[string]string {
	m := s.pattern.FindStringSubmatch(strings.TrimRight(raw, "\n"))
	if m == nil {
		return nil
	}
	result := make(map[string]string, len(s.fields))
	for i, name := range s.pattern.SubexpNames() {
		if name != "" {
			result[name] = m[i]
		}
	}
	return result
}

// Fields returns the list of named fields defined in the schema.
func (s *Schema) Fields() []Field {
	out := make([]Field, len(s.fields))
	copy(out, s.fields)
	return out
}
