package reader

import (
	"io"
	"time"

	"github.com/user/logslice/internal/parser"
)

// SeekToTime performs a binary search over the file to find the byte offset
// of the first log line whose timestamp is >= target. The file must be
// roughly time-ordered for the search to be meaningful.
//
// Returns 0 if no suitable position is found or the file is too small.
func SeekToTime(r *Reader, p *parser.Parser, target time.Time) (int64, error) {
	size, err := r.Size()
	if err != nil {
		return 0, err
	}
	if size == 0 {
		return 0, nil
	}

	lo, hi := int64(0), size
	result := int64(0)

	for lo < hi {
		mid := (lo + hi) / 2

		t, err := timestampAtOffset(r, p, mid)
		if err != nil {
			// Can't parse here; move right
			lo = mid + 1
			continue
		}

		if t.Before(target) {
			lo = mid + 1
		} else {
			result = mid
			hi = mid
		}
	}

	return result, nil
}

// timestampAtOffset seeks to offset, skips the (possibly partial) current line,
// reads the next full line, and returns its parsed timestamp.
func timestampAtOffset(r *Reader, p *parser.Parser, offset int64) (time.Time, error) {
	if err := r.SeekOffset(offset); err != nil {
		return time.Time{}, err
	}
	// Skip the potentially partial line (unless we're at the start)
	if offset != 0 {
		if _, err := r.ReadLine(); err != nil {
			return time.Time{}, err
		}
	}
	line, err := r.ReadLine()
	if err != nil {
		if err == io.EOF {
			return time.Time{}, io.EOF
		}
		return time.Time{}, err
	}
	ll, err := p.Parse(line)
	if err != nil {
		return time.Time{}, err
	}
	return ll.Timestamp, nil
}
