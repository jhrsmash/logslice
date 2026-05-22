package tail_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/tail"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logslice-tail-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	for _, l := range lines {
		fmt.Fprintln(f, l)
	}
	return f.Name()
}

func makeLines(n int) []string {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	lines := make([]string, n)
	for i := 0; i < n; i++ {
		ts := base.Add(time.Duration(i) * time.Minute).Format(time.RFC3339)
		lines[i] = fmt.Sprintf("%s INFO message number %d", ts, i+1)
	}
	return lines
}

func TestTail_LastN_Basic(t *testing.T) {
	lines := makeLines(20)
	path := writeTempLog(t, lines)

	tr := tail.New(path)
	got, err := tr.LastN(5)
	if err != nil {
		t.Fatalf("LastN: %v", err)
	}
	if len(got) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(got))
	}
	if got[0].Raw != lines[15] {
		t.Errorf("first line mismatch: got %q, want %q", got[0].Raw, lines[15])
	}
	if got[4].Raw != lines[19] {
		t.Errorf("last line mismatch: got %q, want %q", got[4].Raw, lines[19])
	}
}

func TestTail_LastN_MoreThanAvailable(t *testing.T) {
	lines := makeLines(3)
	path := writeTempLog(t, lines)

	tr := tail.New(path)
	got, err := tr.LastN(10)
	if err != nil {
		t.Fatalf("LastN: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
}

func TestTail_LastN_EmptyFile(t *testing.T) {
	path := writeTempLog(t, []string{})

	tr := tail.New(path)
	got, err := tr.LastN(5)
	if err != nil {
		t.Fatalf("LastN on empty file: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(got))
	}
}

func TestTail_LastN_InvalidN(t *testing.T) {
	path := writeTempLog(t, makeLines(5))

	tr := tail.New(path)
	_, err := tr.LastN(0)
	if err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
}

func TestTail_LastN_FileNotFound(t *testing.T) {
	tr := tail.New("/nonexistent/path/to/file.log")
	_, err := tr.LastN(5)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
