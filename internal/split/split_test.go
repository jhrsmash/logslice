package split_test

import (
	"sort"
	"testing"
	"time"

	"github.com/example/logslice/internal/parser"
	"github.com/example/logslice/internal/split"
)

func makeLine(sev, msg string) *parser.LogLine {
	return &parser.LogLine{
		Timestamp: time.Now(),
		Severity:  sev,
		Message:   msg,
		Raw:       sev + " " + msg,
	}
}

func TestNew_NilKeyFunc_ReturnsError(t *testing.T) {
	_, err := split.New(nil)
	if err == nil {
		t.Fatal("expected error for nil KeyFunc, got nil")
	}
}

func TestNew_ValidKeyFunc_ReturnsNonNil(t *testing.T) {
	s, err := split.New(func(l *parser.LogLine) string { return l.Severity })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Splitter")
	}
}

func TestAdd_NilLine_IsIgnored(t *testing.T) {
	s, _ := split.New(func(l *parser.LogLine) string { return l.Severity })
	s.Add(nil) // must not panic
	if len(s.Keys()) != 0 {
		t.Fatal("expected no buckets after adding nil line")
	}
}

func TestAdd_EmptyKey_IsIgnored(t *testing.T) {
	s, _ := split.New(func(_ *parser.LogLine) string { return "" })
	s.Add(makeLine("INFO", "hello"))
	if len(s.Keys()) != 0 {
		t.Fatal("expected no buckets when key is empty")
	}
}

func TestAdd_PartitionsBySeverity(t *testing.T) {
	s, _ := split.New(func(l *parser.LogLine) string { return l.Severity })
	s.Add(makeLine("INFO", "a"))
	s.Add(makeLine("INFO", "b"))
	s.Add(makeLine("ERROR", "c"))

	info := s.Bucket("INFO")
	if len(info) != 2 {
		t.Fatalf("expected 2 INFO lines, got %d", len(info))
	}
	errLines := s.Bucket("ERROR")
	if len(errLines) != 1 {
		t.Fatalf("expected 1 ERROR line, got %d", len(errLines))
	}
}

func TestBucket_MissingKey_ReturnsNil(t *testing.T) {
	s, _ := split.New(func(l *parser.LogLine) string { return l.Severity })
	if s.Bucket("WARN") != nil {
		t.Fatal("expected nil for missing bucket")
	}
}

func TestBucket_ReturnsCopy(t *testing.T) {
	s, _ := split.New(func(l *parser.LogLine) string { return l.Severity })
	s.Add(makeLine("DEBUG", "x"))
	copy1 := s.Bucket("DEBUG")
	copy1[0] = nil // mutate the copy
	copy2 := s.Bucket("DEBUG")
	if copy2[0] == nil {
		t.Fatal("Bucket should return an independent copy")
	}
}

func TestKeys_ReturnsSortedKeys(t *testing.T) {
	s, _ := split.New(func(l *parser.LogLine) string { return l.Severity })
	s.Add(makeLine("WARN", "w"))
	s.Add(makeLine("INFO", "i"))
	s.Add(makeLine("ERROR", "e"))
	keys := s.Keys()
	sort.Strings(keys)
	expected := []string{"ERROR", "INFO", "WARN"}
	for i, k := range expected {
		if keys[i] != k {
			t.Fatalf("expected key %q at index %d, got %q", k, i, keys[i])
		}
	}
}

func TestReset_ClearsBuckets(t *testing.T) {
	s, _ := split.New(func(l *parser.LogLine) string { return l.Severity })
	s.Add(makeLine("INFO", "a"))
	s.Reset()
	if len(s.Keys()) != 0 {
		t.Fatal("expected empty buckets after Reset")
	}
}
