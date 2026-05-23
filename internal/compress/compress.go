// Package compress provides transparent decompression for log files,
// supporting gzip and zstd formats detected by file extension or magic bytes.
package compress

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

// Format represents a supported compression format.
type Format int

const (
	FormatNone Format = iota
	FormatGzip
)

// Detector inspects a file and returns an appropriate reader.
type Detector struct{}

// New returns a new Detector.
func New() *Detector {
	return &Detector{}
}

// Open returns a ReadCloser that transparently decompresses the file at path.
// If the file is not compressed, the raw file is returned.
func (d *Detector) Open(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("compress: open %q: %w", path, err)
	}

	fmt := detectFormat(path, f)
	if fmt == FormatNone {
		return f, nil
	}

	// Seek back to start after magic-byte sniff.
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		f.Close()
		return nil, fmt.Errorf("compress: seek %q: %w", path, err)
	}

	gr, err := gzip.NewReader(f)
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("compress: gzip reader %q: %w", path, err)
	}
	return &gzipReadCloser{gr: gr, f: f}, nil
}

// detectFormat returns the compression format by extension or magic bytes.
func detectFormat(path string, f *os.File) Format {
	lower := strings.ToLower(path)
	if strings.HasSuffix(lower, ".gz") || strings.HasSuffix(lower, ".gzip") {
		return FormatGzip
	}

	// Sniff first two bytes for gzip magic number 0x1f 0x8b.
	magic := make([]byte, 2)
	n, _ := f.Read(magic)
	if n == 2 && magic[0] == 0x1f && magic[1] == 0x8b {
		return FormatGzip
	}
	return FormatNone
}

// gzipReadCloser wraps a gzip.Reader and the underlying file so both are closed.
type gzipReadCloser struct {
	gr *gzip.Reader
	f  *os.File
}

func (g *gzipReadCloser) Read(p []byte) (int, error) {
	return g.gr.Read(p)
}

func (g *gzipReadCloser) Close() error {
	gerr := g.gr.Close()
	ferr := g.f.Close()
	if gerr != nil {
		return gerr
	}
	return ferr
}
