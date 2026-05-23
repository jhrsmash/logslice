// Package dedupe implements a consecutive-duplicate suppression filter
// for log lines.
//
// When a sequence of log lines carries the same message text, only the
// first occurrence is forwarded to the downstream consumer. Once the
// run ends (a different message arrives or Flush is called) a single
// summary line annotated with the repetition count is emitted, keeping
// output concise without losing information about repeated events.
package dedupe
