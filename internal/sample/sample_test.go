package sample_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/sample"
)

func makeLine() *parser.LogLine {
	return &parser.LogLine{
		Raw:      "2024-01-01T00:00:00Z INFO hello",
		Time:     time.Now(),
		Severity: parser.SeverityInfo,
		Message:  "hello",
	}
}

func TestSampler_RateOne_AllowsAll(t *testing.T) {
	s := sample.New(1)
	for i := 0; i < 10; i++ {
		if !s.Allow(makeLine()) {
			t.Fatalf("rate=1: expected Allow to return true on iteration %d", i)
		}
	}
}

func TestSampler_RateZero_TreatedAsOne(t *testing.T) {
	s := sample.New(0)
	if s.Rate() != 1 {
		t.Fatalf("expected rate 1, got %d", s.Rate())
	}
	if !s.Allow(makeLine()) {
		t.Fatal("expected Allow to return true")
	}
}

func TestSampler_RateN_EmitsEveryNth(t *testing.T) {
	const rate = 3
	s := sample.New(rate)
	var allowed []int
	for i := 1; i <= 9; i++ {
		if s.Allow(makeLine()) {
			allowed = append(allowed, i)
		}
	}
	expected := []int{3, 6, 9}
	if len(allowed) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, allowed)
	}
	for i, v := range expected {
		if allowed[i] != v {
			t.Errorf("position %d: expected %d, got %d", i, v, allowed[i])
		}
	}
}

func TestSampler_NilLine_ReturnsFalse(t *testing.T) {
	s := sample.New(1)
	if s.Allow(nil) {
		t.Fatal("expected Allow(nil) to return false")
	}
}

func TestSampler_Reset_RestartsCounter(t *testing.T) {
	s := sample.New(3)
	s.Allow(makeLine()) // 1
	s.Allow(makeLine()) // 2
	s.Reset()
	// after reset the next allowed should be at position 3 again
	results := make([]bool, 3)
	for i := range results {
		results[i] = s.Allow(makeLine())
	}
	if results[0] || results[1] {
		t.Error("expected first two calls after reset to be denied")
	}
	if !results[2] {
		t.Error("expected third call after reset to be allowed")
	}
}

func TestSampler_Rate_ReturnsConfigured(t *testing.T) {
	s := sample.New(7)
	if s.Rate() != 7 {
		t.Fatalf("expected 7, got %d", s.Rate())
	}
}
