// Package dedupe provides a sliding-window duplicate-line filter.
// Consecutive log lines with identical message bodies are collapsed
// into a single line; a suppression count annotation is appended when
// more than one duplicate is dropped.
package dedupe

import (
	"fmt"

	"github.com/user/logslice/internal/parser"
)

// Deduplicator tracks the last seen log message and suppresses runs of
// identical messages. Call Flush after the final line to emit any
// pending suppressed entry.
type Deduplicator struct {
	last    *parser.LogLine
	count   int
	emit    func(*parser.LogLine)
}

// New returns a Deduplicator that calls emit for each unique (or
// annotated-duplicate) line.
func New(emit func(*parser.LogLine)) *Deduplicator {
	return &Deduplicator{emit: emit}
}

// Feed processes a single log line. Identical consecutive messages are
// buffered; the first occurrence is emitted immediately.
func (d *Deduplicator) Feed(line *parser.LogLine) {
	if line == nil {
		return
	}
	if d.last != nil && line.Message == d.last.Message {
		d.count++
		return
	}
	d.flushPending()
	d.emit(line)
	d.last = line
	d.count = 0
}

// Flush emits a suppression annotation for any buffered duplicates.
// It must be called once after all lines have been fed.
func (d *Deduplicator) Flush() {
	d.flushPending()
}

// flushPending emits an annotated copy of the last line when duplicates
// were suppressed.
func (d *Deduplicator) flushPending() {
	if d.last == nil || d.count == 0 {
		return
	}
	annotated := *d.last
	annotated.Message = fmt.Sprintf("%s [repeated %d time(s)]", d.last.Message, d.count)
	d.emit(&annotated)
	d.count = 0
}
