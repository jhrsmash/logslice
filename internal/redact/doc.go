// Package redact provides regex-based redaction of sensitive values inside
// raw log lines.
//
// # Usage
//
//	r, err := redact.New([]string{`password=\S+`, `token=\S+`}, "")
//	if err != nil {
//		log.Fatal(err)
//	}
//	clean := r.Redact(line)
//
// Redact never modifies the original LogLine; it returns a shallow copy with
// the Raw field updated when at least one substitution was made.
package redact
