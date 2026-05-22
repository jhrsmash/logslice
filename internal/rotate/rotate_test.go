package rotate

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "app.log")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempLog: %v", err)
	}
	return p
}

func TestNew_FileNotFound(t *testing.T) {
	_, err := New("/nonexistent/path/app.log", time.Second)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestPoll_NoRotation(t *testing.T) {
	p := writeTempLog(t, "line1\nline2\n")
	w, err := New(p, time.Second)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := w.Poll(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestPoll_DetectsRotation_NewInode(t *testing.T) {
	p := writeTempLog(t, "original content\n")
	w, err := New(p, time.Second)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Simulate rotation: remove and recreate the file.
	if err := os.Remove(p); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if err := os.WriteFile(p, []byte("new content\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := w.Poll(); !errors.Is(err, ErrRotated) {
		t.Fatalf("expected ErrRotated, got %v", err)
	}
}

func TestPoll_DetectsTruncation(t *testing.T) {
	p := writeTempLog(t, "some longer content here\n")
	w, err := New(p, time.Second)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Truncate in place (same inode, smaller size).
	if err := os.WriteFile(p, []byte("short\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := w.Poll(); !errors.Is(err, ErrRotated) {
		t.Fatalf("expected ErrRotated after truncation, got %v", err)
	}
}

func TestInterval(t *testing.T) {
	p := writeTempLog(t, "data\n")
	w, err := New(p, 5*time.Second)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if got := w.Interval(); got != 5*time.Second {
		t.Fatalf("Interval: want 5s, got %v", got)
	}
}
