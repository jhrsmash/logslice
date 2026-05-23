// Package throttle provides a line-based throughput throttler that limits
// the number of log lines emitted per second across the processing pipeline.
package throttle

import (
	"sync"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// Throttler limits the rate at which log lines are passed downstream.
// When the configured lines-per-second budget is exhausted the caller
// blocks until the next one-second window opens.
type Throttler struct {
	mu      sync.Mutex
	rate    int // max lines per second; 0 means unlimited
	count   int
	window  time.Time
	clock   func() time.Time
}

// New returns a Throttler that allows at most linesPerSec lines through per
// second. A rate of 0 disables throttling entirely.
func New(linesPerSec int) *Throttler {
	return &Throttler{
		rate:  linesPerSec,
		clock: time.Now,
	}
}

// Allow reports whether the given line should be forwarded. When the rate
// limit is active and the current window is full, Allow blocks until the
// next window begins and then permits the line.
// A nil line is always allowed through unchanged.
func (t *Throttler) Allow(line *parser.LogLine) bool {
	if line == nil {
		return true
	}
	if t.rate <= 0 {
		return true
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.clock()

	// Start or advance the window.
	if t.window.IsZero() || now.After(t.window) {
		t.window = now.Add(time.Second)
		t.count = 0
	}

	if t.count < t.rate {
		t.count++
		return true
	}

	// Budget exhausted — sleep until the window expires then allow.
	sleepFor := t.window.Sub(t.clock())
	if sleepFor > 0 {
		time.Sleep(sleepFor)
	}
	t.window = t.clock().Add(time.Second)
	t.count = 1
	return true
}

// Rate returns the configured lines-per-second limit.
func (t *Throttler) Rate() int {
	return t.rate
}
