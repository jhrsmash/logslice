// Package aggregate groups parsed log lines into fixed-width time buckets
// and counts occurrences by severity level.
//
// Usage:
//
//	agg := aggregate.New(time.Minute)
//	agg.Record(line.Timestamp, line.Severity)
//	for _, bucket := range agg.Snapshot() {
//		fmt.Println(bucket.Start, bucket.Total)
//	}
//
// Aggregator is safe for concurrent use.
package aggregate
