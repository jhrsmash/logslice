package schema

import (
	"testing"
)

func makeSchema(t *testing.T, name, pattern string) *Schema {
	t.Helper()
	s, err := New(name, pattern)
	if err != nil {
		t.Fatalf("makeSchema(%q): %v", name, err)
	}
	return s
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	reg := NewRegistry()
	s := makeSchema(t, "app", samplePattern)
	if err := reg.Register(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := reg.Get("app")
	if !ok || got.Name != "app" {
		t.Errorf("expected to retrieve schema 'app'")
	}
}

func TestRegistry_RegisterNil(t *testing.T) {
	reg := NewRegistry()
	if err := reg.Register(nil); err == nil {
		t.Fatal("expected error registering nil schema")
	}
}

func TestRegistry_DuplicateName(t *testing.T) {
	reg := NewRegistry()
	s := makeSchema(t, "app", samplePattern)
	reg.Register(s)
	if err := reg.Register(s); err == nil {
		t.Fatal("expected error for duplicate schema name")
	}
}

func TestRegistry_GetMissing(t *testing.T) {
	reg := NewRegistry()
	_, ok := reg.Get("missing")
	if ok {
		t.Fatal("expected miss for unknown schema name")
	}
}

func TestRegistry_Names(t *testing.T) {
	reg := NewRegistry()
	reg.Register(makeSchema(t, "a", samplePattern))
	reg.Register(makeSchema(t, "b", `(?P<x>.+)`))
	names := reg.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestRegistry_MatchFirst_Hit(t *testing.T) {
	reg := NewRegistry()
	reg.Register(makeSchema(t, "app", samplePattern))
	name, fields := reg.MatchFirst("2024-01-15T10:30:00 WARN low disk")
	if name != "app" {
		t.Errorf("expected schema 'app', got %q", name)
	}
	if fields["level"] != "WARN" {
		t.Errorf("expected level=WARN, got %q", fields["level"])
	}
}

func TestRegistry_MatchFirst_Miss(t *testing.T) {
	reg := NewRegistry()
	reg.Register(makeSchema(t, "app", samplePattern))
	name, fields := reg.MatchFirst("no match here")
	if name != "" || fields != nil {
		t.Errorf("expected no match, got name=%q fields=%v", name, fields)
	}
}
