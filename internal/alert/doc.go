// Package alert provides threshold-based alerting over a rolling time window.
//
// An Alerter is configured with one or more Rules, each specifying a log
// severity level, a minimum count threshold, and a time window duration.
// As log lines are observed, the Alerter maintains per-severity timestamp
// buckets and fires an alert whenever the number of matching lines within
// the window meets or exceeds the threshold.
//
// Example:
//
//	rules := []alert.Rule{
//		{Severity: "ERROR", Threshold: 5, Window: time.Minute},
//	}
//	a := alert.New(os.Stderr, rules)
//	for _, line := range lines {
//		a.Observe(line)
//	}
package alert
