package lib

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewStreamEvaluator(t *testing.T) {
	se := NewStreamEvaluator()
	if se == nil {
		t.Fatal("expected non-nil StreamEvaluator")
	}
}

func TestEvaluateAndReturnMap_SimpleYaml(t *testing.T) {
	se := NewStreamEvaluator()
	reader := strings.NewReader("key: value\n")
	decoder := NewYamlDecoder(NewDefaultYamlPreferences())
	encoder := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	parser := NewParser(encoder)

	count, results, err := se.EvaluateAndReturnMap("test.yaml", reader, parser, decoder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 document, got %d", count)
	}
	if results["key"] != "value" {
		t.Errorf("expected results[key]=value, got %q", results["key"])
	}
}

func TestEvaluateAndReturnMap_MultiDoc(t *testing.T) {
	se := NewStreamEvaluator()
	reader := strings.NewReader("a: 1\n---\nb: 2\n")
	decoder := NewYamlDecoder(NewDefaultYamlPreferences())
	encoder := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	parser := NewParser(encoder)

	count, results, err := se.EvaluateAndReturnMap("test.yaml", reader, parser, decoder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count < 2 {
		t.Errorf("expected at least 2 documents, got %d", count)
	}
	if results["a"] != "1" {
		t.Errorf("expected results[a]=1, got %q", results["a"])
	}
	if results["b"] != "2" {
		t.Errorf("expected results[b]=2, got %q", results["b"])
	}
}

func TestEvaluateAndReturnMap_InvalidYaml(t *testing.T) {
	se := NewStreamEvaluator()
	reader := strings.NewReader(":\n  :\n    - :\n      invalid: [unterminated\n")
	decoder := NewYamlDecoder(NewDefaultYamlPreferences())
	encoder := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	parser := NewParser(encoder)

	_, _, err := se.EvaluateAndReturnMap("bad.yaml", reader, parser, decoder)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestEvaluateFilesAndReturnMap_RealFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.yaml")
	content := "server: localhost\nport: 8080\n"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	se := NewStreamEvaluator()
	decoder := NewYamlDecoder(NewDefaultYamlPreferences())
	encoder := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	parser := NewParser(encoder)

	results, err := se.EvaluateFilesAndReturnMap([]string{tmpFile}, parser, decoder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results["server"] != "localhost" {
		t.Errorf("expected server=localhost, got %q", results["server"])
	}
	if results["port"] != "8080" {
		t.Errorf("expected port=8080, got %q", results["port"])
	}
}

func TestEvaluateFilesAndReturnMap_FileNotFound(t *testing.T) {
	se := NewStreamEvaluator()
	decoder := NewYamlDecoder(NewDefaultYamlPreferences())
	encoder := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	parser := NewParser(encoder)

	_, err := se.EvaluateFilesAndReturnMap([]string{"/nonexistent/path/file.yaml"}, parser, decoder)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}
