package rewind_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/rewind"
)

// TestRewinder_ChronologicalOrder verifies that lines come back oldest-first
// regardless of chunk boundary alignment.
func TestRewinder_ChronologicalOrder(t *testing.T) {
	const total = 20
	base := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	raw := make([]string, total)
	for i := 0; i < total; i++ {
		raw[i] = fmt.Sprintf("%s WARN event %02d",
			base.Add(time.Duration(i)*time.Minute).Format(time.RFC3339), i)
	}

	path, size := writeTempLog(t, raw)
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Use a tiny chunk to stress boundary handling
	rw, err := rewind.New(f, parser.New(), 32)
	if err != nil {
		t.Fatal(err)
	}

	got, err := rw.Last(size, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(got))
	}

	for i := 1; i < len(got); i++ {
		if !got[i].Timestamp.After(got[i-1].Timestamp) {
			t.Errorf("line %d timestamp not after line %d: %v <= %v",
				i, i-1, got[i].Timestamp, got[i-1].Timestamp)
		}
	}

	// last line should be event 19
	if !strings.Contains(got[len(got)-1].Raw, "event 19") {
		t.Errorf("expected last event to be 19, got: %q", got[len(got)-1].Raw)
	}
}

// TestRewinder_SmallChunkEqualsLargeChunk ensures chunk size doesn't affect results.
func TestRewinder_SmallChunkEqualsLargeChunk(t *testing.T) {
	const total = 15
	base := time.Date(2024, 3, 1, 8, 0, 0, 0, time.UTC)
	raw := make([]string, total)
	for i := 0; i < total; i++ {
		raw[i] = fmt.Sprintf("%s ERROR fault %d",
			base.Add(time.Duration(i)*time.Second).Format(time.RFC3339), i)
	}

	path, size := writeTempLog(t, raw)

	open := func() *os.File {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		return f
	}

	f1 := open()
	defer f1.Close()
	rw1, _ := rewind.New(f1, parser.New(), 16)
	small, _ := rw1.Last(size, 6)

	f2 := open()
	defer f2.Close()
	rw2, _ := rewind.New(f2, parser.New(), 65536)
	large, _ := rw2.Last(size, 6)

	if len(small) != len(large) {
		t.Fatalf("length mismatch: small=%d large=%d", len(small), len(large))
	}
	for i := range small {
		if small[i].Raw != large[i].Raw {
			t.Errorf("line %d mismatch:\n  small: %q\n  large: %q",
				i, small[i].Raw, large[i].Raw)
		}
	}
}
