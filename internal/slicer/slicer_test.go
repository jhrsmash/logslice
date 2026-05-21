package slicer_test

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/slicer"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f.Name()
}

func TestSlicer_NoFilter(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z INFO  app started",
		"2024-01-01T10:01:00Z WARN  disk usage high",
		"2024-01-01T10:02:00Z ERROR connection refused",
	}
	path := writeTempLog(t, lines)

	f, _ := filter.New(filter.Options{})
	var buf bytes.Buffer
	s := slicer.New(slicer.Options{FilePath: path, Filter: f, Writer: &buf})

	n, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 lines written, got %d", n)
	}
}

func TestSlicer_SeverityFilter(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z INFO  app started",
		"2024-01-01T10:01:00Z WARN  disk usage high",
		"2024-01-01T10:02:00Z ERROR connection refused",
	}
	path := writeTempLog(t, lines)

	f, _ := filter.New(filter.Options{MinSeverity: parser.SeverityWarn})
	var buf bytes.Buffer
	s := slicer.New(slicer.Options{FilePath: path, Filter: f, Writer: &buf})

	n, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 lines written, got %d", n)
	}
	if strings.Contains(buf.String(), "INFO") {
		t.Error("INFO line should have been filtered out")
	}
}

func TestSlicer_TimeRangeFilter(t *testing.T) {
	lines := []string{
		"2024-01-01T09:00:00Z INFO  before range",
		"2024-01-01T10:00:00Z INFO  in range",
		"2024-01-01T11:00:00Z INFO  after range",
	}
	path := writeTempLog(t, lines)

	start := time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC)
	f, _ := filter.New(filter.Options{Start: &start, End: &end})
	var buf bytes.Buffer
	s := slicer.New(slicer.Options{FilePath: path, Filter: f, Writer: &buf})

	n, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 line written, got %d", n)
	}
}

func TestSlicer_FileNotFound(t *testing.T) {
	f, _ := filter.New(filter.Options{})
	s := slicer.New(slicer.Options{FilePath: "/nonexistent/file.log", Filter: f, Writer: &bytes.Buffer{}})
	_, err := s.Run()
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
