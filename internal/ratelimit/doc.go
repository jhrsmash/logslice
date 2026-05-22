// Package ratelimit provides a token-bucket rate limiter that controls
// how many log lines logslice emits per second.
//
// Usage:
//
//	l := ratelimit.New(1000) // allow up to 1 000 lines/s
//	defer l.Close()
//	for _, line := range lines {
//		if l.Allow() {
//			output.Write(line)
//		}
//	}
//
// A rate of 0 disables limiting entirely.
package ratelimit
