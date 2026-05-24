package schema

import (
	"testing"
)

const samplePattern = `(?P<ts>\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}) (?P<level>\w+) (?P<msg>.+)`

func TestNew_ValidPattern(t *testing.T) {
	s, err := New("test", samplePattern)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name != "test" {
		t.Errorf("expected name 'test', got %q", s.Name)
	}
	if len(s.Fields()) != 3 {
		t.Errorf("expected 3 fields, got %d", len(s.Fields()))
	}
}

func TestNew_EmptyName(t *testing.T) {
	_, err := New("", samplePattern)
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New("bad", "(?P<x>[")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_NoNamedGroups(t *testing.T) {
	_, err := New("nogroups", `\d+ \w+`)
	if err == nil {
		t.Fatal("expected error when no named groups")
	}
}

func TestSchema_Match_Hit(t *testing.T) {
	s, _ := New("test", samplePattern)
	fields := s.Match("2024-01-15T10:30:00 ERROR something went wrong")
	if fields == nil {
		t.Fatal("expected match, got nil")
	}
	if fields["level"] != "ERROR" {
		t.Errorf("expected level=ERROR, got %q", fields["level"])
	}
	if fields["msg"] != "something went wrong" {
		t.Errorf("unexpected msg: %q", fields["msg"])
	}
}

func TestSchema_Match_Miss(t *testing.T) {
	s, _ := New("test", samplePattern)
	if fields := s.Match("not a log line"); fields != nil {
		t.Errorf("expected nil, got %v", fields)
	}
}

func TestSchema_Match_TrailingNewline(t *testing.T) {
	s, _ := New("test", samplePattern)
	fields := s.Match("2024-01-15T10:30:00 INFO hello\n")
	if fields == nil {
		t.Fatal("expected match with trailing newline")
	}
}
