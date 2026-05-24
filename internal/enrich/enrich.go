// Package enrich provides log line enrichment by attaching derived or
// static key-value fields to each parsed log entry before it reaches
// downstream consumers.
//
// Enrichers are composable: multiple field sources (static tags, hostname,
// environment variables) can be chained together via a single Enricher that
// holds a list of providers.
package enrich

import (
	"fmt"
	"os"
	"sync"

	"logslice/internal/parser"
)

// Provider is a function that returns a set of key-value pairs to attach
// to a log line. It is called once per line so providers may be dynamic.
type Provider func() map[string]string

// Enricher attaches additional fields to log lines produced by the parser.
type Enricher struct {
	mu        sync.RWMutex
	providers []Provider
}

// New creates an Enricher with the given set of field providers.
// Providers are evaluated in order; later providers may overwrite keys
// set by earlier ones.
func New(providers ...Provider) *Enricher {
	return &Enricher{
		providers: providers,
	}
}

// AddProvider appends a provider to the enricher at runtime.
func (e *Enricher) AddProvider(p Provider) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.providers = append(e.providers, p)
}

// Enrich merges fields from all registered providers into line.Fields.
// If line is nil the call is a no-op.
func (e *Enricher) Enrich(line *parser.LogLine) {
	if line == nil {
		return
	}
	e.mu.RLock()
	defer e.mu.RUnlock()

	if line.Fields == nil {
		line.Fields = make(map[string]string)
	}
	for _, p := range e.providers {
		for k, v := range p() {
			line.Fields[k] = v
		}
	}
}

// StaticProvider returns a Provider that always emits the same key-value pairs.
func StaticProvider(fields map[string]string) Provider {
	// copy to avoid external mutation
	copy := make(map[string]string, len(fields))
	for k, v := range fields {
		copy[k] = v
	}
	return func() map[string]string { return copy }
}

// HostnameProvider returns a Provider that attaches the machine hostname
// under the key "host". The hostname is resolved once at construction time.
func HostnameProvider() Provider {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}
	fields := map[string]string{"host": host}
	return func() map[string]string { return fields }
}

// EnvProvider returns a Provider that reads the named environment variables
// and attaches them as fields. The env var name is used as the field key
// unless an explicit alias is supplied via the aliasMap (envVar -> fieldKey).
func EnvProvider(vars []string, aliasMap map[string]string) Provider {
	return func() map[string]string {
		out := make(map[string]string, len(vars))
		for _, v := range vars {
			key := v
			if alias, ok := aliasMap[v]; ok && alias != "" {
				key = alias
			}
			if val, set := os.LookupEnv(v); set {
				out[key] = val
			}
		}
		return out
	}
}

// SequenceProvider returns a Provider that attaches a monotonically increasing
// sequence number under the given key. Useful for ordering enriched lines.
func SequenceProvider(key string) Provider {
	if key == "" {
		key = "seq"
	}
	var mu sync.Mutex
	var n uint64
	return func() map[string]string {
		mu.Lock()
		n++
		v := n
		mu.Unlock()
		return map[string]string{key: fmt.Sprintf("%d", v)}
	}
}
