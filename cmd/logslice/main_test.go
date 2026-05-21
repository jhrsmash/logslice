package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, l := range lines {
		fmt.Fprintln(f, l)
	}
	return f.Name()
}

func TestRun_NoArgs(t *testing.T) {
	var out, errOut bytes.Buffer
	err := run([]string{}, &out, &errOut)
	if err == nil {
		t.Fatal("expected error for missing file argument")
	}
}

func TestRun_FileNotFound(t *testing.T) {
	var out, errOut bytes.Buffer
	err := run([]string{filepath.Join(t.TempDir(), "nonexistent.log")}, &out, &errOut)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_InvalidSince(t *testing.T) {
	logFile := writeTempLog(t, []string{"2024-01-02T10:00:00 INFO hello"})
	var out, errOut bytes.Buffer
	err := run([]string{"--since", "not-a-date", logFile}, &out, &errOut)
	if err == nil {
		t.Fatal("expected error for invalid --since")
	}
}

func TestRun_InvalidUntil(t *testing.T) {
	logFile := writeTempLog(t, []string{"2024-01-02T10:00:00 INFO hello"})
	var out, errOut bytes.Buffer
	err := run([]string{"--until", "bad-date", logFile}, &out, &errOut)
	if err == nil {
		t.Fatal("expected error for invalid --until")
	}
}

func TestRun_NoFilter(t *testing.T) {
	lines := []string{
		"2024-01-02T10:00:00 INFO starting up",
		"2024-01-02T10:01:00 WARN low memory",
		"2024-01-02T10:02:00 ERROR disk full",
	}
	logFile := writeTempLog(t, lines)

	var out, errOut bytes.Buffer
	if err := run([]string{logFile}, &out, &errOut); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := out.String()
	for _, l := range lines {
		if !strings.Contains(result, l) {
			t.Errorf("expected output to contain %q", l)
		}
	}
}

func TestRun_SeverityFilter(t *testing.T) {
	lines := []string{
		"2024-01-02T10:00:00 DEBUG verbose",
		"2024-01-02T10:01:00 INFO normal",
		"2024-01-02T10:02:00 ERROR critical",
	}
	logFile := writeTempLog(t, lines)

	var out, errOut bytes.Buffer
	if err := run([]string{"--severity", "ERROR", logFile}, &out, &errOut); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := out.String()
	if strings.Contains(result, "DEBUG") || strings.Contains(result, "INFO") {
		t.Error("expected DEBUG and INFO lines to be filtered out")
	}
	if !strings.Contains(result, "ERROR") {
		t.Error("expected ERROR line to be present")
	}
}
