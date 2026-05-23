package reader

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/example/logslice/internal/parser"
)

func makeLogFile(t *testing.T, lines []string) *os.File {
	t.Helper()
	f, err := os.CreateTemp("", "seek_test_*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	t.Cleanup(func() {
		f.Close()
		os.Remove(f.Name())
	})
	for _, l := range lines {
		fmt.Fprintln(f, l)
	}
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatalf("seek: %v", err)
	}
	return f
}

// seekAndParse seeks the file to offset and parses the next log line.
// It is a test helper to reduce repetition in seek result assertions.
func seekAndParse(t *testing.T, f *os.File, offset int64, p *parser.Parser) *parser.Entry {
	t.Helper()
	if _, err := f.Seek(offset, 0); err != nil {
		t.Fatalf("seek to offset %d: %v", offset, err)
	}
	line, err := p.Parse(f)
	if err != nil {
		t.Fatalf("parse at offset %d: %v", offset, err)
	}
	return line
}

func TestSeekToTime_BeginningOfFile(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z INFO  startup complete",
		"2024-01-01T10:01:00Z DEBUG checking config",
		"2024-01-01T10:02:00Z WARN  disk usage high",
		"2024-01-01T10:03:00Z ERROR disk full",
	}

	f := makeLogFile(t, lines)
	stat, _ := f.Stat()
	p := parser.New()

	target := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	offset, err := SeekToTime(f, stat.Size(), target, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if offset != 0 {
		t.Errorf("expected offset 0, got %d", offset)
	}
}

func TestSeekToTime_EmptyFile(t *testing.T) {
	f := makeLogFile(t, []string{})
	p := parser.New()

	target := time.Now()
	_, err := SeekToTime(f, 0, target, p)
	if err == nil {
		t.Error("expected io.EOF for empty file, got nil")
	}
}

func TestSeekToTime_MidFile(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z INFO  startup complete",
		"2024-01-01T10:01:00Z DEBUG checking config",
		"2024-01-01T10:02:00Z WARN  disk usage high",
		"2024-01-01T10:03:00Z ERROR disk full",
	}

	f := makeLogFile(t, lines)
	stat, _ := f.Stat()
	p := parser.New()

	target := time.Date(2024, 1, 1, 10, 2, 0, 0, time.UTC)
	offset, err := SeekToTime(f, stat.Size(), target, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := seekAndParse(t, f, offset, p)
	if !line.Timestamp.Equal(target) {
		t.Errorf("expected timestamp %v, got %v", target, line.Timestamp)
	}
}

func TestSeekToTime_AfterLastEntry(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z INFO  startup complete",
		"2024-01-01T10:01:00Z DEBUG checking config",
		"2024-01-01T10:02:00Z WARN  disk usage high",
	}

	f := makeLogFile(t, lines)
	stat, _ := f.Stat()
	p := parser.New()

	// Target is after all log entries; expect the last entry to be returned.
	target := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	offset, err := SeekToTime(f, stat.Size(), target, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := seekAndParse(t, f, offset, p)
	lastEntryTime := time.Date(2024, 1, 1, 10, 2, 0, 0, time.UTC)
	if !line.Timestamp.Equal(lastEntryTime) {
		t.Errorf("expected last entry timestamp %v, got %v", lastEntryTime, line.Timestamp)
	}
}
