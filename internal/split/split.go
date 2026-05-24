// Package split provides a log-line splitter that partitions a stream of
// parsed log lines into named buckets based on a user-supplied key function.
// Each bucket is an independent slice that can be written to separate outputs.
package split

import (
	"errors"
	"sync"

	"github.com/example/logslice/internal/parser"
)

// KeyFunc returns the bucket name for a given log line.
// Lines for which KeyFunc returns an empty string are dropped.
type KeyFunc func(line *parser.LogLine) string

// Splitter partitions log lines into named buckets.
type Splitter struct {
	mu      sync.Mutex
	keyFn   KeyFunc
	buckets map[string][]*parser.LogLine
}

// New creates a Splitter that uses keyFn to assign lines to buckets.
// Returns an error if keyFn is nil.
func New(keyFn KeyFunc) (*Splitter, error) {
	if keyFn == nil {
		return nil, errors.New("split: KeyFunc must not be nil")
	}
	return &Splitter{
		keyFn:   keyFn,
		buckets: make(map[string][]*parser.LogLine),
	}, nil
}

// Add classifies line into a bucket. Nil lines and lines whose key is empty
// are silently ignored.
func (s *Splitter) Add(line *parser.LogLine) {
	if line == nil {
		return
	}
	key := s.keyFn(line)
	if key == "" {
		return
	}
	s.mu.Lock()
	s.buckets[key] = append(s.buckets[key], line)
	s.mu.Unlock()
}

// Bucket returns a copy of the lines stored under key.
// Returns nil if the bucket does not exist.
func (s *Splitter) Bucket(key string) []*parser.LogLine {
	s.mu.Lock()
	defer s.mu.Unlock()
	lines, ok := s.buckets[key]
	if !ok {
		return nil
	}
	out := make([]*parser.LogLine, len(lines))
	copy(out, lines)
	return out
}

// Keys returns the sorted list of bucket names that have at least one line.
func (s *Splitter) Keys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.buckets))
	for k := range s.buckets {
		keys = append(keys, k)
	}
	return keys
}

// Reset clears all buckets.
func (s *Splitter) Reset() {
	s.mu.Lock()
	s.buckets = make(map[string][]*parser.LogLine)
	s.mu.Unlock()
}
