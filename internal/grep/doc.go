// Package grep implements log-line text filtering via regular expressions.
//
// A Matcher is constructed once from a pattern string and reused across
// every line emitted by the reader pipeline. Matching is performed
// against the raw unparsed text of each LogLine, preserving the original
// formatting while still benefiting from the structured parse pass.
//
// Invert mode (analogous to grep -v) causes the Matcher to accept only
// lines whose raw text does NOT match the compiled expression.
//
// An empty pattern compiles to a no-op Matcher that accepts every line,
// so callers do not need to guard against a nil Matcher in hot paths.
package grep
