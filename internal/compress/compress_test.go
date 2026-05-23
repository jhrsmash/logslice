package compress_test

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/compress"
)

func writePlain(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "plain-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func writeGzip(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "logs.gz")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	if _, err := gw.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := gw.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestDetector_PlainFile(t *testing.T) {
	const body = "hello plain\n"
	path := writePlain(t, body)

	d := compress.New()
	rc, err := d.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer rc.Close()

	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != body {
		t.Errorf("got %q, want %q", got, body)
	}
}

func TestDetector_GzipByExtension(t *testing.T) {
	const body = "hello gzip\n"
	path := writeGzip(t, body)

	d := compress.New()
	rc, err := d.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer rc.Close()

	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != body {
		t.Errorf("got %q, want %q", got, body)
	}
}

func TestDetector_GzipByMagicBytes(t *testing.T) {
	const body = "magic bytes detection\n"
	// Write a gzip file but name it without .gz extension.
	path := filepath.Join(t.TempDir(), "logs.log")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	gw := gzip.NewWriter(f)
	_, _ = gw.Write([]byte(body))
	_ = gw.Close()
	f.Close()

	d := compress.New()
	rc, err := d.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer rc.Close()

	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != body {
		t.Errorf("got %q, want %q", got, body)
	}
}

func TestDetector_FileNotFound(t *testing.T) {
	d := compress.New()
	_, err := d.Open("/nonexistent/path/file.log")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
