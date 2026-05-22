// Package stats provides a lightweight Stats type used by the slicer to
// accumulate and report processing metrics for a single logslice run.
//
// Usage:
//
//	s := stats.New()
//	for _, line := range lines {
//		s.RecordLine(len(line.Raw))
//		if matched {
//			s.RecordMatch()
//		}
//	}
//	s.Finish()
//	s.Write(os.Stderr)
package stats
