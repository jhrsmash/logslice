package reader_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/reader"
)

func makeLogFile(t *testing.T, entries []string) string {
	t.Helper()
	return writeTempFile(t, strings.Join(entries, "\n")+"\n")
}

func TestSeekToTime_BeginningOfFile(t *testing.T) {
	entries := []string{
		"2024-03-01T10:00:00Z INFO  first",
		"2024-03-01T10:00:01Z INFO  second",
		"2024-03-01T10:00:02Z INFO  third",
	}
	path := makeLogFile(t, entries)

	r, err := reader.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	p := parser.New()
	target := time.Date(2024, 3, 1, 9, 59, 0, 0, time.UTC)

	offset, err := reader.SeekToTime(r, p, target)
	if err != nil {
		t.Fatalf("SeekToTime: %v", err)
	}
	// Target is before all entries — should return offset 0
	if offset != 0 {
		t.Errorf("expected offset 0, got %d", offset)
	}
}

func TestSeekToTime_EmptyFile(t *testing.T) {
	path := writeTempFile(t, "")

	r, err := reader.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	p := parser.New()
	target := time.Now()

	offset, err := reader.SeekToTime(r, p, target)
	if err != nil {
		t.Fatalf("SeekToTime: %v", err)
	}
	if offset != 0 {
		t.Errorf("expected offset 0 for empty file, got %d", offset)
	}
}
