// Package mask redacts sensitive values from log lines before they reach any
// output sink.
//
// # Overview
//
// A Masker holds a list of compiled regular expressions. Each expression must
// contain exactly one capturing group that isolates the sensitive value within
// the surrounding context. When Mask is called the capturing group's match is
// replaced with the string "[REDACTED]" while the rest of the line is left
// intact.
//
// # Example patterns
//
//	// Mask the value of a "password" key=value pair.
//	`password=([^\s]+)`
//
//	// Mask a Bearer token in an Authorization header.
//	`(?i)Authorization:\s*Bearer\s+(\S+)`
//
// Patterns are applied in registration order. AddPattern is safe for
// concurrent use.
package mask
