package lib

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestFormat(t *testing.T) {
	result := format(color.FgCyan)
	if !strings.HasPrefix(result, escape) {
		t.Errorf("expected ANSI escape prefix, got '%v'", result)
	}
	if !strings.Contains(result, "[") {
		t.Errorf("expected '[' in format string")
	}
	if !strings.HasSuffix(result, "m") {
		t.Errorf("expected 'm' suffix in format string")
	}
}

func TestFormat_Reset(t *testing.T) {
	result := format(color.Reset)
	if result == "" {
		t.Errorf("expected non-empty format string")
	}
}

func TestColorizeAndPrint(t *testing.T) {
	yamlBytes := []byte("name: test\nvalue: 123\n")
	var buf bytes.Buffer
	err := colorizeAndPrint(yamlBytes, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Errorf("expected non-empty output")
	}
}

func TestColorizeAndPrint_EmptyInput(t *testing.T) {
	var buf bytes.Buffer
	err := colorizeAndPrint([]byte(""), &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
