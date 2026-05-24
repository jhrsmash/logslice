// Package fanout provides a multiplexer that forwards each log line to
// multiple downstream consumers concurrently.
package fanout

import (
	"sync"

	"github.com/yourorg/logslice/internal/parser"
)

// Sink is any consumer that can receive a parsed log line.
type Sink interface {
	Write(line *parser.LogLine) error
}

// Fanout distributes a single stream of log lines to N sinks.
// Each sink receives every line; errors are collected but do not stop
// delivery to the remaining sinks.
type Fanout struct {
	mu    sync.RWMutex
	sinks []Sink
}

// New returns an empty Fanout. Use Add to register sinks before sending.
func New(sinks ...Sink) *Fanout {
	f := &Fanout{}
	for _, s := range sinks {
		if s != nil {
			f.sinks = append(f.sinks, s)
		}
	}
	return f
}

// Add registers an additional sink. It is safe to call concurrently.
func (f *Fanout) Add(s Sink) {
	if s == nil {
		return
	}
	f.mu.Lock()
	f.sinks = append(f.sinks, s)
	f.mu.Unlock()
}

// Len returns the number of registered sinks.
func (f *Fanout) Len() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.sinks)
}

// Send forwards line to every registered sink.
// A nil line is silently ignored.
// All errors encountered are returned as a slice; a nil return means
// every sink accepted the line without error.
func (f *Fanout) Send(line *parser.LogLine) []error {
	if line == nil {
		return nil
	}
	f.mu.RLock()
	sinks := f.sinks
	f.mu.RUnlock()

	var errs []error
	for _, s := range sinks {
		if err := s.Write(line); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
