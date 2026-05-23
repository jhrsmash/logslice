package grep_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/grep"
	"github.com/user/logslice/internal/parser"
)

func makeLine(raw string) *parser.LogLine {
	return &parser.LogLine{
		Timestamp: time.Now(),
		Severity:  "INFO",
		Message:   raw,
		Raw:       raw,
	}
}

func TestMatcher_EmptyPattern_MatchesAll(t *testing.T) {
	m, err := grep.New("", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match(makeLine("anything here")) {
		t.Error("expected match for empty pattern")
	}
}

func TestMatcher_NilLine_ReturnsTrue(t *testing.T) {
	m, _ := grep.New("foo", false)
	if !m.Match(nil) {
		t.Error("expected true for nil line")
	}
}

func TestMatcher_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := grep.New("[invalid", false)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestMatcher_Match_Hit(t *testing.T) {
	m, err := grep.New(`error`, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match(makeLine("2024-01-01 ERROR something went wrong")) {
		t.Error("expected pattern to match")
	}
}

func TestMatcher_Match_Miss(t *testing.T) {
	m, _ := grep.New(`panic`, false)
	if m.Match(makeLine("INFO all good")) {
		t.Error("expected no match")
	}
}

func TestMatcher_Invert_HitBecomesmiss(t *testing.T) {
	m, _ := grep.New(`debug`, true)
	if m.Match(makeLine("DEBUG verbose output")) {
		t.Error("invert: expected no match when pattern hits")
	}
}

func TestMatcher_Invert_MissBecomeHit(t *testing.T) {
	m, _ := grep.New(`debug`, true)
	if !m.Match(makeLine("INFO normal line")) {
		t.Error("invert: expected match when pattern misses")
	}
}

func TestMatcher_CaseInsensitive(t *testing.T) {
	m, _ := grep.New(`(?i)warn`, false)
	if !m.Match(makeLine("WARN disk space low")) {
		t.Error("expected case-insensitive match")
	}
}

func TestMatcher_String_NoPattern(t *testing.T) {
	m, _ := grep.New("", false)
	if m.String() != "<no pattern>" {
		t.Errorf("unexpected String(): %s", m.String())
	}
}

func TestMatcher_String_WithPattern(t *testing.T) {
	m, _ := grep.New(`foo`, false)
	if m.String() != "/foo/" {
		t.Errorf("unexpected String(): %s", m.String())
	}
}

func TestMatcher_String_Invert(t *testing.T) {
	m, _ := grep.New(`bar`, true)
	if m.String() != "NOT /bar/" {
		t.Errorf("unexpected String(): %s", m.String())
	}
}
