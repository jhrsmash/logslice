package reader_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/user/logslice/internal/reader"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestReader_ReadAllLines(t *testing.T) {
	lines := []string{
		"2024-01-01T00:00:00Z INFO  hello",
		"2024-01-01T00:00:01Z DEBUG world",
		"2024-01-01T00:00:02Z ERROR boom",
	}
	path := writeTempFile(t, strings.Join(lines, "\n")+"\n")

	r, err := reader.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	for i, want := range lines {
		got, err := r.ReadLine()
		if err != nil {
			t.Fatalf("line %d: unexpected error: %v", i, err)
		}
		if got != want {
			t.Errorf("line %d: got %q, want %q", i, got, want)
		}
	}
	_, err = r.ReadLine()
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}
}

func TestReader_SeekOffset(t *testing.T) {
	content := "AAAA\nBBBB\nCCCC\n"
	path := writeTempFile(t, content)

	r, err := reader.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	// Seek to second line (offset 5)
	if err := r.SeekOffset(5); err != nil {
		t.Fatalf("SeekOffset: %v", err)
	}
	got, err := r.ReadLine()
	if err != nil {
		t.Fatalf("ReadLine: %v", err)
	}
	if got != "BBBB" {
		t.Errorf("got %q, want %q", got, "BBBB")
	}
}

func TestReader_Size(t *testing.T) {
	content := "hello world"
	path := writeTempFile(t, content)

	r, err := reader.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	size, err := r.Size()
	if err != nil {
		t.Fatalf("Size: %v", err)
	}
	if size != int64(len(content)) {
		t.Errorf("got size %d, want %d", size, len(content))
	}
}

func TestReader_NotFound(t *testing.T) {
	_, err := reader.New("/nonexistent/path/to/file.log")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
