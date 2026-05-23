package window

import (
	"sync"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// Window holds a rolling time window of log lines, evicting entries
// older than the configured duration.
type Window struct {
	mu       sync.Mutex
	duration time.Duration
	entries  []*parser.LogLine
}

// New creates a Window that retains lines within the given duration.
// A zero or negative duration retains all lines (no eviction).
func New(d time.Duration) *Window {
	return &Window{duration: d}
}

// Add appends a log line to the window and evicts stale entries.
func (w *Window) Add(line *parser.LogLine) {
	if line == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries = append(w.entries, line)
	w.evict(line.Timestamp)
}

// Snapshot returns a copy of the current window contents.
func (w *Window) Snapshot() []*parser.LogLine {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]*parser.LogLine, len(w.entries))
	copy(out, w.entries)
	return out
}

// Len returns the number of lines currently in the window.
func (w *Window) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.entries)
}

// Reset clears all entries from the window.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries = w.entries[:0]
}

// evict removes entries whose timestamp is older than (pivot - duration).
// Must be called with w.mu held.
func (w *Window) evict(pivot time.Time) {
	if w.duration <= 0 {
		return
	}
	cutoff := pivot.Add(-w.duration)
	i := 0
	for i < len(w.entries) && w.entries[i].Timestamp.Before(cutoff) {
		i++
	}
	if i > 0 {
		w.entries = w.entries[i:]
	}
}
