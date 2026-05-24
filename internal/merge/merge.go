// Package merge provides a multi-file log merger that emits log lines in
// chronological order across two or more sorted log streams.
package merge

import (
	"container/heap"

	"github.com/yourorg/logslice/internal/parser"
)

// Source is a function that returns the next parsed log line, or nil when
// the stream is exhausted.
type Source func() *parser.LogLine

// entry is a heap element coupling a line with the index of its source.
type entry struct {
	line   *parser.LogLine
	srcIdx int
}

// minHeap implements heap.Interface over entry values ordered by timestamp.
type minHeap []entry

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].line.Timestamp.Before(h[j].line.Timestamp) }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(entry)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Merger merges multiple sorted log sources into a single chronological stream.
type Merger struct {
	sources []Source
	h       *minHeap
	ready   bool
}

// New creates a Merger from the provided sources. At least one source is
// required; nil sources are silently ignored.
func New(sources ...Source) *Merger {
	filtered := make([]Source, 0, len(sources))
	for _, s := range sources {
		if s != nil {
			filtered = append(filtered, s)
		}
	}
	return &Merger{sources: filtered, h: &minHeap{}}
}

// init seeds the heap with the first line from every source.
func (m *Merger) init() {
	heap.Init(m.h)
	for i, src := range m.sources {
		if line := src(); line != nil {
			heap.Push(m.h, entry{line: line, srcIdx: i})
		}
	}
	m.ready = true
}

// Next returns the next log line in chronological order across all sources,
// or nil when all sources are exhausted.
func (m *Merger) Next() *parser.LogLine {
	if !m.ready {
		m.init()
	}
	if m.h.Len() == 0 {
		return nil
	}
	e := heap.Pop(m.h).(entry)
	// Refill from the same source.
	if next := m.sources[e.srcIdx](); next != nil {
		heap.Push(m.h, entry{line: next, srcIdx: e.srcIdx})
	}
	return e.line
}
