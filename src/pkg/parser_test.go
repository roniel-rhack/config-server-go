package lib

import (
	"container/list"
	"testing"
)

func TestNewParser(t *testing.T) {
	enc := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	p := NewParser(enc)
	if p == nil {
		t.Fatalf("expected non-nil parser")
	}
}

func TestResultsToMap_EmptyList(t *testing.T) {
	enc := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	p := NewParser(enc)

	results, err := p.ResultsToMap(list.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty map, got %v entries", len(results))
	}
}

func TestResultsToMap_SingleKeyValue(t *testing.T) {
	enc := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	p := NewParser(enc)

	l := list.New()
	node := &CandidateNode{Kind: MappingNode, Tag: "!!map", Content: []*CandidateNode{
		makeScalar("!!str", "key"),
		makeScalar("!!str", "value"),
	}}
	l.PushBack(node)

	results, err := p.ResultsToMap(l)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 entry, got %v", len(results))
	}
	if results["key"] != "value" {
		t.Errorf("expected key=value, got key=%v", results["key"])
	}
}

func TestResultsToMap_SingleMapping(t *testing.T) {
	enc := NewPropertiesEncoder(NewDefaultPropertiesPreferences())
	p := NewParser(enc)

	l := list.New()
	node := &CandidateNode{Kind: MappingNode, Tag: "!!map", Content: []*CandidateNode{
		makeScalar("!!str", "server"),
		makeScalar("!!str", "localhost"),
	}}
	l.PushBack(node)

	results, err := p.ResultsToMap(l)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results["server"] != "localhost" {
		t.Errorf("expected server=localhost, got server=%v", results["server"])
	}
}
