// Package burst detects sudden spikes in log volume over a rolling window.
// It tracks the number of lines seen in the last N seconds and fires a
// callback when the rate exceeds a configured threshold.
package burst

import (
	"sync"
	"time"

	"github.com/user/logslice/internal/parser"
)

// OnBurst is called when a burst is detected. rate is lines-per-second.
type OnBurst func(rate float64)

// Detector tracks log line arrival rate and fires OnBurst when the
// instantaneous rate exceeds Threshold lines per second.
type Detector struct {
	mu        sync.Mutex
	threshold float64
	window    time.Duration
	callback  OnBurst
	buckets   []entry
	fired     bool
}

type entry struct {
	t time.Time
}

// New creates a Detector that fires cb when the rate over the given window
// exceeds threshold lines/sec. window must be > 0.
func New(threshold float64, window time.Duration, cb OnBurst) *Detector {
	if window <= 0 {
		window = time.Second
	}
	return &Detector{
		threshold: threshold,
		window:    window,
		callback:  cb,
	}
}

// Record registers a parsed log line arrival. It evicts stale buckets,
// computes the current rate, and fires the callback if the threshold is
// exceeded. A nil line is silently ignored.
func (d *Detector) Record(line *parser.LogLine) {
	if line == nil {
		return
	}
	now := time.Now()
	d.mu.Lock()
	defer d.mu.Unlock()

	cutoff := now.Add(-d.window)
	d.buckets = append(d.buckets, entry{t: now})

	// evict entries outside the window
	start := 0
	for start < len(d.buckets) && d.buckets[start].t.Before(cutoff) {
		start++
	}
	d.buckets = d.buckets[start:]

	rate := float64(len(d.buckets)) / d.window.Seconds()
	if rate >= d.threshold {
		if !d.fired && d.callback != nil {
			d.fired = true
			go d.callback(rate)
		}
	} else {
		d.fired = false
	}
}

// Rate returns the current lines-per-second rate over the configured window.
func (d *Detector) Rate() float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	cutoff := time.Now().Add(-d.window)
	count := 0
	for _, e := range d.buckets {
		if !e.t.Before(cutoff) {
			count++
		}
	}
	return float64(count) / d.window.Seconds()
}

// Reset clears all recorded entries and resets the fired flag.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.buckets = d.buckets[:0]
	d.fired = false
}
