package throttle

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

func makeLine(msg string) *parser.LogLine {
	return &parser.LogLine{
		Raw:      msg,
		Message:  msg,
		Severity: parser.SeverityInfo,
	}
}

func TestThrottler_Disabled_AlwaysAllows(t *testing.T) {
	th := New(0)
	for i := 0; i < 1000; i++ {
		if !th.Allow(makeLine("msg")) {
			t.Fatal("expected all lines to be allowed when rate is 0")
		}
	}
}

func TestThrottler_NilLine_AlwaysAllows(t *testing.T) {
	th := New(1)
	if !th.Allow(nil) {
		t.Fatal("expected nil line to be allowed")
	}
}

func TestThrottler_AllowsUpToRate(t *testing.T) {
	const rate = 5
	now := time.Now()
	th := New(rate)
	// Override clock so window never advances during the test.
	th.clock = func() time.Time { return now }

	for i := 0; i < rate; i++ {
		if !th.Allow(makeLine("msg")) {
			t.Fatalf("line %d should have been allowed", i)
		}
	}
	if th.count != rate {
		t.Fatalf("expected count=%d, got %d", rate, th.count)
	}
}

func TestThrottler_ResetsOnNewWindow(t *testing.T) {
	const rate = 3
	now := time.Now()
	th := New(rate)
	th.clock = func() time.Time { return now }

	// Exhaust the first window.
	for i := 0; i < rate; i++ {
		th.Allow(makeLine("msg"))
	}
	if th.count != rate {
		t.Fatalf("expected count=%d after first window, got %d", rate, th.count)
	}

	// Advance clock past the window boundary.
	advanced := now.Add(2 * time.Second)
	th.clock = func() time.Time { return advanced }

	if !th.Allow(makeLine("msg")) {
		t.Fatal("expected allow after window reset")
	}
	if th.count != 1 {
		t.Fatalf("expected count=1 after window reset, got %d", th.count)
	}
}

func TestThrottler_Rate_ReturnsConfiguredValue(t *testing.T) {
	th := New(42)
	if th.Rate() != 42 {
		t.Fatalf("expected rate=42, got %d", th.Rate())
	}
}
