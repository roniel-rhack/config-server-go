package lib

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteString_Success(t *testing.T) {
	var buf bytes.Buffer
	err := writeString(&buf, "hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "hello world" {
		t.Errorf("expected 'hello world', got '%v'", buf.String())
	}
}

func TestIsTruthyNode_Nil(t *testing.T) {
	if isTruthyNode(nil) {
		t.Errorf("expected nil node to be falsy")
	}
}

func TestIsTruthyNode_NullTag(t *testing.T) {
	node := makeScalar("!!null", "")
	if isTruthyNode(node) {
		t.Errorf("expected null node to be falsy")
	}
}

func TestIsTruthyNode_BoolTrue(t *testing.T) {
	for _, val := range []string{"true", "True", "TRUE", "y", "Y", "yes", "Yes", "YES", "on", "On", "ON"} {
		node := makeScalar("!!bool", val)
		if !isTruthyNode(node) {
			t.Errorf("expected '%v' to be truthy", val)
		}
	}
}

func TestIsTruthyNode_BoolFalse(t *testing.T) {
	for _, val := range []string{"false", "False", "FALSE", "n", "N", "no", "No", "NO", "off", "Off", "OFF"} {
		node := makeScalar("!!bool", val)
		if isTruthyNode(node) {
			t.Errorf("expected '%v' to be falsy", val)
		}
	}
}

func TestIsTruthyNode_NonBoolScalar(t *testing.T) {
	node := makeScalar("!!str", "hello")
	if !isTruthyNode(node) {
		t.Errorf("expected non-bool scalar to be truthy")
	}
}

func TestIsTruthyNode_IntScalar(t *testing.T) {
	node := makeScalar("!!int", "0")
	if !isTruthyNode(node) {
		t.Errorf("expected int scalar to be truthy (non-null, non-bool)")
	}
}

func TestReadStream_File(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(tmpFile, []byte("key: value\n"), 0644); err != nil {
		t.Fatal(err)
	}

	reader, err := readStream(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reader == nil {
		t.Fatal("expected non-nil reader")
	}
}

func TestReadStream_FileNotFound(t *testing.T) {
	_, err := readStream("/nonexistent/path/file.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestSafelyCloseFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	// Should not panic
	safelyCloseFile(f)
}

func TestSafelyCloseFile_AlreadyClosed(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	// Should log error but not panic
	safelyCloseFile(f)
}
