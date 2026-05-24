// Package coalesce merges consecutive log lines that belong to the same
// logical event (e.g. multi-line stack traces) into a single LogLine.
//
// A new logical event begins whenever a line matches the primary timestamp
// pattern. Lines that do not start with a timestamp are appended to the
// previous line's message, separated by a newline.
package coalesce

import (
	"strings"
	"sync"

	"github.com/yourorg/logslice/internal/parser"
)

// Coalescer buffers incoming LogLines and emits merged events.
type Coalescer struct {
	mu      sync.Mutex
	pending *parser.LogLine
	maxLines int
}

// New returns a Coalescer. maxLines is the maximum number of raw lines that
// may be merged into a single event (0 means unlimited).
func New(maxLines int) *Coalescer {
	return &Coalescer{maxLines: maxLines}
}

// Push accepts a parsed LogLine. If line is nil it is ignored.
// When a new timestamped line arrives the previously buffered event is
// returned (flushed); otherwise nil is returned and the continuation line
// is appended to the buffer.
func (c *Coalescer) Push(line *parser.LogLine) *parser.LogLine {
	if line == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	// A line is a continuation if it has no timestamp (zero value).
	isContinuation := line.Timestamp.IsZero()

	if isContinuation && c.pending != nil {
		// Check merge cap.
		if c.maxLines <= 0 || strings.Count(c.pending.Message, "\n") < c.maxLines-1 {
			c.pending.Message += "\n" + line.Message
		}
		return nil
	}

	// New timestamped line — flush whatever was pending.
	flushed := c.pending
	c.pending = line
	return flushed
}

// Flush returns the currently buffered event and clears the buffer.
// Call this after the last line has been pushed to drain any remaining event.
func (c *Coalescer) Flush() *parser.LogLine {
	c.mu.Lock()
	defer c.mu.Unlock()
	flushed := c.pending
	c.pending = nil
	return flushed
}
