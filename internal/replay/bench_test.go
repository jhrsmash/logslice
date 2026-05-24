package replay

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

// BenchmarkReplayer_HighSpeed measures throughput when the speed multiplier
// is large enough that no real sleeping occurs.
func BenchmarkReplayer_HighSpeed(b *testing.B) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	const n = 1000
	lines := make([]*parser.LogLine, n)
	for i := range lines {
		lines[i] = makeLine(base.Add(time.Duration(i)*time.Millisecond), "benchmark line")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := New(nil, 1e12)
		out := r.Run(feedLines(lines))
		for range out {
		}
	}
}
