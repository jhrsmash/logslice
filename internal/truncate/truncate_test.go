package truncate_test

import (
	"strings"
	"testing"

	"github.com/example/logslice/internal/truncate"
)

func TestTruncator_Disabled(t *testing.T) {
	tr := truncate.New(0, "...")
	if tr.Enabled() {
		t.Fatal("expected truncator to be disabled")
	}
	s := strings.Repeat("a", 200)
	if got := tr.Apply(s); got != s {
		t.Fatalf("expected unchanged string, got length %d", len(got))
	}
}

func TestTruncator_ShortString(t *testing.T) {
	tr := truncate.New(80, "...")
	s := "short line"
	if got := tr.Apply(s); got != s {
		t.Fatalf("expected %q, got %q", s, got)
	}
}

func TestTruncator_ExactLength(t *testing.T) {
	tr := truncate.New(10, "...")
	s := "1234567890" // exactly 10 bytes
	if got := tr.Apply(s); got != s {
		t.Fatalf("expected unchanged string at exact limit, got %q", got)
	}
}

func TestTruncator_LongASCII(t *testing.T) {
	tr := truncate.New(10, "...")
	s := "hello world, this is a long log line"
	got := tr.Apply(s)
	if len(got) > 10 {
		t.Fatalf("expected len <= 10, got %d: %q", len(got), got)
	}
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("expected suffix '...', got %q", got)
	}
}

func TestTruncator_UTF8Boundary(t *testing.T) {
	// Each Japanese rune is 3 bytes; limit to 10 bytes with 3-byte suffix "→".
	tr := truncate.New(10, "→") // "→" is 3 bytes (UTF-8: e2 86 92)
	s := "日本語テスト" // 6 runes × 3 bytes = 18 bytes
	got := tr.Apply(s)
	if len(got) > 10 {
		t.Fatalf("result exceeds maxBytes: len=%d %q", len(got), got)
	}
	// Result must be valid UTF-8.
	if !isValidUTF8(got) {
		t.Fatalf("result is not valid UTF-8: %q", got)
	}
}

func TestTruncator_SuffixLargerThanMax(t *testing.T) {
	tr := truncate.New(2, "...")
	s := "hello"
	got := tr.Apply(s)
	if len(got) > 2 {
		t.Fatalf("expected len <= 2, got %d: %q", len(got), got)
	}
}

func TestTruncator_MaxBytes(t *testing.T) {
	tr := truncate.New(42, "...")
	if tr.MaxBytes() != 42 {
		t.Fatalf("expected MaxBytes()=42, got %d", tr.MaxBytes())
	}
}

// isValidUTF8 is a helper that checks all runes decoded successfully.
func isValidUTF8(s string) bool {
	for _, r := range s {
		if r == '\uFFFD' {
			return false
		}
	}
	return true
}
