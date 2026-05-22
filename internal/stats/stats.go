// Package stats collects and reports processing statistics for a logslice run.
package stats

import (
	"fmt"
	"io"
	"time"
)

// Stats tracks counters accumulated during log slicing.
type Stats struct {
	LinesRead    int64
	LinesMatched int64
	BytesRead    int64
	Started      time.Time
	Finished     time.Time
}

// New returns a new Stats instance with the start time recorded.
func New() *Stats {
	return &Stats{Started: time.Now()}
}

// RecordLine increments the lines-read counter and adds n to bytes read.
func (s *Stats) RecordLine(n int) {
	s.LinesRead++
	s.BytesRead += int64(n)
}

// RecordMatch increments the lines-matched counter.
func (s *Stats) RecordMatch() {
	s.LinesMatched++
}

// Finish records the finish time. Safe to call multiple times; only the first
// call has effect.
func (s *Stats) Finish() {
	if s.Finished.IsZero() {
		s.Finished = time.Now()
	}
}

// Elapsed returns the duration between start and finish (or now if not finished).
func (s *Stats) Elapsed() time.Duration {
	end := s.Finished
	if end.IsZero() {
		end = time.Now()
	}
	return end.Sub(s.Started)
}

// Write prints a human-readable summary to w.
func (s *Stats) Write(w io.Writer) error {
	_, err := fmt.Fprintf(
		w,
		"lines read: %d | lines matched: %d | bytes read: %d | elapsed: %s\n",
		s.LinesRead,
		s.LinesMatched,
		s.BytesRead,
		s.Elapsed().Round(time.Millisecond),
	)
	return err
}
