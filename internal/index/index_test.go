package index_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/index"
	"github.com/user/logslice/internal/parser"
)

func baseTime() time.Time {
	t, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	return t
}

func TestIndex_FloorOffset_Empty(t *testing.T) {
	idx := index.New()
	_, ok := idx.FloorOffset(baseTime())
	if ok {
		t.Fatal("expected false for empty index")
	}
}

func TestIndex_FloorOffset_BeforeAll(t *testing.T) {
	idx := index.New()
	t0 := baseTime()
	idx.Add(t0.Add(time.Hour), 100)
	idx.Add(t0.Add(2*time.Hour), 200)

	_, ok := idx.FloorOffset(t0) // t0 is before every entry
	if ok {
		t.Fatal("expected false when target is before all entries")
	}
}

func TestIndex_FloorOffset_Exact(t *testing.T) {
	idx := index.New()
	t0 := baseTime()
	idx.Add(t0, 0)
	idx.Add(t0.Add(time.Hour), 512)

	off, ok := idx.FloorOffset(t0)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if off != 0 {
		t.Fatalf("expected offset 0, got %d", off)
	}
}

func TestIndex_FloorOffset_Mid(t *testing.T) {
	idx := index.New()
	t0 := baseTime()
	idx.Add(t0, 0)
	idx.Add(t0.Add(time.Hour), 1000)
	idx.Add(t0.Add(2*time.Hour), 2000)

	off, ok := idx.FloorOffset(t0.Add(90 * time.Minute))
	if !ok {
		t.Fatal("expected ok=true")
	}
	if off != 1000 {
		t.Fatalf("expected offset 1000, got %d", off)
	}
}

func TestBuild_SamplesEntries(t *testing.T) {
	var sb strings.Builder
	t0 := baseTime()
	for i := 0; i < 10; i++ {
		ts := t0.Add(time.Duration(i) * time.Minute).Format(time.RFC3339)
		sb.WriteString(ts + " INFO message\n")
	}

	p := parser.New()
	r := strings.NewReader(sb.String())

	idx, err := index.Build(r, p, index.BuildOptions{SampleEvery: 3})
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	if idx.Len() == 0 {
		t.Fatal("expected at least one index entry")
	}
}

func TestIndex_Entries_ReturnsCopy(t *testing.T) {
	idx := index.New()
	idx.Add(baseTime(), 0)
	entries := idx.Entries()
	entries[0].Offset = 9999
	if idx.Entries()[0].Offset == 9999 {
		t.Fatal("Entries should return a copy, not a reference")
	}
}
