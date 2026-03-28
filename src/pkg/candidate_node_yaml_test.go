package lib

import (
	"testing"

	yaml "gopkg.in/yaml.v3"
)

// --- MapYamlStyle ---

func TestMapYamlStyle_AllStyles(t *testing.T) {
	tests := []struct {
		input    yaml.Style
		expected Style
	}{
		{yaml.TaggedStyle, TaggedStyle},
		{yaml.DoubleQuotedStyle, DoubleQuotedStyle},
		{yaml.SingleQuotedStyle, SingleQuotedStyle},
		{yaml.LiteralStyle, LiteralStyle},
		{yaml.FoldedStyle, FoldedStyle},
		{yaml.FlowStyle, FlowStyle},
		{0, 0},
	}
	for _, tc := range tests {
		result := MapYamlStyle(tc.input)
		if result != tc.expected {
			t.Errorf("MapYamlStyle(%v): expected %v, got %v", tc.input, tc.expected, result)
		}
	}
}

func TestMapYamlStyle_Unknown(t *testing.T) {
	result := MapYamlStyle(yaml.Style(255))
	if result != Style(255) {
		t.Errorf("expected Style(255), got %v", result)
	}
}

// --- MapToYamlStyle ---

func TestMapToYamlStyle_AllStyles(t *testing.T) {
	tests := []struct {
		input    Style
		expected yaml.Style
	}{
		{TaggedStyle, yaml.TaggedStyle},
		{DoubleQuotedStyle, yaml.DoubleQuotedStyle},
		{SingleQuotedStyle, yaml.SingleQuotedStyle},
		{LiteralStyle, yaml.LiteralStyle},
		{FoldedStyle, yaml.FoldedStyle},
		{FlowStyle, yaml.FlowStyle},
		{0, 0},
	}
	for _, tc := range tests {
		result := MapToYamlStyle(tc.input)
		if result != tc.expected {
			t.Errorf("MapToYamlStyle(%v): expected %v, got %v", tc.input, tc.expected, result)
		}
	}
}

func TestMapToYamlStyle_Unknown(t *testing.T) {
	result := MapToYamlStyle(Style(255))
	if result != yaml.Style(255) {
		t.Errorf("expected yaml.Style(255), got %v", result)
	}
}

// --- UnmarshalYAML / MarshalYAML round trip ---

func TestUnmarshalMarshalYAML_Scalar(t *testing.T) {
	yamlNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: "hello",
	}
	candidate := &CandidateNode{}
	anchorMap := make(map[string]*CandidateNode)
	err := candidate.UnmarshalYAML(yamlNode, anchorMap)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if candidate.Kind != ScalarNode {
		t.Errorf("expected ScalarNode, got %v", candidate.Kind)
	}
	if candidate.Value != "hello" {
		t.Errorf("expected 'hello', got '%v'", candidate.Value)
	}

	// Marshal back
	result, err := candidate.MarshalYAML()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if result.Kind != yaml.ScalarNode {
		t.Errorf("expected yaml.ScalarNode")
	}
	if result.Value != "hello" {
		t.Errorf("expected 'hello' in marshaled node")
	}
}

func TestUnmarshalMarshalYAML_Mapping(t *testing.T) {
	yamlNode := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: "!!str", Value: "key"},
			{Kind: yaml.ScalarNode, Tag: "!!str", Value: "value"},
		},
	}
	candidate := &CandidateNode{}
	anchorMap := make(map[string]*CandidateNode)
	err := candidate.UnmarshalYAML(yamlNode, anchorMap)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if candidate.Kind != MappingNode {
		t.Errorf("expected MappingNode")
	}
	if len(candidate.Content) != 2 {
		t.Fatalf("expected 2 content nodes, got %v", len(candidate.Content))
	}

	// Marshal back
	result, err := candidate.MarshalYAML()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if result.Kind != yaml.MappingNode {
		t.Errorf("expected yaml.MappingNode")
	}
	if len(result.Content) != 2 {
		t.Errorf("expected 2 content nodes in marshaled result")
	}
}

func TestUnmarshalMarshalYAML_Sequence(t *testing.T) {
	yamlNode := &yaml.Node{
		Kind: yaml.SequenceNode,
		Tag:  "!!seq",
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: "!!str", Value: "item1"},
			{Kind: yaml.ScalarNode, Tag: "!!str", Value: "item2"},
		},
	}
	candidate := &CandidateNode{}
	anchorMap := make(map[string]*CandidateNode)
	err := candidate.UnmarshalYAML(yamlNode, anchorMap)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if candidate.Kind != SequenceNode {
		t.Errorf("expected SequenceNode")
	}
	if len(candidate.Content) != 2 {
		t.Fatalf("expected 2 content nodes, got %v", len(candidate.Content))
	}

	// Marshal back
	result, err := candidate.MarshalYAML()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if result.Kind != yaml.SequenceNode {
		t.Errorf("expected yaml.SequenceNode")
	}
}

func TestUnmarshalMarshalYAML_Alias(t *testing.T) {
	anchorNode := &yaml.Node{
		Kind:   yaml.ScalarNode,
		Tag:    "!!str",
		Value:  "anchored",
		Anchor: "myanchor",
	}
	aliasNode := &yaml.Node{
		Kind:  yaml.AliasNode,
		Alias: anchorNode,
	}

	// First unmarshal the anchor
	anchorMap := make(map[string]*CandidateNode)
	anchorCandidate := &CandidateNode{}
	err := anchorCandidate.UnmarshalYAML(anchorNode, anchorMap)
	if err != nil {
		t.Fatalf("unmarshal anchor error: %v", err)
	}

	// Then unmarshal the alias
	aliasCandidate := &CandidateNode{}
	err = aliasCandidate.UnmarshalYAML(aliasNode, anchorMap)
	if err != nil {
		t.Fatalf("unmarshal alias error: %v", err)
	}
	if aliasCandidate.Kind != AliasNode {
		t.Errorf("expected AliasNode")
	}

	// Marshal alias
	result, err := aliasCandidate.MarshalYAML()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if result.Kind != yaml.AliasNode {
		t.Errorf("expected yaml.AliasNode")
	}
}

func TestUnmarshalYAML_NullNode(t *testing.T) {
	yamlNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!null",
		Value: "",
	}
	parent := &CandidateNode{}
	child, err := parent.decodeIntoChild(yamlNode, make(map[string]*CandidateNode))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if child.Kind != ScalarNode {
		t.Errorf("expected ScalarNode for null")
	}
	if child.Tag != "!!null" {
		t.Errorf("expected !!null tag, got %v", child.Tag)
	}
}

func TestUnmarshalYAML_KindZero(t *testing.T) {
	yamlNode := &yaml.Node{
		Kind: 0,
		Tag:  "!!str",
	}
	candidate := &CandidateNode{}
	anchorMap := make(map[string]*CandidateNode)
	err := candidate.UnmarshalYAML(yamlNode, anchorMap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnmarshalYAML_InvalidKind(t *testing.T) {
	yamlNode := &yaml.Node{
		Kind: yaml.Kind(99),
	}
	candidate := &CandidateNode{}
	anchorMap := make(map[string]*CandidateNode)
	err := candidate.UnmarshalYAML(yamlNode, anchorMap)
	if err == nil {
		t.Errorf("expected error for invalid kind")
	}
}

func TestMarshalYAML_DefaultKind(t *testing.T) {
	// A node with Kind 0 (default, not matching any case)
	node := &CandidateNode{Kind: Kind(0), Tag: "!!str", Value: "test"}
	result, err := node.MarshalYAML()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Errorf("expected non-nil result")
	}
}
