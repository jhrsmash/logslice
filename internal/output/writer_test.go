package output_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/example/logslice/internal/output"
	"github.com/example/logslice/internal/parser"
)

func makeLine(ts time.Time, sev, msg, raw string) *parser.LogLine {
	return &parser.LogLine{
		Timestamp: ts,
		Severity:  sev,
		Message:   msg,
		Raw:       raw,
	}
}

func TestWriter_NilLine(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatRaw)
	if err := w.WriteLine(nil); err != nil {
		t.Fatalf("unexpected error on nil line: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for nil line, got %q", buf.String())
	}
}

func TestWriter_FormatRaw(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatRaw)
	ts := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	line := makeLine(ts, "INFO", "server started", "2024-01-15T10:00:00Z INFO server started")

	if err := w.WriteLine(line); err != nil {
		t.Fatalf("WriteLine error: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush error: %v", err)
	}

	got := strings.TrimRight(buf.String(), "\n")
	if got != line.Raw {
		t.Errorf("expected %q, got %q", line.Raw, got)
	}
}

func TestWriter_FormatJSON(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatJSON)
	ts := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	line := makeLine(ts, "ERROR", "disk full", "2024-01-15T10:00:00Z ERROR disk full")

	if err := w.WriteLine(line); err != nil {
		t.Fatalf("WriteLine error: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush error: %v", err)
	}

	got := strings.TrimRight(buf.String(), "\n")
	want := `{"timestamp":"2024-01-15T10:00:00Z","severity":"ERROR","message":"disk full"}`
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestWriter_Count(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.FormatRaw)
	ts := time.Now()

	for i := 0; i < 5; i++ {
		line := makeLine(ts, "DEBUG", "tick", "DEBUG tick")
		if err := w.WriteLine(line); err != nil {
			t.Fatalf("WriteLine error at i=%d: %v", i, err)
		}
	}
	if w.Count() != 5 {
		t.Errorf("expected count=5, got %d", w.Count())
	}
}
