// Package normalize rewrites log line field keys and values to a canonical
// form before the lines are passed downstream to filters, writers, or
// aggregators.
//
// Usage:
//
//	rules := []normalize.Rule{
//		{From: "lvl",      To: "level"},
//		{From: "msg",      To: "message"},
//		{From: "ts",       To: "timestamp"},
//		{From: "hostname", To: "host",
//		 Transform: strings.ToLower},
//	}
//	norm := normalize.New(rules)
//	outLine := norm.Apply(inLine)
//
// Rules are matched case-insensitively on the source key. When a rule
// provides a Transform function the field value is passed through it before
// being stored under the canonical key. Fields with no matching rule are
// copied unchanged.
package normalize
