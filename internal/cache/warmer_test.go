package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	for _, l := range lines {
		fmt.Fprintln(f, l)
	}
	return f.Name()
}

func TestWarmer_WarmBuildsIndex(t *testing.T) {
	lines := []string{
		"2024-01-15T10:00:00Z INFO  startup complete",
		"2024-01-15T10:01:00Z DEBUG checking config",
		"2024-01-15T10:02:00Z ERROR disk full",
	}
	path := writeTempLog(t, lines)

	c := New(5 * time.Minute)
	w := NewWarmer(c)

	idx, err := w.Warm(path, time.Minute)
	if err != nil {
		t.Fatalf("Warm: %v", err)
	}
	if idx == nil {
		t.Fatal("expected non-nil index")
	}
	if c.Len() != 1 {
		t.Errorf("expected 1 cache entry, got %d", c.Len())
	}
}

func TestWarmer_WarmReturnsCachedEntry(t *testing.T) {
	lines := []string{
		"2024-01-15T10:00:00Z INFO  hello",
	}
	path := writeTempLog(t, lines)

	c := New(5 * time.Minute)
	w := NewWarmer(c)

	idx1, err := w.Warm(path, time.Minute)
	if err != nil {
		t.Fatalf("first Warm: %v", err)
	}

	idx2, err := w.Warm(path, time.Minute)
	if err != nil {
		t.Fatalf("second Warm: %v", err)
	}

	if idx1 != idx2 {
		t.Error("expected the same index pointer on cache hit")
	}
}

func TestWarmer_WarmFileNotFound(t *testing.T) {
	c := New(5 * time.Minute)
	w := NewWarmer(c)

	missing := filepath.Join(t.TempDir(), "no-such-file.log")
	_, err := w.Warm(missing, time.Minute)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
