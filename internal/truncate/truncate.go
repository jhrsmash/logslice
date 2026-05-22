// Package truncate provides utilities for truncating long log lines
// to a configurable maximum byte length, preserving valid UTF-8 boundaries.
package truncate

import "unicode/utf8"

// Truncator truncates log line messages that exceed a maximum length.
type Truncator struct {
	maxBytes int
	suffix   string
}

// New returns a Truncator that clips messages to maxBytes.
// If maxBytes <= 0 truncation is disabled (original string is returned as-is).
// suffix is appended when a line is clipped (e.g. "...").
func New(maxBytes int, suffix string) *Truncator {
	return &Truncator{
		maxBytes: maxBytes,
		suffix:   suffix,
	}
}

// Apply returns s unchanged when truncation is disabled or s fits within
// maxBytes. Otherwise it clips s at a valid UTF-8 rune boundary and appends
// the configured suffix.
func (t *Truncator) Apply(s string) string {
	if t.maxBytes <= 0 || len(s) <= t.maxBytes {
		return s
	}

	// Reserve space for the suffix so the total stays within maxBytes.
	cutAt := t.maxBytes - len(t.suffix)
	if cutAt <= 0 {
		// suffix alone exceeds the limit; return just the suffix truncated.
		if len(t.suffix) > t.maxBytes {
			return t.suffix[:t.maxBytes]
		}
		return t.suffix
	}

	// Walk back to a valid rune boundary.
	for cutAt > 0 && !utf8.RuneStart(s[cutAt]) {
		cutAt--
	}

	return s[:cutAt] + t.suffix
}

// Enabled reports whether truncation is active.
func (t *Truncator) Enabled() bool {
	return t.maxBytes > 0
}

// MaxBytes returns the configured byte limit (0 means disabled).
func (t *Truncator) MaxBytes() int {
	return t.maxBytes
}
