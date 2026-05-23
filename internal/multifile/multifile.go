// Package multifile provides ordered iteration over multiple log files,
// merging their lines in chronological order.
package multifile

import (
	"fmt"
	"io"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/reader"
)

// Source wraps a reader and its current pending line.
type Source struct {
	r    *reader.Reader
	line *parser.LogLine
	done bool
}

// MultiReader merges lines from multiple log files in timestamp order.
type MultiReader struct {
	sources []*Source
	p       *parser.Parser
}

// New creates a MultiReader over the given file paths.
// Files must each be individually sorted by timestamp.
func New(paths []string) (*MultiReader, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("multifile: at least one path required")
	}
	mr := &MultiReader{p: parser.New()}
	for _, path := range paths {
		r, err := reader.New(path)
		if err != nil {
			mr.Close()
			return nil, fmt.Errorf("multifile: open %q: %w", path, err)
		}
		mr.sources = append(mr.sources, &Source{r: r})
	}
	// Prime each source with its first line.
	for _, s := range mr.sources {
		mr.advance(s)
	}
	return mr, nil
}

// Next returns the earliest-timestamped line across all sources.
// Returns (nil, io.EOF) when all sources are exhausted.
func (mr *MultiReader) Next() (*parser.LogLine, error) {
	var best *Source
	for _, s := range mr.sources {
		if s.done || s.line == nil {
			continue
		}
		if best == nil || s.line.Timestamp.Before(best.line.Timestamp) {
			best = s
		}
	}
	if best == nil {
		return nil, io.EOF
	}
	line := best.line
	mr.advance(best)
	return line, nil
}

// Close releases all underlying readers.
func (mr *MultiReader) Close() {
	for _, s := range mr.sources {
		if s.r != nil {
			s.r.Close()
		}
	}
}

func (mr *MultiReader) advance(s *Source) {
	raw, err := s.r.ReadLine()
	if err != nil {
		s.line = nil
		s.done = true
		return
	}
	line, err := mr.p.Parse(raw)
	if err != nil {
		s.line = nil
		s.done = true
		return
	}
	s.line = line
}
