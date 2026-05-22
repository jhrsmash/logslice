package reader

import (
	"io"
	"os"
	"time"

	"github.com/example/logslice/internal/parser"
)

// SeekToTime performs a binary search within the file to find the first line
// whose timestamp is >= target. Returns the byte offset of that line, or
// io.EOF if no such line exists.
func SeekToTime(f *os.File, size int64, target time.Time, p *parser.Parser) (int64, error) {
	if size == 0 {
		return 0, io.EOF
	}

	lo, hi := int64(0), size

	for lo < hi {
		mid := (lo + hi) / 2
		offset, err := snapToLineStart(f, mid, size)
		if err != nil {
			return 0, err
		}

		ts, err := timestampAtOffset(f, offset, p)
		if err != nil {
			// Can't parse; move forward
			lo = offset + 1
			continue
		}

		if ts.Before(target) {
			lo = offset + 1
		} else {
			hi = offset
		}
	}

	if lo >= size {
		return 0, io.EOF
	}

	// Snap lo to the start of its line
	final, err := snapToLineStart(f, lo, size)
	if err != nil {
		return 0, err
	}
	return final, nil
}

// timestampAtOffset reads the line at the given offset and parses its timestamp.
func timestampAtOffset(f *os.File, offset int64, p *parser.Parser) (time.Time, error) {
	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return time.Time{}, err
	}

	line, err := p.Parse(f)
	if err != nil {
		return time.Time{}, err
	}
	if line == nil {
		return time.Time{}, io.EOF
	}
	return line.Timestamp, nil
}

// snapToLineStart moves backward from offset until a newline (or BOF) is found,
// returning the offset of the first character of the line.
func snapToLineStart(f *os.File, offset int64, size int64) (int64, error) {
	if offset == 0 {
		return 0, nil
	}

	// Read backward in small chunks to find the preceding newline
	buf := make([]byte, 1)
	pos := offset
	for pos > 0 {
		pos--
		if _, err := f.Seek(pos, io.SeekStart); err != nil {
			return 0, err
		}
		if _, err := f.Read(buf); err != nil {
			return 0, err
		}
		if buf[0] == '\n' {
			return pos + 1, nil
		}
	}
	return 0, nil
}
