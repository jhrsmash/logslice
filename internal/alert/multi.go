package alert

import (
	"io"
	"sync"

	"github.com/logslice/logslice/internal/parser"
)

// MultiAlerter fans a single log line out to multiple Alerters and
// collects all fired alerts from each.
type MultiAlerter struct {
	mu       sync.Mutex
	alerters []*Alerter
}

// NewMulti creates a MultiAlerter from a set of (writer, rules) pairs.
// Each pair becomes an independent Alerter with its own bucket state.
func NewMulti(targets []struct {
	Out   io.Writer
	Rules []Rule
}) *MultiAlerter {
	ma := &MultiAlerter{}
	for _, t := range targets {
		ma.alerters = append(ma.alerters, New(t.Out, t.Rules))
	}
	return ma
}

// Observe passes line to every contained Alerter and returns the
// combined slice of all fired alerts.
func (ma *MultiAlerter) Observe(line *parser.LogLine) []Alert {
	if line == nil {
		return nil
	}
	ma.mu.Lock()
	defer ma.mu.Unlock()

	var all []Alert
	for _, a := range ma.alerters {
		all = append(all, a.Observe(line)...)
	}
	return all
}

// Reset resets all contained Alerters.
func (ma *MultiAlerter) Reset() {
	ma.mu.Lock()
	defer ma.mu.Unlock()
	for _, a := range ma.alerters {
		a.Reset()
	}
}

// Add appends a new Alerter to the MultiAlerter at runtime.
func (ma *MultiAlerter) Add(a *Alerter) {
	ma.mu.Lock()
	defer ma.mu.Unlock()
	ma.alerters = append(ma.alerters, a)
}
