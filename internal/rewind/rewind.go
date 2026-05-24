package rewind

import (
	"fmt"
	"io"

	"github.com/yourorg/logslice/internal/parser"
)

// Rewinder scans a log file backwards from a given offset, yielding
// up to N lines in chronological order.
type Rewinder struct {
	r      io.ReadSeeker
	p      *parser.Parser
	chunkSz int64
}

// New returns a Rewinder that reads from r using p to parse lines.
// chunkSize controls how many bytes are read per backward step.
func New(r io.ReadSeeker, p *parser.Parser, chunkSize int64) (*Rewinder, error) {
	if r == nil {
		return nil, fmt.Errorf("rewind: reader must not be nil")
	}
	if p == nil {
		return nil, fmt.Errorf("rewind: parser must not be nil")
	}
	if chunkSize <= 0 {
		chunkSize = 4096
	}
	return &Rewinder{r: r, p: p, chunkSz: chunkSize}, nil
}

// Last returns up to n log lines ending at endOffset (exclusive),
// returned in chronological (oldest-first) order.
func (rw *Rewinder) Last(endOffset int64, n int) ([]*parser.LogLine, error) {
	if n <= 0 {
		return nil, nil
	}

	lines, err := rw.collectLines(endOffset, n)
	if err != nil {
		return nil, err
	}
	return lines, nil
}

// collectLines walks backwards through the file gathering raw lines,
// then parses and returns up to n of them in forward order.
func (rw *Rewinder) collectLines(endOffset int64, n int) ([]*parser.LogLine, error) {
	var raw []string
	pos := endOffset
	remainder := ""

	for pos > 0 && len(raw) <= n {
		step := rw.chunkSz
		if pos < step {
			step = pos
		}
		pos -= step

		buf := make([]byte, step)
		if _, err := rw.r.Seek(pos, io.SeekStart); err != nil {
			return nil, fmt.Errorf("rewind: seek: %w", err)
		}
		if _, err := io.ReadFull(rw.r, buf); err != nil {
			return nil, fmt.Errorf("rewind: read: %w", err)
		}

		chunk := string(buf) + remainder
		parts := splitLines(chunk)
		remainder = parts[0]
		for i := len(parts) - 1; i >= 1; i-- {
			raw = append(raw, parts[i])
		}
	}

	if remainder != "" {
		raw = append(raw, remainder)
	}

	// reverse so oldest first
	reverse(raw)
	if len(raw) > n {
		raw = raw[len(raw)-n:]
	}

	var out []*parser.LogLine
	for _, s := range raw {
		if s == "" {
			continue
		}
		line, err := rw.p.Parse(s)
		if err != nil {
			continue
		}
		out = append(out, line)
	}
	return out, nil
}

func splitLines(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

func reverse(ss []string) {
	for i, j := 0, len(ss)-1; i < j; i, j = i+1, j-1 {
		ss[i], ss[j] = ss[j], ss[i]
	}
}
