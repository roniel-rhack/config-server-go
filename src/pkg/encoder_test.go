package lib

import (
	"testing"
)

func TestMapKeysToStrings_MappingNode(t *testing.T) {
	node := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		{Kind: ScalarNode, Tag: "!!int", Value: "1"},
		{Kind: ScalarNode, Tag: "!!str", Value: "one"},
		{Kind: ScalarNode, Tag: "!!int", Value: "2"},
		{Kind: ScalarNode, Tag: "!!str", Value: "two"},
	}}
	mapKeysToStrings(node)

	// Keys (even indices) should have !!str tag
	if node.Content[0].Tag != "!!str" {
		t.Errorf("expected key at index 0 to have !!str tag, got %v", node.Content[0].Tag)
	}
	if node.Content[2].Tag != "!!str" {
		t.Errorf("expected key at index 2 to have !!str tag, got %v", node.Content[2].Tag)
	}
	// Values (odd indices) should remain unchanged
	if node.Content[1].Tag != "!!str" {
		t.Errorf("expected value at index 1 to keep !!str tag")
	}
	if node.Content[3].Tag != "!!str" {
		t.Errorf("expected value at index 3 to keep !!str tag")
	}
}

func TestMapKeysToStrings_RecursesIntoChildren(t *testing.T) {
	child := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		{Kind: ScalarNode, Tag: "!!int", Value: "inner"},
		{Kind: ScalarNode, Tag: "!!str", Value: "val"},
	}}
	node := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		{Kind: ScalarNode, Tag: "!!str", Value: "outer"},
		child,
	}}
	mapKeysToStrings(node)

	if child.Content[0].Tag != "!!str" {
		t.Errorf("expected nested key to have !!str tag, got %v", child.Content[0].Tag)
	}
}

func TestMapKeysToStrings_NonMapping(t *testing.T) {
	// Should not panic on non-mapping nodes
	node := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "hello"}
	mapKeysToStrings(node) // should not panic
}

func TestMapKeysToStrings_SequenceWithMappingChild(t *testing.T) {
	child := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		{Kind: ScalarNode, Tag: "!!int", Value: "key"},
		{Kind: ScalarNode, Tag: "!!str", Value: "val"},
	}}
	node := &CandidateNode{Kind: SequenceNode, Content: []*CandidateNode{child}}
	mapKeysToStrings(node)

	// Sequence node itself shouldn't modify keys, but recursion into children should
	if child.Content[0].Tag != "!!str" {
		t.Errorf("expected child mapping key to have !!str tag, got %v", child.Content[0].Tag)
	}
}
