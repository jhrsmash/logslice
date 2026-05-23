// Package context provides a cancellable processing context that carries
// a deadline, cancellation signal, and shared metadata for a logslice run.
package context

import (
	"context"
	"sync"
	"time"
)

// Meta holds arbitrary string key/value pairs attached to a run context.
type Meta map[string]string

// RunContext wraps a standard context.Context and adds logslice-specific
// metadata that pipeline stages can read without coupling to each other.
type RunContext struct {
	ctx    context.Context
	cancel context.CancelFunc

	mu   sync.RWMutex
	meta Meta

	StartedAt time.Time
}

// New creates a RunContext derived from parent. Call Cancel when the run
// is complete to release associated resources.
func New(parent context.Context) *RunContext {
	ctx, cancel := context.WithCancel(parent)
	return &RunContext{
		ctx:       ctx,
		cancel:    cancel,
		meta:      make(Meta),
		StartedAt: time.Now(),
	}
}

// WithTimeout creates a RunContext that is automatically cancelled after d.
func WithTimeout(parent context.Context, d time.Duration) *RunContext {
	ctx, cancel := context.WithTimeout(parent, d)
	return &RunContext{
		ctx:       ctx,
		cancel:    cancel,
		meta:      make(Meta),
		StartedAt: time.Now(),
	}
}

// Done returns the underlying done channel so RunContext satisfies select
// patterns used by pipeline stages.
func (r *RunContext) Done() <-chan struct{} { return r.ctx.Done() }

// Err returns the context error (nil until cancelled or deadline exceeded).
func (r *RunContext) Err() error { return r.ctx.Err() }

// Cancel cancels the RunContext, unblocking any stage waiting on Done.
func (r *RunContext) Cancel() { r.cancel() }

// Set stores a metadata value under key. Safe for concurrent use.
func (r *RunContext) Set(key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.meta[key] = value
}

// Get retrieves a metadata value. Returns "", false when key is absent.
func (r *RunContext) Get(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.meta[key]
	return v, ok
}

// Snapshot returns a shallow copy of all current metadata.
func (r *RunContext) Snapshot() Meta {
	r.mu.RLock()
	defer r.mu.RUnlock()
	copy := make(Meta, len(r.meta))
	for k, v := range r.meta {
		copy[k] = v
	}
	return copy
}

// Elapsed returns the duration since the RunContext was created.
func (r *RunContext) Elapsed() time.Duration {
	return time.Since(r.StartedAt)
}
