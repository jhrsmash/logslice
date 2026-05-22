package index

import (
	"bufio"
	"io"
	"time"

	"github.com/user/logslice/internal/parser"
)

const defaultSampleEvery = 500

// BuildOptions controls how the index is constructed.
type BuildOptions struct {
	// SampleEvery indexes one entry per this many lines (default 500).
	SampleEvery int
}

// Build reads r from its current position, parsing timestamps and recording
// one index entry every opts.SampleEvery lines. The returned Index can then
// be used for fast offset lookups.
//
// r must support io.ReadSeeker so that the current offset can be tracked.
func Build(r io.ReadSeeker, p *parser.Parser, opts BuildOptions) (*Index, error) {
	if opts.SampleEvery <= 0 {
		opts.SampleEvery = defaultSampleEvery
	}

	idx := New()
	scanner := bufio.NewScanner(r)

	var (
		lineNum    int
		byteOffset int64
		lastTS     time.Time
	)

	for scanner.Scan() {
		raw := scanner.Text()
		lineLen := int64(len(raw)) + 1 // +1 for newline

		if lineNum%opts.SampleEvery == 0 {
			line, err := p.Parse(raw)
			if err == nil && line != nil {
				if line.Timestamp.After(lastTS) || lineNum == 0 {
					idx.Add(line.Timestamp, byteOffset)
					lastTS = line.Timestamp
				}
			}
		}

		byteOffset += lineLen
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return idx, nil
}
