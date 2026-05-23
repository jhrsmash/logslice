// Package levelcount tracks per-severity line counts across a stream of log lines.
package levelcount

import (
	"sync"

	"github.com/user/logslice/internal/parser"
)

// Counter accumulates counts of log lines by severity level.
type Counter struct {
	mu     sync.Mutex
	counts map[parser.Severity]int64
}

// New returns an initialised Counter.
func New() *Counter {
	return &Counter{
		counts: make(map[parser.Severity]int64),
	}
}

// Record increments the count for the severity of line.
// A nil line is silently ignored.
func (c *Counter) Record(line *parser.LogLine) {
	if line == nil {
		return
	}
	c.mu.Lock()
	c.counts[line.Severity]++
	c.mu.Unlock()
}

// Snapshot returns a copy of the current counts keyed by severity.
func (c *Counter) Snapshot() map[parser.Severity]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[parser.Severity]int64, len(c.counts))
	for k, v := range c.counts {
		out[k] = v
	}
	return out
}

// Reset zeroes all accumulated counts.
func (c *Counter) Reset() {
	c.mu.Lock()
	c.counts = make(map[parser.Severity]int64)
	c.mu.Unlock()
}

// Total returns the sum of all severity counts.
func (c *Counter) Total() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	var n int64
	for _, v := range c.counts {
		n += v
	}
	return n
}
