package output

import "fmt"

// ParseFormat converts a user-supplied string (e.g. from a CLI flag) into a
// Format constant. Accepted values are "raw" and "json" (case-insensitive).
// An error is returned for any unrecognised value.
func ParseFormat(s string) (Format, error) {
	switch s {
	case "raw", "Raw", "RAW", "":
		return FormatRaw, nil
	case "json", "Json", "JSON":
		return FormatJSON, nil
	default:
		return FormatRaw, fmt.Errorf("output: unknown format %q (want \"raw\" or \"json\")", s)
	}
}

// String implements the fmt.Stringer interface for Format.
func (f Format) String() string {
	switch f {
	case FormatJSON:
		return "json"
	default:
		return "raw"
	}
}
