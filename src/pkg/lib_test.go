package lib

import (
	clog "configTest/custom_logguer"
	"container/list"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	clog.Initialize()
	os.Exit(m.Run())
}

func makeScalar(tag, value string) *CandidateNode {
	return &CandidateNode{Kind: ScalarNode, Tag: tag, Value: value}
}

// --- recurseNodeArrayEqual ---

func TestRecurseNodeArrayEqual_Equal(t *testing.T) {
	lhs := &CandidateNode{Kind: SequenceNode, Content: []*CandidateNode{
		makeScalar("!!str", "a"), makeScalar("!!str", "b"),
	}}
	rhs := &CandidateNode{Kind: SequenceNode, Content: []*CandidateNode{
		makeScalar("!!str", "a"), makeScalar("!!str", "b"),
	}}
	if !recurseNodeArrayEqual(lhs, rhs) {
		t.Errorf("expected equal arrays to be equal")
	}
}

func TestRecurseNodeArrayEqual_DifferentLengths(t *testing.T) {
	lhs := &CandidateNode{Kind: SequenceNode, Content: []*CandidateNode{
		makeScalar("!!str", "a"),
	}}
	rhs := &CandidateNode{Kind: SequenceNode, Content: []*CandidateNode{
		makeScalar("!!str", "a"), makeScalar("!!str", "b"),
	}}
	if recurseNodeArrayEqual(lhs, rhs) {
		t.Errorf("expected different length arrays to not be equal")
	}
}

func TestRecurseNodeArrayEqual_DifferentContent(t *testing.T) {
	lhs := &CandidateNode{Kind: SequenceNode, Content: []*CandidateNode{
		makeScalar("!!str", "a"),
	}}
	rhs := &CandidateNode{Kind: SequenceNode, Content: []*CandidateNode{
		makeScalar("!!str", "b"),
	}}
	if recurseNodeArrayEqual(lhs, rhs) {
		t.Errorf("expected different content arrays to not be equal")
	}
}

// --- findInArray ---

func TestFindInArray_Found(t *testing.T) {
	arr := &CandidateNode{Kind: SequenceNode, Content: []*CandidateNode{
		makeScalar("!!str", "x"), makeScalar("!!str", "y"),
	}}
	idx := findInArray(arr, makeScalar("!!str", "y"))
	if idx != 1 {
		t.Errorf("expected index 1, got %v", idx)
	}
}

func TestFindInArray_NotFound(t *testing.T) {
	arr := &CandidateNode{Kind: SequenceNode, Content: []*CandidateNode{
		makeScalar("!!str", "x"),
	}}
	idx := findInArray(arr, makeScalar("!!str", "z"))
	if idx != -1 {
		t.Errorf("expected -1, got %v", idx)
	}
}

// --- findKeyInMap ---

func TestFindKeyInMap_Found(t *testing.T) {
	m := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		makeScalar("!!str", "key1"), makeScalar("!!str", "val1"),
		makeScalar("!!str", "key2"), makeScalar("!!str", "val2"),
	}}
	idx := findKeyInMap(m, makeScalar("!!str", "key2"))
	if idx != 2 {
		t.Errorf("expected index 2, got %v", idx)
	}
}

func TestFindKeyInMap_NotFound(t *testing.T) {
	m := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		makeScalar("!!str", "key1"), makeScalar("!!str", "val1"),
	}}
	idx := findKeyInMap(m, makeScalar("!!str", "missing"))
	if idx != -1 {
		t.Errorf("expected -1, got %v", idx)
	}
}

// --- recurseNodeObjectEqual ---

func TestRecurseNodeObjectEqual_Equal(t *testing.T) {
	lhs := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		makeScalar("!!str", "a"), makeScalar("!!str", "1"),
	}}
	rhs := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		makeScalar("!!str", "a"), makeScalar("!!str", "1"),
	}}
	if !recurseNodeObjectEqual(lhs, rhs) {
		t.Errorf("expected equal objects")
	}
}

func TestRecurseNodeObjectEqual_DifferentSize(t *testing.T) {
	lhs := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		makeScalar("!!str", "a"), makeScalar("!!str", "1"),
		makeScalar("!!str", "b"), makeScalar("!!str", "2"),
	}}
	rhs := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		makeScalar("!!str", "a"), makeScalar("!!str", "1"),
	}}
	if recurseNodeObjectEqual(lhs, rhs) {
		t.Errorf("expected not equal due to different sizes")
	}
}

func TestRecurseNodeObjectEqual_MissingKey(t *testing.T) {
	lhs := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		makeScalar("!!str", "a"), makeScalar("!!str", "1"),
	}}
	rhs := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		makeScalar("!!str", "b"), makeScalar("!!str", "1"),
	}}
	if recurseNodeObjectEqual(lhs, rhs) {
		t.Errorf("expected not equal due to missing key")
	}
}

func TestRecurseNodeObjectEqual_DifferentValue(t *testing.T) {
	lhs := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		makeScalar("!!str", "a"), makeScalar("!!str", "1"),
	}}
	rhs := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		makeScalar("!!str", "a"), makeScalar("!!str", "2"),
	}}
	if recurseNodeObjectEqual(lhs, rhs) {
		t.Errorf("expected not equal due to different value")
	}
}

// --- parseSnippet ---

func TestParseSnippet_EmptyString(t *testing.T) {
	node, err := parseSnippet("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Tag != "!!null" {
		t.Errorf("expected !!null tag, got %v", node.Tag)
	}
}

func TestParseSnippet_ValidString(t *testing.T) {
	node, err := parseSnippet("hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Tag != "!!str" {
		t.Errorf("expected !!str tag, got %v", node.Tag)
	}
	if node.Value != "hello world" {
		t.Errorf("expected 'hello world', got %v", node.Value)
	}
}

func TestParseSnippet_ValidInt(t *testing.T) {
	node, err := parseSnippet("42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Tag != "!!int" {
		t.Errorf("expected !!int tag, got %v", node.Tag)
	}
}

func TestParseSnippet_ValidBool(t *testing.T) {
	node, err := parseSnippet("true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Tag != "!!bool" {
		t.Errorf("expected !!bool tag, got %v", node.Tag)
	}
}

// --- recursiveNodeEqual ---

func TestRecursiveNodeEqual_SameScalars(t *testing.T) {
	a := makeScalar("!!str", "hello")
	b := makeScalar("!!str", "hello")
	if !recursiveNodeEqual(a, b) {
		t.Errorf("expected equal scalars")
	}
}

func TestRecursiveNodeEqual_DifferentScalars(t *testing.T) {
	a := makeScalar("!!str", "hello")
	b := makeScalar("!!str", "world")
	if recursiveNodeEqual(a, b) {
		t.Errorf("expected different scalars to not be equal")
	}
}

func TestRecursiveNodeEqual_DifferentKinds(t *testing.T) {
	a := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "a"}
	b := &CandidateNode{Kind: SequenceNode, Tag: "!!seq"}
	if recursiveNodeEqual(a, b) {
		t.Errorf("expected different kinds to not be equal")
	}
}

func TestRecursiveNodeEqual_NullNodes(t *testing.T) {
	a := makeScalar("!!null", "")
	b := makeScalar("!!null", "")
	if !recursiveNodeEqual(a, b) {
		t.Errorf("expected null nodes to be equal")
	}
}

func TestRecursiveNodeEqual_Sequences(t *testing.T) {
	a := &CandidateNode{Kind: SequenceNode, Content: []*CandidateNode{
		makeScalar("!!str", "x"),
	}}
	b := &CandidateNode{Kind: SequenceNode, Content: []*CandidateNode{
		makeScalar("!!str", "x"),
	}}
	if !recursiveNodeEqual(a, b) {
		t.Errorf("expected equal sequences")
	}
}

func TestRecursiveNodeEqual_Mappings(t *testing.T) {
	a := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		makeScalar("!!str", "k"), makeScalar("!!str", "v"),
	}}
	b := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		makeScalar("!!str", "k"), makeScalar("!!str", "v"),
	}}
	if !recursiveNodeEqual(a, b) {
		t.Errorf("expected equal mappings")
	}
}

func TestRecursiveNodeEqual_DifferentTags(t *testing.T) {
	a := makeScalar("!!str", "1")
	b := makeScalar("!!int", "1")
	if recursiveNodeEqual(a, b) {
		t.Errorf("expected scalars with different tags to not be equal")
	}
}

// --- parseInt64 ---

func TestParseInt64_Decimal(t *testing.T) {
	format, val, err := parseInt64("42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}
	if format != "%v" {
		t.Errorf("expected %%v format, got %v", format)
	}
}

func TestParseInt64_Hex(t *testing.T) {
	format, val, err := parseInt64("0xFF")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 255 {
		t.Errorf("expected 255, got %v", val)
	}
	if format != "0x%X" {
		t.Errorf("expected 0x%%X format, got %v", format)
	}
}

func TestParseInt64_Octal(t *testing.T) {
	format, val, err := parseInt64("0o17")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 15 {
		t.Errorf("expected 15, got %v", val)
	}
	if format != "0o%o" {
		t.Errorf("expected 0o%%o format, got %v", format)
	}
}

func TestParseInt64_Invalid(t *testing.T) {
	_, _, err := parseInt64("notanumber")
	if err == nil {
		t.Errorf("expected error for invalid number")
	}
}

// --- parseInt ---

func TestParseInt_Valid(t *testing.T) {
	val, err := parseInt("100")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 100 {
		t.Errorf("expected 100, got %v", val)
	}
}

func TestParseInt_Invalid(t *testing.T) {
	_, err := parseInt("abc")
	if err == nil {
		t.Errorf("expected error for invalid number")
	}
}

// --- comment functions ---

func TestHeadAndLineComment(t *testing.T) {
	node := &CandidateNode{HeadComment: "# head", LineComment: "# line"}
	result := headAndLineComment(node)
	if result != " head line" {
		t.Errorf("expected ' head line', got '%v'", result)
	}
}

func TestHeadComment_WithHash(t *testing.T) {
	node := &CandidateNode{HeadComment: "# my comment"}
	if headComment(node) != " my comment" {
		t.Errorf("unexpected result: %v", headComment(node))
	}
}

func TestHeadComment_WithoutHash(t *testing.T) {
	node := &CandidateNode{HeadComment: "no hash"}
	if headComment(node) != "no hash" {
		t.Errorf("unexpected result: %v", headComment(node))
	}
}

func TestLineComment_WithHash(t *testing.T) {
	node := &CandidateNode{LineComment: "# inline"}
	if lineComment(node) != " inline" {
		t.Errorf("unexpected result: %v", lineComment(node))
	}
}

func TestLineComment_WithoutHash(t *testing.T) {
	node := &CandidateNode{LineComment: "inline"}
	if lineComment(node) != "inline" {
		t.Errorf("unexpected result: %v", lineComment(node))
	}
}

func TestFootComment_WithHash(t *testing.T) {
	node := &CandidateNode{FootComment: "# foot"}
	if footComment(node) != " foot" {
		t.Errorf("unexpected result: %v", footComment(node))
	}
}

func TestFootComment_WithoutHash(t *testing.T) {
	node := &CandidateNode{FootComment: "foot"}
	if footComment(node) != "foot" {
		t.Errorf("unexpected result: %v", footComment(node))
	}
}

// --- NodesToString / NodeToString ---

func TestNodeToString_Nil(t *testing.T) {
	result := NodeToString(nil)
	if result != "-- nil --" {
		t.Errorf("expected '-- nil --', got '%v'", result)
	}
}

func TestNodeToString_Scalar(t *testing.T) {
	node := makeScalar("!!str", "hello")
	result := NodeToString(node)
	if result == "" {
		t.Errorf("expected non-empty string")
	}
}

func TestNodeToString_Alias(t *testing.T) {
	node := &CandidateNode{Kind: AliasNode, Tag: "!!str", Value: "ref"}
	result := NodeToString(node)
	if result == "" {
		t.Errorf("expected non-empty string")
	}
}

func TestNodeToString_WithContent(t *testing.T) {
	node := &CandidateNode{Kind: MappingNode, Tag: "!!map", Content: []*CandidateNode{
		makeScalar("!!str", "k"), makeScalar("!!str", "v"),
	}}
	result := NodeToString(node)
	if result == "" {
		t.Errorf("expected non-empty string for node with content")
	}
}

func TestNodesToString(t *testing.T) {
	l := list.New()
	l.PushBack(makeScalar("!!str", "hello"))
	result := NodesToString(l)
	if result == "" {
		t.Errorf("expected non-empty string")
	}
}

// --- NodeContentToString ---

func TestNodeContentToString(t *testing.T) {
	parent := &CandidateNode{Kind: MappingNode, Content: []*CandidateNode{
		makeScalar("!!str", "child1"),
	}}
	result := NodeContentToString(parent, 0)
	if result == "" {
		t.Errorf("expected non-empty content string")
	}
}

// --- KindString ---

func TestKindString_AllKinds(t *testing.T) {
	tests := []struct {
		kind     Kind
		expected string
	}{
		{ScalarNode, "ScalarNode"},
		{SequenceNode, "SequenceNode"},
		{MappingNode, "MappingNode"},
		{AliasNode, "AliasNode"},
		{Kind(999), "unknown!"},
	}
	for _, tc := range tests {
		result := KindString(tc.kind)
		if result != tc.expected {
			t.Errorf("KindString(%v): expected %v, got %v", tc.kind, tc.expected, result)
		}
	}
}
