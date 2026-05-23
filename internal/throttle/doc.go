// Package throttle implements a sliding-window throughput limiter for log
// lines. It is intended to sit between the reader and the output writer so
// that high-volume log files do not overwhelm downstream consumers.
//
// Usage:
//
//	th := throttle.New(500) // allow at most 500 lines/sec
//	for _, line := range lines {
//		if th.Allow(line) {
//			writer.Write(line)
//		}
//	}
//
// A rate of 0 disables throttling and every line is forwarded immediately.
package throttle
