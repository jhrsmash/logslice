package reader

import (
	"os"
	"strings"
	"testing"
	"time"
)

func makeLogFile(t *testing.T, lines []string) *os.File {
	t.Helper()
	f, err := os.CreateTemp("", "logslice-seek-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })

	if _, err := f.WriteString(strings.Join(lines, "\n") + "\n"); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatalf("seek temp file: %v", err)
	}
	return f
}

func TestSeekToTime_BeginningOfFile(t *testing.T) {
	lines := []string{
		"2024-01-01T10:00:00Z INFO  first message",
		"2024-01-01T10:01:00Z INFO  second message",
		"2024-01-01T10:02:00Z WARN  third message",
	}
	f := makeLogFile(t, lines)
	defer f.Close()

	info, _ := f.Stat()
	target := time.Date(2024, 1, 1, 9, 59, 0, 0, time.UTC)

	offset, err := SeekToTime(f, info.Size(), target)
	if err != nil {
		t.Fatalf("SeekToTime error: %v", err)
	}
	if offset != 0 {
		t.Errorf("expected offset 0 for target before all lines, got %d", offset)
	}
}

func TestSeekToTime_EmptyFile(t *testing.T) {
	f, err := os.CreateTemp("", "logslice-empty-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	target := time.Now()
	offset, err := SeekToTime(f, 0, target)
	if err != nil {
		t.Fatalf("SeekToTime error on empty file: %v", err)
	}
	if offset != 0 {
		t.Errorf("expected offset 0 for empty file, got %d", offset)
	}
}
