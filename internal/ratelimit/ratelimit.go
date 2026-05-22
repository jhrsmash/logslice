package ratelimit

import (
	"sync"
	"time"
)

// Limiter throttles the number of log lines emitted per second.
// A zero MaxRate means no limiting is applied.
type Limiter struct {
	mu      sync.Mutex
	rate    int           // max lines per second
	bucket  int           // current token count
	last    time.Time
	ticker  *time.Ticker
	stop    chan struct{}
}

// New creates a Limiter that allows at most ratePerSec lines per second.
// If ratePerSec is 0 the limiter is disabled (all lines pass immediately).
func New(ratePerSec int) *Limiter {
	l := &Limiter{
		rate:   ratePerSec,
		bucket: ratePerSec,
		last:   time.Now(),
		stop:   make(chan struct{}),
	}
	if ratePerSec > 0 {
		l.ticker = time.NewTicker(time.Second)
		go l.refill()
	}
	return l
}

// Allow returns true if the caller may emit a line right now.
// When the limiter is disabled it always returns true.
func (l *Limiter) Allow() bool {
	if l.rate == 0 {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.bucket > 0 {
		l.bucket--
		return true
	}
	return false
}

// Close stops the background refill goroutine.
func (l *Limiter) Close() {
	if l.ticker != nil {
		l.ticker.Stop()
		close(l.stop)
	}
}

func (l *Limiter) refill() {
	for {
		select {
		case <-l.ticker.C:
			l.mu.Lock()
			l.bucket = l.rate
			l.mu.Unlock()
		case <-l.stop:
			return
		}
	}
}
