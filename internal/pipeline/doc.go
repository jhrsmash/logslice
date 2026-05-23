// Package pipeline assembles the individual logslice processing stages into a
// single, ordered execution chain.
//
// Stage order
//
//  1. Parser      – convert raw bytes into a structured LogLine.
//  2. Filter      – apply time-range and severity constraints.
//  3. Deduper     – suppress repeated identical messages within a window.
//  4. Sampler     – emit only every N-th matching line.
//  5. RateLimiter – cap output throughput (lines / second).
//  6. Truncator   – clip long message fields to a configured maximum.
//  7. Writer      – serialise to the configured output format (raw / JSON).
//
// Each stage is independently configurable via [config.Config]. Stages that
// are not configured (zero values) are no-ops and add negligible overhead.
package pipeline
