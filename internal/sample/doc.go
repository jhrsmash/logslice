// Package sample implements deterministic log line sampling.
//
// When processing very large log files it can be useful to emit only a
// representative fraction of matching lines rather than every match.
// Sampler keeps a running counter and forwards one line for every N
// lines it receives, where N is the configured rate.
//
// A rate of 1 (the default) disables sampling and every line is forwarded.
// A rate of 10 forwards the 10th, 20th, 30th … matched line.
package sample
