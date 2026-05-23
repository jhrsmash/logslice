// Package pipeline wires together the core processing stages of logslice:
// reading, filtering, deduplication, sampling, rate-limiting, and writing.
package pipeline

import (
	"fmt"
	"io"

	"github.com/yourorg/logslice/internal/config"
	"github.com/yourorg/logslice/internal/dedupe"
	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/ratelimit"
	"github.com/yourorg/logslice/internal/sample"
	"github.com/yourorg/logslice/internal/stats"
	"github.com/yourorg/logslice/internal/truncate"
)

// Pipeline holds all processing stages for a single run.
type Pipeline struct {
	parser    *parser.Parser
	filter    *filter.Filter
	deduper   *dedupe.Deduper
	sampler   *sample.Sampler
	limiter   *ratelimit.Limiter
	truncator *truncate.Truncator
	writer    *output.Writer
	stats     *stats.Stats
}

// New constructs a Pipeline from the given config, writing results to w.
func New(cfg *config.Config, w io.Writer) (*Pipeline, error) {
	if cfg == nil {
		return nil, fmt.Errorf("pipeline: config must not be nil")
	}
	if w == nil {
		return nil, fmt.Errorf("pipeline: writer must not be nil")
	}

	p := &Pipeline{
		parser:    parser.New(),
		filter:    filter.New(cfg),
		deduper:   dedupe.New(cfg.DedupeWindow),
		sampler:   sample.New(cfg.SampleRate),
		limiter:   ratelimit.New(cfg.RateLimit),
		truncator: truncate.New(cfg.MaxLineLen),
		writer:    output.New(w, cfg.Format),
		stats:     stats.New(),
	}
	return p, nil
}

// Process parses rawLine, runs it through every stage, and writes it if it
// passes all gates. It returns the number of bytes written and any error.
func (p *Pipeline) Process(rawLine []byte) (int, error) {
	line, err := p.parser.Parse(rawLine)
	if err != nil || line == nil {
		return 0, nil
	}
	p.stats.RecordLine()

	if !p.filter.Match(line) {
		return 0, nil
	}
	if !p.deduper.Allow(line) {
		return 0, nil
	}
	if !p.sampler.Allow(line) {
		return 0, nil
	}
	if !p.limiter.Allow() {
		return 0, nil
	}

	line = p.truncator.Apply(line)
	p.stats.RecordMatch()
	return p.writer.Write(line)
}

// Stats returns the accumulated run statistics.
func (p *Pipeline) Stats() *stats.Stats { return p.stats }

// Close flushes and releases resources held by pipeline stages.
func (p *Pipeline) Close() error {
	p.limiter.Close()
	return p.writer.Flush()
}
