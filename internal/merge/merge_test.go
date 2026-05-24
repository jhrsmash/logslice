package merge_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/merge"
	"github.com/yourorg/logslice/internal/parser"
)

var base = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func makeSource(lines []*parser.LogLine) merge.Source {
	i := 0
	return func() *parser.LogLine {
		if i >= len(lines) {
			return nil
		}
		l := lines[i]
		i++
		return l
	}
}

func line(offsetSec int, msg string) *parser.LogLine {
	return &parser.LogLine{
		Timestamp: base.Add(time.Duration(offsetSec) * time.Second),
		Raw:       msg,
	}
}

func collect(m *merge.Merger) []*parser.LogLine {
	var out []*parser.LogLine
	for {
		l := m.Next()
		if l == nil {
			break
		}
		out = append(out, l)
	}
	return out
}

func TestMerger_NoSources(t *testing.T) {
	m := merge.New()
	if got := m.Next(); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestMerger_NilSourcesIgnored(t *testing.T) {
	m := merge.New(nil, nil)
	if got := m.Next(); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestMerger_SingleSource(t *testing.T) {
	src := makeSource([]*parser.LogLine{line(1, "a"), line(3, "b"), line(5, "c")})
	m := merge.New(src)
	got := collect(m)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	for i, want := range []string{"a", "b", "c"} {
		if got[i].Raw != want {
			t.Errorf("line %d: want %q, got %q", i, want, got[i].Raw)
		}
	}
}

func TestMerger_TwoSourcesInterleaved(t *testing.T) {
	src1 := makeSource([]*parser.LogLine{line(1, "a"), line(4, "c"), line(6, "e")})
	src2 := makeSource([]*parser.LogLine{line(2, "b"), line(3, "bb"), line(5, "d")})
	m := merge.New(src1, src2)
	got := collect(m)
	want := []string{"a", "b", "bb", "c", "d", "e"}
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i].Raw != w {
			t.Errorf("line %d: want %q, got %q", i, w, got[i].Raw)
		}
	}
}

func TestMerger_OneExhaustedEarly(t *testing.T) {
	src1 := makeSource([]*parser.LogLine{line(1, "only")})
	src2 := makeSource([]*parser.LogLine{line(2, "x"), line(3, "y")})
	m := merge.New(src1, src2)
	got := collect(m)
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
	if got[0].Raw != "only" {
		t.Errorf("first line should be 'only', got %q", got[0].Raw)
	}
}
