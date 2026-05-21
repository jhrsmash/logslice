package reader

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// makeLogFile writes a temp log file with lines spaced one minute apart
// starting from baseTime.
func makeLogFile(t *testing.T, baseTime time.Time, count int) string {
	t.Helper()
	f, err := os.CreateTemp("", "seek_test_*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()

	for i := 0; i < count; i++ {
		ts := baseTime.Add(time.Duration(i) * time.Minute)
		line := fmt.Sprintf("%s INFO  message number %d\n", ts.Format(time.RFC3339), i)
		if _, err := f.WriteString(line); err != nil {
			t.Fatalf("write line: %v", err)
		}
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestSeekToTime_BeginningOfFile(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	path := makeLogFile(t, base, 10)

	r, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	// Seek to a time before all entries — should land at offset 0.
	offset, err := SeekToTime(r, base.Add(-time.Hour))
	if err != nil {
		t.Fatalf("SeekToTime: %v", err)
	}
	if offset != 0 {
		t.Errorf("expected offset 0, got %d", offset)
	}
}

func TestSeekToTime_EmptyFile(t *testing.T) {
	f, err := os.CreateTemp("", "empty_seek_*.log")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	r, err := New(f.Name())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	_, err = SeekToTime(r, time.Now())
	if err == nil {
		t.Error("expected EOF error for empty file, got nil")
	}
}
