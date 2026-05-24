package highlight_test

import (
	"strings"
	"testing"

	"github.com/logslice/logslice/internal/highlight"
	"github.com/logslice/logslice/internal/parser"
)

func TestHighlighter_Disabled_ReturnsRaw(t *testing.T) {
	h := highlight.New(false)
	raw := "2024-01-01T00:00:00Z ERROR something broke"
	got := h.Format(raw, parser.ParseSeverity("ERROR"))
	if got != raw {
		t.Fatalf("expected raw line, got %q", got)
	}
}

func TestHighlighter_Enabled_ContainsRaw(t *testing.T) {
	h := highlight.New(true)
	raw := "2024-01-01T00:00:00Z INFO hello world"
	got := h.Format(raw, parser.ParseSeverity("INFO"))
	if !strings.Contains(got, raw) {
		t.Fatalf("formatted output should contain original line; got %q", got)
	}
}

func TestHighlighter_Enabled_ContainsResetCode(t *testing.T) {
	h := highlight.New(true)
	raw := "line"
	got := h.Format(raw, parser.ParseSeverity("WARN"))
	if !strings.Contains(got, "\033[") {
		t.Fatalf("expected ANSI escape in output; got %q", got)
	}
}

func TestHighlighter_FormatSeverity_Disabled(t *testing.T) {
	h := highlight.New(false)
	sev := parser.ParseSeverity("ERROR")
	got := h.FormatSeverity(sev)
	if got != sev.String() {
		t.Fatalf("expected %q, got %q", sev.String(), got)
	}
}

func TestHighlighter_FormatSeverity_Enabled(t *testing.T) {
	h := highlight.New(true)
	sev := parser.ParseSeverity("FATAL")
	got := h.FormatSeverity(sev)
	if !strings.Contains(got, sev.String()) {
		t.Fatalf("severity label missing from %q", got)
	}
	if !strings.Contains(got, "\033[") {
		t.Fatalf("expected ANSI escape in severity label; got %q", got)
	}
}

func TestHighlighter_SeverityColors_Distinct(t *testing.T) {
	h := highlight.New(true)
	raw := "line"
	errorFmt := h.Format(raw, parser.ParseSeverity("ERROR"))
	infoFmt := h.Format(raw, parser.ParseSeverity("INFO"))
	warnFmt := h.Format(raw, parser.ParseSeverity("WARN"))
	if errorFmt == infoFmt || infoFmt == warnFmt || errorFmt == warnFmt {
		t.Fatal("expected distinct color codes for different severities")
	}
}

// TestHighlighter_Enabled_EndsWithReset verifies that formatted output ends
// with an ANSI reset sequence so that colors do not bleed into subsequent output.
func TestHighlighter_Enabled_EndsWithReset(t *testing.T) {
	const ansiReset = "\033[0m"
	h := highlight.New(true)
	for _, level := range []string{"ERROR", "WARN", "INFO", "DEBUG", "FATAL"} {
		got := h.Format("some log line", parser.ParseSeverity(level))
		if !strings.HasSuffix(got, ansiReset) {
			t.Errorf("level %s: expected output to end with ANSI reset %q; got %q", level, ansiReset, got)
		}
	}
}
