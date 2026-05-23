package multifile_test

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/multifile"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "log-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, l := range lines {
		fmt.Fprintln(f, l)
	}
	return f.Name()
}

func ts(offset int) string {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return base.Add(time.Duration(offset) * time.Second).Format(time.RFC3339)
}

func TestMultiReader_NoPaths(t *testing.T) {
	_, err := multifile.New(nil)
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestMultiReader_FileNotFound(t *testing.T) {
	_, err := multifile.New([]string{"/nonexistent/path.log"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestMultiReader_SingleFile(t *testing.T) {
	path := writeTempLog(t, []string{
		ts(0) + " INFO  first",
		ts(1) + " INFO  second",
	})
	mr, err := multifile.New([]string{path})
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	line, err := mr.Next()
	if err != nil || line == nil {
		t.Fatalf("expected first line, got err=%v", err)
	}
	if line.Message != "first" {
		t.Errorf("expected 'first', got %q", line.Message)
	}
	line, err = mr.Next()
	if err != nil || line.Message != "second" {
		t.Errorf("expected 'second', got %q / err=%v", line.Message, err)
	}
	_, err = mr.Next()
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}
}

func TestMultiReader_MergesInOrder(t *testing.T) {
	// File A: t=0, t=2, t=4
	pathA := writeTempLog(t, []string{
		ts(0) + " INFO  a0",
		ts(2) + " INFO  a2",
		ts(4) + " INFO  a4",
	})
	// File B: t=1, t=3
	pathB := writeTempLog(t, []string{
		ts(1) + " WARN  b1",
		ts(3) + " WARN  b3",
	})

	mr, err := multifile.New([]string{pathA, pathB})
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	want := []string{"a0", "b1", "a2", "b3", "a4"}
	for i, w := range want {
		line, err := mr.Next()
		if err != nil {
			t.Fatalf("step %d: unexpected error: %v", i, err)
		}
		if line.Message != w {
			t.Errorf("step %d: expected %q, got %q", i, w, line.Message)
		}
	}
	_, err = mr.Next()
	if err != io.EOF {
		t.Errorf("expected EOF after all lines")
	}
}
