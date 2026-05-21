package reader

import (
	"io"
	"time"

	"github.com/example/logslice/internal/parser"
)

// SeekToTime performs a binary search over the file to find the offset of the
// first log line whose timestamp is >= target. Returns the byte offset, or
// io.EOF if no such line exists.
func SeekToTime(r *Reader, target time.Time) (int64, error) {
	size := r.Size()
	if size == 0 {
		return 0, io.EOF
	}

	lo, hi := int64(0), size

	for lo < hi {
		mid := (lo + hi) / 2
		offset, err := snapToLineStart(r, mid)
		if err != nil {
			return 0, err
		}

		t, err := timestampAtOffset(r, offset)
		if err != nil {
			// No valid line found in upper half; shrink
			hi = mid
			continue
		}

		if t.Before(target) {
			lo = offset + 1
		} else {
			hi = offset
		}
	}

	if lo >= size {
		return 0, io.EOF
	}

	// Snap lo to a clean line boundary
	offset, err := snapToLineStart(r, lo)
	if err != nil {
		return 0, err
	}
	return offset, nil
}

// timestampAtOffset reads the log line starting at offset and parses its timestamp.
func timestampAtOffset(r *Reader, offset int64) (time.Time, error) {
	line, err := r.ReadLineAt(offset)
	if err != nil {
		return time.Time{}, err
	}
	p := parser.New()
	parsed, err := p.Parse(line)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.Timestamp, nil
}

// snapToLineStart scans backwards from offset to find the start of the
// enclosing line. If offset is already at position 0 it is returned as-is.
func snapToLineStart(r *Reader, offset int64) (int64, error) {
	if offset <= 0 {
		return 0, nil
	}

	// Walk backwards until we find a newline or the beginning of the file.
	for offset > 0 {
		offset--
		b, err := r.ReadByteAt(offset)
		if err != nil {
			return 0, err
		}
		if b == '\n' {
			return offset + 1, nil
		}
	}
	return 0, nil
}
