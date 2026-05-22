// Package tail provides functionality to read the last N lines of a log file
// efficiently without loading the entire file into memory.
package tail

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/logslice/logslice/internal/parser"
)

const defaultChunkSize = 4096

// Tailer reads the last N lines from a log file.
type Tailer struct {
	filePath  string
	chunkSize int64
	p         *parser.Parser
}

// New returns a new Tailer for the given file path.
func New(filePath string) *Tailer {
	return &Tailer{
		filePath:  filePath,
		chunkSize: defaultChunkSize,
		p:         parser.New(),
	}
}

// LastN returns the last n parsed log lines from the file.
// Lines that cannot be parsed are skipped.
func (t *Tailer) LastN(n int) ([]*parser.LogLine, error) {
	if n <= 0 {
		return nil, fmt.Errorf("tail: n must be greater than zero, got %d", n)
	}

	f, err := os.Open(t.filePath)
	if err != nil {
		return nil, fmt.Errorf("tail: open file: %w", err)
	}
	defer f.Close()

	size, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("tail: seek end: %w", err)
	}

	offset, err := findTailOffset(f, size, n, t.chunkSize)
	if err != nil {
		return nil, err
	}

	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("tail: seek to offset: %w", err)
	}

	var lines []*parser.LogLine
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line, err := t.p.Parse(scanner.Text())
		if err != nil {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("tail: scan: %w", err)
	}

	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return lines, nil
}

// findTailOffset seeks backward through the file to find the byte offset
// that contains at least n newline-terminated lines.
func findTailOffset(f *os.File, size int64, n int, chunkSize int64) (int64, error) {
	newlines := 0
	offset := size
	buf := make([]byte, chunkSize)

	for offset > 0 && newlines <= n {
		read := chunkSize
		if offset < chunkSize {
			read = offset
		}
		offset -= read
		if _, err := f.Seek(offset, io.SeekStart); err != nil {
			return 0, fmt.Errorf("tail: backward seek: %w", err)
		}
		slice := buf[:read]
		if _, err := io.ReadFull(f, slice); err != nil {
			return 0, fmt.Errorf("tail: read chunk: %w", err)
		}
		for i := int64(read) - 1; i >= 0; i-- {
			if slice[i] == '\n' {
				newlines++
				if newlines > n {
					return offset + i + 1, nil
				}
			}
		}
	}
	return 0, nil
}
