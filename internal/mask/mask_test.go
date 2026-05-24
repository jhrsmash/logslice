package mask

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeLine(raw string) *parser.LogLine {
	return &parser.LogLine{
		Timestamp: time.Now(),
		Severity:  "INFO",
		Raw:       raw,
	}
}

func TestMasker_NilLine(t *testing.T) {
	m, _ := New(nil)
	if got := m.Mask(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestMasker_NoPatterns_ReturnsOriginal(t *testing.T) {
	m, _ := New(nil)
	line := makeLine("user=alice password=secret")
	got := m.Mask(line)
	if got != line {
		t.Fatal("expected same pointer when no patterns are registered")
	}
}

func TestMasker_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New([]string{`(?P<bad`})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestMasker_MasksPassword(t *testing.T) {
	m, err := New([]string{`password=(\S+)`})
	if err != nil {
		t.Fatal(err)
	}
	line := makeLine("user=alice password=hunter2 host=db")
	got := m.Mask(line)
	if got == line {
		t.Fatal("expected a new copy, not the same pointer")
	}
	if got.Raw == line.Raw {
		t.Fatal("expected raw to be modified")
	}
	if contains(got.Raw, "hunter2") {
		t.Errorf("sensitive value still present: %s", got.Raw)
	}
	if !contains(got.Raw, redacted) {
		t.Errorf("redacted placeholder missing: %s", got.Raw)
	}
}

func TestMasker_MultiplePatterns(t *testing.T) {
	m, err := New([]string{
		`password=(\S+)`,
		`token=(\S+)`,
	})
	if err != nil {
		t.Fatal(err)
	}
	line := makeLine("password=s3cr3t token=abc123 ok=yes")
	got := m.Mask(line)
	if contains(got.Raw, "s3cr3t") || contains(got.Raw, "abc123") {
		t.Errorf("sensitive values still present: %s", got.Raw)
	}
}

func TestMasker_AddPattern(t *testing.T) {
	m, _ := New(nil)
	if err := m.AddPattern(`apikey=(\S+)`); err != nil {
		t.Fatal(err)
	}
	line := makeLine("apikey=TOPSECRET action=read")
	got := m.Mask(line)
	if contains(got.Raw, "TOPSECRET") {
		t.Errorf("sensitive value still present: %s", got.Raw)
	}
}

func TestMasker_NoMatch_ReturnsSamePointer(t *testing.T) {
	m, _ := New([]string{`password=(\S+)`})
	line := makeLine("user=alice action=login")
	got := m.Mask(line)
	if got != line {
		t.Fatal("expected same pointer when no match found")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
