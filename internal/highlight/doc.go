// Package highlight provides ANSI terminal colorization for log lines
// based on their severity level.
//
// When output is redirected to a file or pipe the highlighter should be
// constructed with enabled=false so that raw text is emitted without
// escape sequences.
package highlight
