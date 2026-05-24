// Package fanout implements a one-to-many log-line multiplexer.
//
// A Fanout holds an ordered list of Sink implementations and forwards
// every incoming *parser.LogLine to each of them in registration order.
// Sinks that return errors do not interrupt delivery to subsequent sinks;
// all errors are collected and returned together so the caller can decide
// how to handle partial failures.
//
// Typical usage:
//
//	fw := fanout.New(primaryWriter, auditWriter)
//	fw.Add(metricsWriter)
//	for _, line := range lines {
//		if errs := fw.Send(line); len(errs) > 0 {
//			log.Println("fanout errors:", errs)
//		}
//	}
package fanout
