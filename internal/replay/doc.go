// Package replay provides time-accurate replay of historical log streams.
//
// A Replayer accepts a channel of *parser.LogLine values (as produced by
// the reader or multifile packages) and re-emits them on a new channel,
// sleeping between consecutive lines to honour the gap recorded in their
// timestamps.
//
// The playback speed is controlled by a multiplier:
//
//	1.0  – real time
//	2.0  – double speed
//	0.5  – half speed (slow motion)
//
// Lines whose timestamps are zero (unparseable) are forwarded immediately
// without introducing any delay.
//
// Example:
//
//	r := replay.New(p, 10.0)          // 10× speed
//	out := r.Run(linesCh)
//	for line := range out {
//	    fmt.Println(line.Raw)
//	}
package replay
