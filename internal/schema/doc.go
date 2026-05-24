// Package schema provides named-capture-group pattern matching for log lines.
//
// A Schema wraps a compiled regular expression whose named capture groups
// define the fields present in a log line format (e.g. timestamp, severity,
// message). Multiple schemas can be collected in a Registry, which can
// attempt to match an incoming line against all registered schemas and
// return the first successful match.
//
// Typical usage:
//
//	s, err := schema.New("nginx", `(?P<ts>[\d/: ]+) \[(?P<level>\w+)\] (?P<msg>.+)`)
//	if err != nil { ... }
//	reg := schema.NewRegistry()
//	reg.Register(s)
//	name, fields := reg.MatchFirst(line)
package schema
