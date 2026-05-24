package replay

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeLine(ts time.Time, raw string) *parser.LogLine {
	return &parser.LogLine{Timestamp: ts, Raw: raw}
}

func feedLines(lines []*parser.LogLine) <-chan *parser.LogLine {
	ch := make(chan *parser.LogLine, len(lines))
	for _, l := range lines {
		ch <- l
	}
	close(ch)
	return ch
}

func TestReplayer_NilLinesSkipped(t *testing.T) {
	r := New(nil, 1.0)
	src := make(chan *parser.LogLine, 2)
	src <- nil
	src <- makeLine(time.Time{}, "hello")
	close(src)

	out := r.Run(src)
	var got []*parser.LogLine
	for l := range out {
		got = append(got, l)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 line, got %d", len(got))
	}
}

func TestReplayer_DefaultSpeed(t *testing.T) {
	r := New(nil, 0) // 0 → treated as 1.0
	if r.speed != 1.0 {
		t.Fatalf("expected speed 1.0, got %f", r.speed)
	}
}

func TestReplayer_ForwardsAllLines(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	lines := []*parser.LogLine{
		makeLine(base, "line1"),
		makeLine(base, "line2"),
		makeLine(base, "line3"),
	}
	// Use very high speed so gaps are negligible.
	r := New(nil, 1e9)
	out := r.Run(feedLines(lines))

	var got []string
	for l := range out {
		got = append(got, l.Raw)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
}

func TestReplayer_Stop(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	// Two lines 10 s apart at speed 1 → 10 s sleep; Stop should cancel.
	lines := []*parser.LogLine{
		makeLine(base, "a"),
		makeLine(base.Add(10*time.Second), "b"),
	}
	r := New(nil, 1.0)
	src := feedLines(lines)
	out := r.Run(src)

	// Drain first line then stop.
	<-out
	r.Stop()

	// Channel should close promptly.
	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected channel to be closed after Stop")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Stop to take effect")
	}
}

func TestReplayer_ZeroTimestampNoDelay(t *testing.T) {
	// Lines with zero timestamps must not cause a sleep.
	lines := []*parser.LogLine{
		makeLine(time.Time{}, "x"),
		makeLine(time.Time{}, "y"),
	}
	r := New(nil, 1.0)
	start := time.Now()
	out := r.Run(feedLines(lines))
	for range out {
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Fatalf("zero-timestamp lines caused unexpected delay: %v", elapsed)
	}
}
