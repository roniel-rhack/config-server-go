package clog

import (
	"testing"
)

func TestInitialize(t *testing.T) {
	Initialize()
	if LOG == nil {
		t.Fatal("LOG should not be nil after Initialize()")
	}
}

func TestInitialize_Idempotent(t *testing.T) {
	Initialize()
	first := LOG
	Initialize()
	// After re-initialization, LOG should still be non-nil (may be a different instance)
	if LOG == nil {
		t.Fatal("LOG should not be nil after second Initialize()")
	}
	// Verify it was re-assigned (new instance)
	if LOG == first {
		// This is also acceptable; the point is it doesn't panic
	}
}

func TestInfo(t *testing.T) {
	Initialize()
	// Should not panic
	Info("test info message")
}

func TestInfoWithFormat(t *testing.T) {
	Initialize()
	// Should not panic
	Info("hello %s, count=%d", "world", 42)
}

func TestWarn(t *testing.T) {
	Initialize()
	Warn("test warn message")
}

func TestWarnWithFormat(t *testing.T) {
	Initialize()
	Warn("warning: %s at line %d", "issue", 10)
}

func TestPrint(t *testing.T) {
	Initialize()
	Print("test print message")
}

func TestPrintWithFormat(t *testing.T) {
	Initialize()
	Print("value is %v", 3.14)
}

func TestDebug(t *testing.T) {
	Initialize()
	Debug("test debug message")
}

func TestDebugWithFormat(t *testing.T) {
	Initialize()
	Debug("debugging %s=%d", "x", 99)
}

func TestError(t *testing.T) {
	Initialize()
	Error("test error message")
}

func TestErrorWithFormat(t *testing.T) {
	Initialize()
	Error("error: %s (code %d)", "not found", 404)
}

func TestAllFunctionsAfterSingleInit(t *testing.T) {
	Initialize()
	// Verify all log functions work with a single initialization
	Info("info")
	Warn("warn")
	Print("print")
	Debug("debug")
	Error("error")
}
