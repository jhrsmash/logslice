package alert

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// Rule defines a threshold-based alert condition.
type Rule struct {
	Severity  string
	Threshold int
	Window    time.Duration
}

// Alert represents a fired alert event.
type Alert struct {
	Rule      Rule
	Count     int
	FiredAt   time.Time
	Message   string
}

// Alerter watches log lines and fires alerts when severity counts
// exceed configured thresholds within a rolling time window.
type Alerter struct {
	mu      sync.Mutex
	rules   []Rule
	buckets map[string][]time.Time
	out     io.Writer
}

// New creates an Alerter that writes fired alerts to out.
func New(out io.Writer, rules []Rule) *Alerter {
	return &Alerter{
		rules:   rules,
		buckets: make(map[string][]time.Time),
		out:     out,
	}
}

// Observe records a log line and checks all rules, writing any fired
// alerts to the configured writer. Returns the list of fired alerts.
func (a *Alerter) Observe(line *parser.LogLine) []Alert {
	if line == nil {
		return nil
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	sev := line.Severity
	now := line.Timestamp
	a.buckets[sev] = append(a.buckets[sev], now)

	var fired []Alert
	for _, rule := range a.rules {
		if rule.Severity != sev {
			continue
		}
		cutoff := now.Add(-rule.Window)
		filtered := a.buckets[sev][:0]
		for _, t := range a.buckets[sev] {
			if !t.Before(cutoff) {
				filtered = append(filtered, t)
			}
		}
		a.buckets[sev] = filtered
		if len(filtered) >= rule.Threshold {
			al := Alert{
				Rule:    rule,
				Count:   len(filtered),
				FiredAt: now,
				Message: fmt.Sprintf("alert: %s count %d >= threshold %d in last %s", sev, len(filtered), rule.Threshold, rule.Window),
			}
			fired = append(fired, al)
			fmt.Fprintln(a.out, al.Message)
		}
	}
	return fired
}

// Reset clears all internal buckets, resetting counts for all severities.
func (a *Alerter) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.buckets = make(map[string][]time.Time)
}
