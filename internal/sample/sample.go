// Package sample provides log line sampling by emitting every Nth matched line.
// This is useful for high-volume log files where only a representative subset
// of matching lines is needed.
package sample

import "github.com/yourorg/logslice/internal/parser"

// Sampler wraps an output channel and forwards only every Nth line.
type Sampler struct {
	rate    int
	counter int
}

// New creates a Sampler that emits one line for every n lines it receives.
// A rate of 0 or 1 disables sampling (all lines are forwarded).
func New(rate int) *Sampler {
	if rate < 1 {
		rate = 1
	}
	return &Sampler{rate: rate}
}

// Allow returns true if the current line should be forwarded to output.
// It increments an internal counter on every call, regardless of the result.
func (s *Sampler) Allow(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	s.counter++
	return s.counter%s.rate == 0
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() int {
	return s.rate
}

// Reset resets the internal counter, restarting the sampling window.
func (s *Sampler) Reset() {
	s.counter = 0
}
