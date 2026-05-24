package rewind_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/rewind"
)

func writeTempLog(t *testing.T, lines []string) (path string, size int64) {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "rewind-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, l := range lines {
		fmt.Fprintln(f, l)
	}
	info, _ := f.Stat()
	return f.Name(), info.Size()
}

func makeLines(n int) []string {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = fmt.Sprintf("%s INFO message number %d",
			base.Add(time.Duration(i)*time.Second).Format(time.RFC3339), i)
	}
	return out
}

func TestRewinder_NilReader(t *testing.T) {
	_, err := rewind.New(nil, parser.New(), 0)
	if err == nil {
		t.Fatal("expected error for nil reader")
	}
}

func TestRewinder_NilParser(t *testing.T) {
	_, err := rewind.New(strings.NewReader(""), nil, 0)
	if err == nil {
		t.Fatal("expected error for nil parser")
	}
}

func TestRewinder_EmptyFile(t *testing.T) {
	rw, err := rewind.New(bytes.NewReader(nil), parser.New(), 512)
	if err != nil {
		t.Fatal(err)
	}
	lines, err := rw.Last(0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(lines))
	}
}

func TestRewinder_LastN_Basic(t *testing.T) {
	rawLines := makeLines(10)
	path, size := writeTempLog(t, rawLines)

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	rw, err := rewind.New(f, parser.New(), 512)
	if err != nil {
		t.Fatal(err)
	}

	got, err := rw.Last(size, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	// oldest of the last 3 is index 7
	if !strings.Contains(got[0].Raw, "number 7") {
		t.Errorf("unexpected first line: %q", got[0].Raw)
	}
}

func TestRewinder_LastN_MoreThanAvailable(t *testing.T) {
	rawLines := makeLines(5)
	path, size := writeTempLog(t, rawLines)

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	rw, err := rewind.New(f, parser.New(), 64)
	if err != nil {
		t.Fatal(err)
	}

	got, err := rw.Last(size, 20)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(got))
	}
}

func TestRewinder_ZeroN(t *testing.T) {
	rw, err := rewind.New(io.NopCloser(strings.NewReader("")), parser.New(), 0)
	if err != nil {
		t.Fatal(err)
	}
	lines, err := rw.Last(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if lines != nil {
		t.Fatal("expected nil slice for n=0")
	}
}
