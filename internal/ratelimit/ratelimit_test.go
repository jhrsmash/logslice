package ratelimit

import (
	"testing"
	"time"
)

func TestLimiter_Disabled_AlwaysAllows(t *testing.T) {
	l := New(0)
	defer l.Close()
	for i := 0; i < 10_000; i++ {
		if !l.Allow() {
			t.Fatal("disabled limiter should always allow")
		}
	}
}

func TestLimiter_Enabled_AllowsUpToRate(t *testing.T) {
	const rate = 5
	l := New(rate)
	defer l.Close()

	allowed := 0
	for i := 0; i < rate*2; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != rate {
		t.Fatalf("expected %d allowed, got %d", rate, allowed)
	}
}

func TestLimiter_Enabled_DeniesAfterBucketEmpty(t *testing.T) {
	l := New(3)
	defer l.Close()

	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("call %d should be allowed", i)
		}
	}
	if l.Allow() {
		t.Fatal("4th call should be denied")
	}
}

func TestLimiter_Refills_AfterOneSec(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping timing test in short mode")
	}
	const rate = 4
	l := New(rate)
	defer l.Close()

	// drain the bucket
	for i := 0; i < rate; i++ {
		l.Allow()
	}
	if l.Allow() {
		t.Fatal("bucket should be empty")
	}

	// wait for refill
	time.Sleep(1100 * time.Millisecond)

	allowed := 0
	for i := 0; i < rate; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != rate {
		t.Fatalf("after refill expected %d, got %d", rate, allowed)
	}
}

func TestLimiter_Close_IdempotentOnDisabled(t *testing.T) {
	l := New(0)
	// should not panic
	l.Close()
	l.Close()
}
