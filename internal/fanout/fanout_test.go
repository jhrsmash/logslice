package fanout_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/fanout"
	"github.com/yourorg/logslice/internal/parser"
)

// stubSink records every line it receives.
type stubSink struct {
	mu    sync.Mutex
	lines []*parser.LogLine
	err   error
}

func (s *stubSink) Write(line *parser.LogLine) error {
	s.mu.Lock()
	s.lines = append(s.lines, line)
	s.mu.Unlock()
	return s.err
}

func (s *stubSink) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.lines)
}

func makeLine(msg string) *parser.LogLine {
	return &parser.LogLine{Raw: msg, Timestamp: time.Now()}
}

func TestFanout_NilLine_Ignored(t *testing.T) {
	s := &stubSink{}
	f := fanout.New(s)
	errs := f.Send(nil)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if s.count() != 0 {
		t.Fatalf("expected 0 lines received, got %d", s.count())
	}
}

func TestFanout_NoSinks(t *testing.T) {
	f := fanout.New()
	errs := f.Send(makeLine("hello"))
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestFanout_SingleSink_ReceivesLine(t *testing.T) {
	s := &stubSink{}
	f := fanout.New(s)
	f.Send(makeLine("test"))
	if s.count() != 1 {
		t.Fatalf("expected 1 line, got %d", s.count())
	}
}

func TestFanout_MultipleSinks_AllReceive(t *testing.T) {
	a, b, c := &stubSink{}, &stubSink{}, &stubSink{}
	f := fanout.New(a, b, c)
	for i := 0; i < 5; i++ {
		f.Send(makeLine("msg"))
	}
	for _, s := range []*stubSink{a, b, c} {
		if s.count() != 5 {
			t.Fatalf("expected 5 lines, got %d", s.count())
		}
	}
}

func TestFanout_SinkError_CollectedAndDeliveryContines(t *testing.T) {
	bad := &stubSink{err: errors.New("write failed")}
	good := &stubSink{}
	f := fanout.New(bad, good)
	errs := f.Send(makeLine("oops"))
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if good.count() != 1 {
		t.Fatalf("good sink should still receive line, got %d", good.count())
	}
}

func TestFanout_Add_NilIgnored(t *testing.T) {
	f := fanout.New()
	f.Add(nil)
	if f.Len() != 0 {
		t.Fatalf("expected 0 sinks, got %d", f.Len())
	}
}

func TestFanout_Add_IncreasesLen(t *testing.T) {
	f := fanout.New()
	f.Add(&stubSink{})
	f.Add(&stubSink{})
	if f.Len() != 2 {
		t.Fatalf("expected 2 sinks, got %d", f.Len())
	}
}
