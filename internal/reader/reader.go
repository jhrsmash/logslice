// Package reader provides buffered log file reading with binary search
// capabilities for efficient time-based seeking.
package reader

import (
	"bufio"
	"io"
	"os"
)

// Reader wraps a file and provides line-by-line reading.
type Reader struct {
	f       *os.File
	scanner *bufio.Scanner
	path    string
}

// New opens the file at path and returns a Reader ready to scan lines.
func New(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 64*1024)
	return &Reader{
		f:       f,
		scanner: scanner,
		path:    path,
	}, nil
}

// ReadLine returns the next line from the file, or io.EOF when exhausted.
func (r *Reader) ReadLine() (string, error) {
	if r.scanner.Scan() {
		return r.scanner.Text(), nil
	}
	if err := r.scanner.Err(); err != nil {
		return "", err
	}
	return "", io.EOF
}

// SeekOffset repositions the underlying file to the given byte offset
// and resets the scanner so subsequent ReadLine calls read from that position.
func (r *Reader) SeekOffset(offset int64) error {
	_, err := r.f.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	r.scanner = bufio.NewScanner(r.f)
	r.scanner.Buffer(make([]byte, 64*1024), 64*1024)
	return nil
}

// Size returns the total byte size of the underlying file.
func (r *Reader) Size() (int64, error) {
	info, err := r.f.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// Close releases the file handle.
func (r *Reader) Close() error {
	return r.f.Close()
}
