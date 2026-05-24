// Package replay re-emits log lines at a controlled speed,
// simulating real-time log streaming from a historical file.
package replay

import (
	"time"

	"github.com/user/logslice/internal/parser"
)

// Replayer replays log lines with timing derived from their timestamps,
// scaled by a configurable speed multiplier.
type Replayer struct {
	speed   float64
	parser  *parser.Parser
	stop    chan struct{}
}

// New creates a Replayer with the given speed multiplier.
// A speed of 1.0 replays in real time; 2.0 replays at double speed;
// 0 or negative values are treated as 1.0.
func New(p *parser.Parser, speed float64) *Replayer {
	if speed <= 0 {
		speed = 1.0
	}
	return &Replayer{
		speed:  speed,
		parser: p,
		stop:   make(chan struct{}),
	}
}

// Stop signals the replayer to cease emitting lines.
func (r *Replayer) Stop() {
	select {
	case <-r.stop:
	default:
		close(r.stop)
	}
}

// Run reads lines from src, sleeping between them to honour the
// original inter-line gap scaled by r.speed, and sends each
// *parser.LogLine on the returned channel.
// The channel is closed when src is exhausted or Stop is called.
func (r *Replayer) Run(src <-chan *parser.LogLine) <-chan *parser.LogLine {
	out := make(chan *parser.LogLine)
	go func() {
		defer close(out)
		var prev time.Time
		for line := range src {
			if line == nil {
				continue
			}
			if !prev.IsZero() && !line.Timestamp.IsZero() {
				gap := line.Timestamp.Sub(prev)
				if gap > 0 {
					scaled := time.Duration(float64(gap) / r.speed)
					select {
					case <-time.After(scaled):
					case <-r.stop:
						return
					}
				}
			}
			if !line.Timestamp.IsZero() {
				prev = line.Timestamp
			}
			select {
			case out <- line:
			case <-r.stop:
				return
			}
		}
	}()
	return out
}
