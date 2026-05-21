package output

import (
	"bufio"
	"fmt"
	"io"

	"github.com/example/logslice/internal/parser"
)

// Format controls how matched log lines are written to the output.
type Format int

const (
	// FormatRaw writes the original log line as-is.
	FormatRaw Format = iota
	// FormatJSON writes each log line as a JSON object.
	FormatJSON
)

// Writer wraps an io.Writer and writes filtered log lines in the chosen format.
type Writer struct {
	bw     *bufio.Writer
	format Format
	count  int
}

// New creates a new Writer that writes to w using the given Format.
func New(w io.Writer, format Format) *Writer {
	return &Writer{
		bw:     bufio.NewWriter(w),
		format: format,
	}
}

// WriteLine writes a single parsed log line to the underlying writer.
func (w *Writer) WriteLine(line *parser.LogLine) error {
	if line == nil {
		return nil
	}
	var err error
	switch w.format {
	case FormatJSON:
		err = w.writeJSON(line)
	default:
		err = w.writeRaw(line)
	}
	if err != nil {
		return err
	}
	w.count++
	return nil
}

// Flush flushes any buffered data to the underlying writer.
func (w *Writer) Flush() error {
	return w.bw.Flush()
}

// Count returns the number of lines written so far.
func (w *Writer) Count() int {
	return w.count
}

func (w *Writer) writeRaw(line *parser.LogLine) error {
	_, err := fmt.Fprintln(w.bw, line.Raw)
	return err
}

func (w *Writer) writeJSON(line *parser.LogLine) error {
	_, err := fmt.Fprintf(w.bw,
		`{"timestamp":%q,"severity":%q,"message":%q}`+"\n",
		line.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		line.Severity,
		line.Message,
	)
	return err
}
