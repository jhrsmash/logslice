// Package aggregate provides log line counting and grouping by severity
// over a sliding or fixed time window.
package aggregate

import (
	"sync"
	"time"

	"github.com/user/logslice/internal/parser"
)

// Bucket holds counts for a single time window bucket.
type Bucket struct {
	Start    time.Time
	Counts   map[parser.Severity]int
	Total    int
}

// Aggregator groups log lines into fixed-width time buckets and counts
// occurrences per severity level.
type Aggregator struct {
	mu       sync.Mutex
	buckets  map[int64]*Bucket
	window   time.Duration
}

// New creates a new Aggregator that groups lines into buckets of the given
// window duration (e.g. time.Minute).
func New(window time.Duration) *Aggregator {
	if window <= 0 {
		window = time.Minute
	}
	return &Aggregator{
		buckets: make(map[int64]*Bucket),
		window:  window,
	}
}

// Record adds a log line's severity into the appropriate time bucket.
// Lines with a zero timestamp are ignored.
func (a *Aggregator) Record(ts time.Time, sev parser.Severity) {
	if ts.IsZero() {
		return
	}
	key := ts.Truncate(a.window).UnixNano()

	a.mu.Lock()
	defer a.mu.Unlock()

	b, ok := a.buckets[key]
	if !ok {
		b = &Bucket{
			Start:  ts.Truncate(a.window),
			Counts: make(map[parser.Severity]int),
		}
		a.buckets[key] = b
	}
	b.Counts[sev]++
	b.Total++
}

// Snapshot returns a sorted copy of all buckets collected so far.
func (a *Aggregator) Snapshot() []Bucket {
	a.mu.Lock()
	defer a.mu.Unlock()

	result := make([]Bucket, 0, len(a.buckets))
	for _, b := range a.buckets {
		copy := Bucket{
			Start:  b.Start,
			Total:  b.Total,
			Counts: make(map[parser.Severity]int, len(b.Counts)),
		}
		for k, v := range b.Counts {
			copy.Counts[k] = v
		}
		result = append(result, copy)
	}
	sortBuckets(result)
	return result
}

// Reset clears all accumulated buckets.
func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.buckets = make(map[int64]*Bucket)
}

func sortBuckets(bs []Bucket) {
	for i := 1; i < len(bs); i++ {
		for j := i; j > 0 && bs[j].Start.Before(bs[j-1].Start); j-- {
			bs[j], bs[j-1] = bs[j-1], bs[j]
		}
	}
}
