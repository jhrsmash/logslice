package reader

import (
	"bufio"
	"io"
	"time"

	"github.com/example/logslice/internal/parser"
)

// SeekToTime performs a binary search over the log file to find the byte
// offset of the first line whose timestamp is >= target. It returns 0 if
// the file is empty or every line precedes target.
func SeekToTime(rs io.ReadSeeker, size int64, target time.Time) (int64, error) {
	if size == 0 {
		return 0, nil
	}

	lo, hi := int64(0), size

	for lo < hi {
		mid := (lo + hi) / 2

		t, err := timestampAtOffset(rs, mid)
		if err != nil {
			// Could not parse a timestamp at this offset; move forward.
			lo = mid + 1
			continue
		}

		if t.Before(target) {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	// Snap lo back to the start of the line that contains byte lo.
	return snapToLineStart(rs, lo)
}

// timestampAtOffset seeks to offset, skips to the next newline (to avoid
// landing mid-line), then parses the timestamp of the following line.
func timestampAtOffset(rs io.ReadSeeker, offset int64) (time.Time, error) {
	if _, err := rs.Seek(offset, io.SeekStart); err != nil {
		return time.Time{}, err
	}

	scanner := bufio.NewScanner(rs)

	// If offset > 0, the first scan advances past the partial line we landed on.
	if offset > 0 {
		if !scanner.Scan() {
			return time.Time{}, io.EOF
		}
	}

	if !scanner.Scan() {
		return time.Time{}, io.EOF
	}

	p := parser.New()
	line, err := p.Parse(scanner.Text())
	if err != nil {
		return time.Time{}, err
	}
	return line.Timestamp, nil
}

// snapToLineStart moves backward from offset until a newline (or BOF) is
// found, returning the offset of the first byte of that line.
func snapToLineStart(rs io.ReadSeeker, offset int64) (int64, error) {
	if offset == 0 {
		return 0, nil
	}

	// Walk backward one byte at a time to find the preceding newline.
	for offset > 0 {
		offset--
		if _, err := rs.Seek(offset, io.SeekStart); err != nil {
			return 0, err
		}
		buf := make([]byte, 1)
		if _, err := rs.Read(buf); err != nil {
			return 0, err
		}
		if buf[0] == '\n' {
			return offset + 1, nil
		}
	}
	return 0, nil
}
