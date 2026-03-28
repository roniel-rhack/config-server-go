package lib

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewPropertiesEncoder(t *testing.T) {
	enc := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	if enc == nil {
		t.Fatalf("expected non-nil encoder")
	}
}

func TestPropertiesEncoder_CanHandleAliases(t *testing.T) {
	enc := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	if enc.CanHandleAliases() {
		t.Errorf("expected CanHandleAliases to return false")
	}
}

func TestPropertiesEncoder_PrintDocumentSeparator(t *testing.T) {
	enc := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	var buf bytes.Buffer
	err := enc.PrintDocumentSeparator(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got '%v'", buf.String())
	}
}

func TestPropertiesEncoder_PrintLeadingContent_Plain(t *testing.T) {
	enc := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	var buf bytes.Buffer
	err := enc.PrintLeadingContent(&buf, "# a comment\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "# a comment") {
		t.Errorf("expected comment in output, got '%v'", buf.String())
	}
}

func TestPropertiesEncoder_PrintLeadingContent_DocSeparator(t *testing.T) {
	enc := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	var buf bytes.Buffer
	err := enc.PrintLeadingContent(&buf, "$yqDocSeparator$\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// The doc separator for properties is a no-op, so nothing should be written
	if buf.Len() != 0 {
		t.Errorf("expected empty output for doc separator, got '%v'", buf.String())
	}
}

func TestPropertiesEncoder_PrintLeadingContent_NoTrailingNewline(t *testing.T) {
	enc := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	var buf bytes.Buffer
	err := enc.PrintLeadingContent(&buf, "# comment without newline")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(buf.String(), "\n") {
		t.Errorf("expected trailing newline to be added")
	}
}

func TestPropertiesEncoder_Encode_Scalar(t *testing.T) {
	enc := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	var buf bytes.Buffer
	node := makeScalar("!!str", "hello")
	err := enc.Encode(&buf, node)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("expected 'hello' in output, got '%v'", buf.String())
	}
}

func TestPropertiesEncoder_Encode_Mapping(t *testing.T) {
	enc := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	var buf bytes.Buffer

	node := &CandidateNode{Kind: MappingNode, Tag: "!!map", Content: []*CandidateNode{
		makeScalar("!!str", "server"),
		{Kind: MappingNode, Tag: "!!map", Content: []*CandidateNode{
			makeScalar("!!str", "port"),
			makeScalar("!!str", "8080"),
		}},
	}}

	err := enc.Encode(&buf, node)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "server.port") {
		t.Errorf("expected 'server.port' in output, got '%v'", output)
	}
	if !strings.Contains(output, "8080") {
		t.Errorf("expected '8080' in output, got '%v'", output)
	}
}

func TestPropertiesEncoder_EncodeToMap_SimpleMapping(t *testing.T) {
	enc := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	node := &CandidateNode{Kind: MappingNode, Tag: "!!map", Content: []*CandidateNode{
		makeScalar("!!str", "name"),
		makeScalar("!!str", "value"),
	}}
	result, err := enc.EncodeToMap(node)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["name"] != "value" {
		t.Errorf("expected name=value, got name=%v", result["name"])
	}
}

func TestPropertiesEncoder_EncodeToMap_NestedMapping(t *testing.T) {
	enc := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	node := &CandidateNode{Kind: MappingNode, Tag: "!!map", Content: []*CandidateNode{
		makeScalar("!!str", "parent"),
		{Kind: MappingNode, Tag: "!!map", Content: []*CandidateNode{
			makeScalar("!!str", "child"),
			makeScalar("!!str", "val"),
		}},
	}}
	result, err := enc.EncodeToMap(node)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["parent.child"] != "val" {
		t.Errorf("expected parent.child=val, got %v", result["parent.child"])
	}
}

func TestPropertiesEncoder_EncodeToMap_Sequence(t *testing.T) {
	enc := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	node := &CandidateNode{Kind: MappingNode, Tag: "!!map", Content: []*CandidateNode{
		makeScalar("!!str", "items"),
		{Kind: SequenceNode, Tag: "!!seq", Content: []*CandidateNode{
			makeScalar("!!str", "a"),
			makeScalar("!!str", "b"),
		}},
	}}
	result, err := enc.EncodeToMap(node)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["items[0]"] != "a" {
		t.Errorf("expected items[0]=a, got %v", result["items[0]"])
	}
	if result["items[1]"] != "b" {
		t.Errorf("expected items[1]=b, got %v", result["items[1]"])
	}
}

func TestPropertiesEncoder_AppendPath_NoBrackets(t *testing.T) {
	prefs := NewDefaultPropertiesPreferences()
	prefs.UseArrayBrackets = false
	enc := NewPropertiesEncoder(prefs)

	node := &CandidateNode{Kind: MappingNode, Tag: "!!map", Content: []*CandidateNode{
		makeScalar("!!str", "items"),
		{Kind: SequenceNode, Tag: "!!seq", Content: []*CandidateNode{
			makeScalar("!!str", "a"),
		}},
	}}
	result, err := enc.EncodeToMap(node)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["items.0"] != "a" {
		t.Errorf("expected items.0=a, got %v", result["items.0"])
	}
}
