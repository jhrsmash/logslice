package slicer

import (
	"io"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/reader"
)

// Options holds the configuration for a slice operation.
type Options struct {
	FilePath  string
	Filter    *filter.Filter
	Writer    io.Writer
}

// Slicer reads a log file and writes matching lines to the configured writer.
type Slicer struct {
	opts Options
	p    *parser.Parser
}

// New creates a new Slicer with the given options.
func New(opts Options) *Slicer {
	return &Slicer{
		opts: opts,
		p:    parser.New(),
	}
}

// Run performs the slice operation: it iterates over log lines, applies the
// filter, and writes matching lines to the writer. It returns the number of
// lines written and any error encountered.
func (s *Slicer) Run() (int, error) {
	r, err := reader.New(s.opts.FilePath)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	written := 0
	for {
		raw, readErr := r.ReadLine()
		if readErr != nil && readErr != io.EOF {
			return written, readErr
		}

		if len(raw) > 0 {
			line, parseErr := s.p.Parse(raw)
			if parseErr == nil && s.opts.Filter.Match(line) {
				if _, wErr := s.opts.Writer.Write(append(raw, '\n')); wErr != nil {
					return written, wErr
				}
				written++
			}
		}

		if readErr == io.EOF {
			break
		}
	}

	return written, nil
}
