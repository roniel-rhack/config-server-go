package lib

import (
	"container/list"
	"strings"
	"testing"
	"time"
)

func TestSingleReadonlyChildContext(t *testing.T) {
	parent := &Context{MatchingNodes: list.New()}
	node := makeScalar("!!str", "test")
	ctx := parent.SingleReadonlyChildContext(node)
	if ctx.MatchingNodes.Len() != 1 {
		t.Errorf("expected 1 matching node, got %v", ctx.MatchingNodes.Len())
	}
	if !ctx.DontAutoCreate {
		t.Errorf("expected DontAutoCreate true")
	}
}

func TestSingleChildContext(t *testing.T) {
	parent := &Context{MatchingNodes: list.New()}
	node := makeScalar("!!str", "test")
	ctx := parent.SingleChildContext(node)
	if ctx.MatchingNodes.Len() != 1 {
		t.Errorf("expected 1 matching node, got %v", ctx.MatchingNodes.Len())
	}
	if ctx.DontAutoCreate {
		t.Errorf("expected DontAutoCreate false")
	}
}

func TestSetGetDateTimeLayout_Custom(t *testing.T) {
	ctx := &Context{MatchingNodes: list.New()}
	ctx.SetDateTimeLayout("2006-01-02")
	result := ctx.GetDateTimeLayout()
	if result != "2006-01-02" {
		t.Errorf("expected '2006-01-02', got '%v'", result)
	}
}

func TestGetDateTimeLayout_Default(t *testing.T) {
	ctx := &Context{MatchingNodes: list.New()}
	result := ctx.GetDateTimeLayout()
	if result != time.RFC3339 {
		t.Errorf("expected RFC3339, got '%v'", result)
	}
}

func TestGetVariable_NilMap(t *testing.T) {
	ctx := &Context{MatchingNodes: list.New()}
	result := ctx.GetVariable("foo")
	if result != nil {
		t.Errorf("expected nil for nil variables map")
	}
}

func TestSetGetVariable(t *testing.T) {
	ctx := &Context{MatchingNodes: list.New()}
	l := list.New()
	l.PushBack(makeScalar("!!str", "val"))
	ctx.SetVariable("myvar", l)

	result := ctx.GetVariable("myvar")
	if result == nil {
		t.Fatalf("expected non-nil result")
	}
	if result.Len() != 1 {
		t.Errorf("expected 1 element, got %v", result.Len())
	}
}

func TestSetVariable_NilMapCreatesMap(t *testing.T) {
	ctx := &Context{MatchingNodes: list.New()}
	if ctx.Variables != nil {
		t.Fatalf("expected nil variables initially")
	}
	ctx.SetVariable("x", list.New())
	if ctx.Variables == nil {
		t.Errorf("expected variables map to be created")
	}
}

func TestChildContext_PreservesVariables(t *testing.T) {
	ctx := &Context{MatchingNodes: list.New()}
	l := list.New()
	l.PushBack(makeScalar("!!str", "val"))
	ctx.SetVariable("myvar", l)

	child := ctx.ChildContext(list.New())
	result := child.GetVariable("myvar")
	if result == nil {
		t.Fatalf("expected variable to be preserved in child")
	}
	if result.Len() != 1 {
		t.Errorf("expected 1 element in cloned variable")
	}
}

func TestToString(t *testing.T) {
	ctx := &Context{MatchingNodes: list.New(), DontAutoCreate: true}
	ctx.MatchingNodes.PushBack(makeScalar("!!str", "hello"))
	result := ctx.ToString()
	if !strings.Contains(result, "DontAutoCreate: true") {
		t.Errorf("expected DontAutoCreate in string")
	}
	if !strings.Contains(result, "1 results") {
		t.Errorf("expected '1 results' in string, got '%v'", result)
	}
}

func TestDeepClone(t *testing.T) {
	ctx := &Context{MatchingNodes: list.New()}
	node := makeScalar("!!str", "original")
	ctx.MatchingNodes.PushBack(node)

	clone := ctx.DeepClone()
	// Modify original
	node.Value = "modified"

	clonedNode := clone.MatchingNodes.Front().Value.(*CandidateNode)
	if clonedNode.Value != "original" {
		t.Errorf("deep clone should be independent, got '%v'", clonedNode.Value)
	}
}

func TestClone(t *testing.T) {
	ctx := &Context{MatchingNodes: list.New(), DontAutoCreate: true}
	ctx.MatchingNodes.PushBack(makeScalar("!!str", "hello"))
	clone := ctx.Clone()
	if clone.MatchingNodes.Len() != 1 {
		t.Errorf("expected 1 matching node in clone")
	}
	if !clone.DontAutoCreate {
		t.Errorf("expected DontAutoCreate preserved")
	}
}

func TestReadOnlyClone(t *testing.T) {
	ctx := &Context{MatchingNodes: list.New(), DontAutoCreate: false}
	clone := ctx.ReadOnlyClone()
	if !clone.DontAutoCreate {
		t.Errorf("expected DontAutoCreate true in readonly clone")
	}
}

func TestWritableClone(t *testing.T) {
	ctx := &Context{MatchingNodes: list.New(), DontAutoCreate: true}
	clone := ctx.WritableClone()
	if clone.DontAutoCreate {
		t.Errorf("expected DontAutoCreate false in writable clone")
	}
}
