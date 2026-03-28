package lib

import (
	"io"
	"strings"
	"testing"
)

func TestNewYamlDecoder(t *testing.T) {
	dec := NewYamlDecoder(NewDefaultYamlPreferences())
	if dec == nil {
		t.Fatalf("expected non-nil decoder")
	}
}

func TestYamlDecoder_SimpleScalar(t *testing.T) {
	dec := NewYamlDecoder(NewDefaultYamlPreferences())
	err := dec.Init(strings.NewReader("hello"))
	if err != nil {
		t.Fatalf("init error: %v", err)
	}
	node, err := dec.Decode()
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if node.Kind != ScalarNode {
		t.Errorf("expected ScalarNode, got %v", KindString(node.Kind))
	}
	if node.Value != "hello" {
		t.Errorf("expected 'hello', got '%v'", node.Value)
	}
}

func TestYamlDecoder_Mapping(t *testing.T) {
	dec := NewYamlDecoder(NewDefaultYamlPreferences())
	err := dec.Init(strings.NewReader("name: test\nport: 8080"))
	if err != nil {
		t.Fatalf("init error: %v", err)
	}
	node, err := dec.Decode()
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if node.Kind != MappingNode {
		t.Errorf("expected MappingNode, got %v", KindString(node.Kind))
	}
	if len(node.Content) < 4 {
		t.Errorf("expected at least 4 content nodes, got %v", len(node.Content))
	}
}

func TestYamlDecoder_Sequence(t *testing.T) {
	dec := NewYamlDecoder(NewDefaultYamlPreferences())
	err := dec.Init(strings.NewReader("- item1\n- item2"))
	if err != nil {
		t.Fatalf("init error: %v", err)
	}
	node, err := dec.Decode()
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if node.Kind != SequenceNode {
		t.Errorf("expected SequenceNode, got %v", KindString(node.Kind))
	}
	if len(node.Content) != 2 {
		t.Errorf("expected 2 content nodes, got %v", len(node.Content))
	}
}

func TestYamlDecoder_MultiDoc(t *testing.T) {
	dec := NewYamlDecoder(NewDefaultYamlPreferences())
	err := dec.Init(strings.NewReader("hello\n---\nworld"))
	if err != nil {
		t.Fatalf("init error: %v", err)
	}

	node1, err := dec.Decode()
	if err != nil {
		t.Fatalf("decode doc1 error: %v", err)
	}
	if node1.Value != "hello" {
		t.Errorf("expected 'hello', got '%v'", node1.Value)
	}

	node2, err := dec.Decode()
	if err != nil {
		t.Fatalf("decode doc2 error: %v", err)
	}
	if node2.Value != "world" {
		t.Errorf("expected 'world', got '%v'", node2.Value)
	}
}

func TestYamlDecoder_EOF(t *testing.T) {
	dec := NewYamlDecoder(NewDefaultYamlPreferences())
	err := dec.Init(strings.NewReader("hello"))
	if err != nil {
		t.Fatalf("init error: %v", err)
	}
	_, err = dec.Decode()
	if err != nil {
		t.Fatalf("first decode error: %v", err)
	}
	_, err = dec.Decode()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestYamlDecoder_LeadingComment(t *testing.T) {
	dec := NewYamlDecoder(NewDefaultYamlPreferences())
	err := dec.Init(strings.NewReader("# a comment\nname: test"))
	if err != nil {
		t.Fatalf("init error: %v", err)
	}
	node, err := dec.Decode()
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if node.LeadingContent == "" {
		t.Errorf("expected leading content with comment")
	}
}

func TestYamlDecoder_LeadingDirective(t *testing.T) {
	dec := NewYamlDecoder(NewDefaultYamlPreferences())
	err := dec.Init(strings.NewReader("%YAML 1.2\n---\nname: test"))
	if err != nil {
		t.Fatalf("init error: %v", err)
	}
	node, err := dec.Decode()
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if node.LeadingContent == "" {
		t.Errorf("expected leading content with directive")
	}
}

func TestYamlDecoder_OnlyComments(t *testing.T) {
	dec := NewYamlDecoder(NewDefaultYamlPreferences())
	err := dec.Init(strings.NewReader("# just a comment\n"))
	if err != nil {
		t.Fatalf("init error: %v", err)
	}
	node, err := dec.Decode()
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if node.LeadingContent == "" {
		t.Errorf("expected leading content from comment-only document")
	}
}

func TestYamlDecoder_NoLeadingContentPreProcessing(t *testing.T) {
	prefs := NewDefaultYamlPreferences()
	prefs.LeadingContentPreProcessing = false
	dec := NewYamlDecoder(prefs)
	err := dec.Init(strings.NewReader("name: test"))
	if err != nil {
		t.Fatalf("init error: %v", err)
	}
	node, err := dec.Decode()
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if node.Kind != MappingNode {
		t.Errorf("expected MappingNode, got %v", KindString(node.Kind))
	}
}

func TestYamlDecoder_NoLeadingContentPreProcessing_OnlyComments(t *testing.T) {
	prefs := NewDefaultYamlPreferences()
	prefs.LeadingContentPreProcessing = false
	dec := NewYamlDecoder(prefs)
	err := dec.Init(strings.NewReader("# just a comment\n"))
	if err != nil {
		t.Fatalf("init error: %v", err)
	}
	node, err := dec.Decode()
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if node.LeadingContent == "" {
		t.Errorf("expected leading content from buffered read")
	}
}

func TestDecode_DocumentSeparatorWithSpace(t *testing.T) {
	dec := NewYamlDecoder(NewDefaultYamlPreferences())
	err := dec.Init(strings.NewReader("--- \nkey: val\n"))
	if err != nil {
		t.Fatalf("init error: %v", err)
	}
	node, err := dec.Decode()
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if node.Kind != MappingNode {
		t.Errorf("expected MappingNode, got %v", KindString(node.Kind))
	}
}

func TestDecode_DocumentSeparatorWithNewline(t *testing.T) {
	dec := NewYamlDecoder(NewDefaultYamlPreferences())
	err := dec.Init(strings.NewReader("---\nkey: val\n"))
	if err != nil {
		t.Fatalf("init error: %v", err)
	}
	node, err := dec.Decode()
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if node.Kind != MappingNode {
		t.Errorf("expected MappingNode, got %v", KindString(node.Kind))
	}
}

func TestDecode_YamlDirective(t *testing.T) {
	dec := NewYamlDecoder(NewDefaultYamlPreferences())
	err := dec.Init(strings.NewReader("%YAML 1.2\n---\nkey: val\n"))
	if err != nil {
		t.Fatalf("init error: %v", err)
	}
	node, err := dec.Decode()
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if node.Kind != MappingNode {
		t.Errorf("expected MappingNode, got %v", KindString(node.Kind))
	}
	if node.LeadingContent == "" {
		t.Errorf("expected leading content with YAML directive")
	}
}

func TestDecode_EmptyLines(t *testing.T) {
	dec := NewYamlDecoder(NewDefaultYamlPreferences())
	err := dec.Init(strings.NewReader("\n\nkey: val\n"))
	if err != nil {
		t.Fatalf("init error: %v", err)
	}
	node, err := dec.Decode()
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if node.Kind != MappingNode {
		t.Errorf("expected MappingNode, got %v", KindString(node.Kind))
	}
}
