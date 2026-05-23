package burst

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeLine() *parser.LogLine {
	return &parser.LogLine{
		Timestamp: time.Now(),
		Severity:  parser.INFO,
		Raw:       "2024-01-01T00:00:00Z INFO test message",
	}
}

func TestDetector_NilLine(t *testing.T) {
	d := New(1.0, time.Second, nil)
	// Should not panic
	d.Record(nil)
	if d.Rate() != 0 {
		t.Fatalf("expected rate 0 after nil record, got %f", d.Rate())
	}
}

func TestDetector_RateBelowThreshold(t *testing.T) {
	d := New(100.0, time.Second, nil)
	d.Record(makeLine())
	d.Record(makeLine())
	if d.Rate() < 0 {
		t.Fatal("rate should be non-negative")
	}
}

func TestDetector_FiresCallbackOnBurst(t *testing.T) {
	var fired int32
	cb := func(rate float64) {
		atomic.StoreInt32(&fired, 1)
	}
	// threshold of 2 lines/sec over 1 second window
	d := New(2.0, time.Second, cb)

	for i := 0; i < 5; i++ {
		d.Record(makeLine())
	}

	// Give goroutine time to execute
	time.Sleep(20 * time.Millisecond)
	if atomic.LoadInt32(&fired) != 1 {
		t.Fatal("expected burst callback to fire")
	}
}

func TestDetector_DoesNotFireBelowThreshold(t *testing.T) {
	var fired int32
	cb := func(rate float64) {
		atomic.StoreInt32(&fired, 1)
	}
	// Very high threshold — should never fire for a single record
	d := New(1_000_000.0, time.Second, cb)
	d.Record(makeLine())
	time.Sleep(20 * time.Millisecond)
	if atomic.LoadInt32(&fired) != 0 {
		t.Fatal("callback should not have fired below threshold")
	}
}

func TestDetector_Reset_ClearsState(t *testing.T) {
	var fired int32
	cb := func(rate float64) { atomic.AddInt32(&fired, 1) }
	d := New(2.0, time.Second, cb)

	for i := 0; i < 5; i++ {
		d.Record(makeLine())
	}
	time.Sleep(20 * time.Millisecond)

	d.Reset()
	if d.Rate() != 0 {
		t.Fatalf("expected rate 0 after reset, got %f", d.Rate())
	}
}

func TestDetector_ZeroWindow_DefaultsToOneSecond(t *testing.T) {
	d := New(100.0, 0, nil)
	if d.window != time.Second {
		t.Fatalf("expected default window of 1s, got %v", d.window)
	}
}

func TestDetector_CallbackFiredOncePerBurst(t *testing.T) {
	var count int32
	cb := func(rate float64) { atomic.AddInt32(&count, 1) }
	d := New(2.0, time.Second, cb)

	for i := 0; i < 10; i++ {
		d.Record(makeLine())
	}
	time.Sleep(30 * time.Millisecond)

	if atomic.LoadInt32(&count) != 1 {
		t.Fatalf("expected callback fired exactly once, got %d", count)
	}
}
