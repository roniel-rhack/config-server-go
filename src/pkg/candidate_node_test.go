package lib

import (
	"testing"
)

// --- createStringScalarNode ---

func TestCreateStringScalarNode(t *testing.T) {
	node := createStringScalarNode("hello")
	if node.Kind != ScalarNode {
		t.Errorf("expected ScalarNode kind")
	}
	if node.Tag != "!!str" {
		t.Errorf("expected !!str tag, got %v", node.Tag)
	}
	if node.Value != "hello" {
		t.Errorf("expected 'hello', got '%v'", node.Value)
	}
}

// --- createScalarNode ---

func TestCreateScalarNode_Float(t *testing.T) {
	node := createScalarNode(3.14, "3.14")
	if node.Tag != "!!float" {
		t.Errorf("expected !!float tag, got %v", node.Tag)
	}
}

func TestCreateScalarNode_Int(t *testing.T) {
	node := createScalarNode(42, "42")
	if node.Tag != "!!int" {
		t.Errorf("expected !!int tag, got %v", node.Tag)
	}
}

func TestCreateScalarNode_Bool(t *testing.T) {
	node := createScalarNode(true, "true")
	if node.Tag != "!!bool" {
		t.Errorf("expected !!bool tag, got %v", node.Tag)
	}
}

func TestCreateScalarNode_String(t *testing.T) {
	node := createScalarNode("text", "text")
	if node.Tag != "!!str" {
		t.Errorf("expected !!str tag, got %v", node.Tag)
	}
}

func TestCreateScalarNode_Nil(t *testing.T) {
	node := createScalarNode(nil, "")
	if node.Tag != "!!null" {
		t.Errorf("expected !!null tag, got %v", node.Tag)
	}
}

// --- CreateChild ---

func TestCreateChild(t *testing.T) {
	parent := &CandidateNode{Kind: MappingNode}
	child := parent.CreateChild()
	if child.Parent != parent {
		t.Errorf("expected child parent to be set")
	}
}

// --- SetDocument / GetDocument ---

func TestSetGetDocument(t *testing.T) {
	node := &CandidateNode{}
	node.SetDocument(5)
	if node.GetDocument() != 5 {
		t.Errorf("expected document 5, got %v", node.GetDocument())
	}
}

func TestGetDocument_WithParent(t *testing.T) {
	parent := &CandidateNode{}
	parent.SetDocument(3)
	child := &CandidateNode{Parent: parent}
	if child.GetDocument() != 3 {
		t.Errorf("expected document 3 from parent, got %v", child.GetDocument())
	}
}

// --- SetFilename / GetFilename ---

func TestSetGetFilename(t *testing.T) {
	node := &CandidateNode{}
	node.SetFilename("test.yaml")
	if node.GetFilename() != "test.yaml" {
		t.Errorf("expected 'test.yaml', got '%v'", node.GetFilename())
	}
}

func TestGetFilename_WithParent(t *testing.T) {
	parent := &CandidateNode{}
	parent.SetFilename("parent.yaml")
	child := &CandidateNode{Parent: parent}
	if child.GetFilename() != "parent.yaml" {
		t.Errorf("expected 'parent.yaml' from parent, got '%v'", child.GetFilename())
	}
}

// --- SetFileIndex / GetFileIndex ---

func TestSetGetFileIndex(t *testing.T) {
	node := &CandidateNode{}
	node.SetFileIndex(7)
	if node.GetFileIndex() != 7 {
		t.Errorf("expected 7, got %v", node.GetFileIndex())
	}
}

func TestGetFileIndex_WithParent(t *testing.T) {
	parent := &CandidateNode{}
	parent.SetFileIndex(2)
	child := &CandidateNode{Parent: parent}
	if child.GetFileIndex() != 2 {
		t.Errorf("expected 2 from parent, got %v", child.GetFileIndex())
	}
}

// --- GetKey ---

func TestGetKey_Normal(t *testing.T) {
	key := makeScalar("!!str", "mykey")
	node := &CandidateNode{Key: key}
	result := node.GetKey()
	if result != "0 - mykey" {
		t.Errorf("expected '0 - mykey', got '%v'", result)
	}
}

func TestGetKey_MapKey(t *testing.T) {
	node := &CandidateNode{IsMapKey: true, Value: "thekey"}
	result := node.GetKey()
	if result != "key-thekey-0 - " {
		t.Errorf("expected 'key-thekey-0 - ', got '%v'", result)
	}
}

// --- GetPath ---

func TestGetPath(t *testing.T) {
	key := makeScalar("!!str", "root")
	node := &CandidateNode{Key: key}
	path := node.GetPath()
	if len(path) != 1 || path[0] != "root" {
		t.Errorf("expected [root], got %v", path)
	}
}

// --- GetNicePath ---

func TestGetNicePath_Dots(t *testing.T) {
	parentKey := makeScalar("!!str", "parent")
	parent := &CandidateNode{Key: parentKey}

	childKey := makeScalar("!!str", "child")
	child := &CandidateNode{Key: childKey, Parent: parent}

	result := child.GetNicePath()
	if result != "parent.child" {
		t.Errorf("expected 'parent.child', got '%v'", result)
	}
}

func TestGetNicePath_IntIndex(t *testing.T) {
	key := makeScalar("!!int", "0")
	node := &CandidateNode{Key: key}
	result := node.GetNicePath()
	if result != "[0]" {
		t.Errorf("expected '[0]', got '%v'", result)
	}
}

func TestGetNicePath_Nested(t *testing.T) {
	rootKey := makeScalar("!!str", "root")
	root := &CandidateNode{Key: rootKey}

	midKey := makeScalar("!!str", "mid")
	mid := &CandidateNode{Key: midKey, Parent: root}

	leafKey := makeScalar("!!str", "leaf")
	leaf := &CandidateNode{Key: leafKey, Parent: mid}

	result := leaf.GetNicePath()
	if result != "root.mid.leaf" {
		t.Errorf("expected 'root.mid.leaf', got '%v'", result)
	}
}

// --- AsList ---

func TestAsList(t *testing.T) {
	node := makeScalar("!!str", "test")
	l := node.AsList()
	if l.Len() != 1 {
		t.Errorf("expected list of length 1, got %v", l.Len())
	}
	if l.Front().Value.(*CandidateNode) != node {
		t.Errorf("expected the node in the list")
	}
}

// --- SetParent ---

func TestSetParent(t *testing.T) {
	parent := &CandidateNode{Kind: MappingNode}
	child := &CandidateNode{Kind: ScalarNode}
	child.SetParent(parent)
	if child.Parent != parent {
		t.Errorf("expected parent to be set")
	}
}

// --- AddKeyValueChild ---

func TestAddKeyValueChild(t *testing.T) {
	parent := &CandidateNode{Kind: MappingNode}
	key := makeScalar("!!str", "k")
	value := makeScalar("!!str", "v")
	parent.AddKeyValueChild(key, value)

	if len(parent.Content) != 2 {
		t.Fatalf("expected 2 children, got %v", len(parent.Content))
	}
	if parent.Content[0].Value != "k" {
		t.Errorf("expected key 'k'")
	}
	if parent.Content[1].Value != "v" {
		t.Errorf("expected value 'v'")
	}
	if !parent.Content[0].IsMapKey {
		t.Errorf("expected key to be marked as map key")
	}
	if parent.Content[1].Key == nil {
		t.Errorf("expected value's Key to be set")
	}
}

// --- AddChild ---

func TestAddChild(t *testing.T) {
	parent := &CandidateNode{Kind: SequenceNode}
	child := makeScalar("!!str", "item")
	parent.AddChild(child)

	if len(parent.Content) != 1 {
		t.Fatalf("expected 1 child, got %v", len(parent.Content))
	}
	if parent.Content[0].Value != "item" {
		t.Errorf("expected 'item'")
	}
	if parent.Content[0].Key == nil {
		t.Errorf("expected key to be auto-generated")
	}
}

// --- AddChildren ---

func TestAddChildren_Mapping(t *testing.T) {
	parent := &CandidateNode{Kind: MappingNode}
	children := []*CandidateNode{
		makeScalar("!!str", "k1"),
		makeScalar("!!str", "v1"),
		makeScalar("!!str", "k2"),
		makeScalar("!!str", "v2"),
	}
	parent.AddChildren(children)
	if len(parent.Content) != 4 {
		t.Errorf("expected 4 content nodes, got %v", len(parent.Content))
	}
}

func TestAddChildren_Sequence(t *testing.T) {
	parent := &CandidateNode{Kind: SequenceNode}
	children := []*CandidateNode{
		makeScalar("!!str", "a"),
		makeScalar("!!str", "b"),
	}
	parent.AddChildren(children)
	if len(parent.Content) != 2 {
		t.Errorf("expected 2 content nodes, got %v", len(parent.Content))
	}
}

// --- GetValueRep ---

func TestGetValueRep_Int(t *testing.T) {
	node := makeScalar("!!int", "42")
	val, err := node.GetValueRep()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val.(int64) != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

func TestGetValueRep_Float(t *testing.T) {
	node := makeScalar("!!float", "3.14")
	val, err := node.GetValueRep()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	f, ok := val.(float64)
	if !ok || f != 3.14 {
		t.Errorf("expected 3.14, got %v", val)
	}
}

func TestGetValueRep_BoolTrue(t *testing.T) {
	node := makeScalar("!!bool", "true")
	val, err := node.GetValueRep()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != true {
		t.Errorf("expected true, got %v", val)
	}
}

func TestGetValueRep_BoolFalse(t *testing.T) {
	node := makeScalar("!!bool", "false")
	val, err := node.GetValueRep()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != false {
		t.Errorf("expected false, got %v", val)
	}
}

func TestGetValueRep_Null(t *testing.T) {
	node := makeScalar("!!null", "")
	val, err := node.GetValueRep()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}
}

func TestGetValueRep_String(t *testing.T) {
	node := makeScalar("!!str", "hello")
	val, err := node.GetValueRep()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "hello" {
		t.Errorf("expected 'hello', got '%v'", val)
	}
}

// --- CreateReplacement ---

func TestCreateReplacement(t *testing.T) {
	original := makeScalar("!!str", "orig")
	original.Key = makeScalar("!!str", "mykey")
	replacement := original.CreateReplacement(ScalarNode, "!!int", "42")
	if replacement.Tag != "!!int" {
		t.Errorf("expected !!int tag, got %v", replacement.Tag)
	}
	if replacement.Value != "42" {
		t.Errorf("expected '42', got '%v'", replacement.Value)
	}
	if replacement.Key == nil {
		t.Errorf("expected key to be copied from original")
	}
}

// --- CopyAsReplacement ---

func TestCopyAsReplacement(t *testing.T) {
	original := makeScalar("!!str", "orig")
	original.IsMapKey = true
	original.Value = "thekey"

	replacement := makeScalar("!!int", "99")
	result := original.CopyAsReplacement(replacement)
	if result.Key != original {
		t.Errorf("expected Key to be original when original is map key")
	}
}

func TestCopyAsReplacement_NotMapKey(t *testing.T) {
	original := makeScalar("!!str", "orig")
	original.Key = makeScalar("!!str", "existingkey")

	replacement := makeScalar("!!int", "99")
	result := original.CopyAsReplacement(replacement)
	if result.Key.Value != "existingkey" {
		t.Errorf("expected Key from original's key")
	}
}

// --- CreateReplacementWithComments ---

func TestCreateReplacementWithComments(t *testing.T) {
	original := makeScalar("!!str", "orig")
	original.HeadComment = "# head"
	original.LineComment = "# line"
	original.FootComment = "# foot"
	original.LeadingContent = "leading"

	replacement := original.CreateReplacementWithComments(MappingNode, "!!map", FlowStyle)
	if replacement.HeadComment != "# head" {
		t.Errorf("expected head comment copied")
	}
	if replacement.LineComment != "# line" {
		t.Errorf("expected line comment copied")
	}
	if replacement.FootComment != "# foot" {
		t.Errorf("expected foot comment copied")
	}
	if replacement.LeadingContent != "leading" {
		t.Errorf("expected leading content copied")
	}
	if replacement.Style != FlowStyle {
		t.Errorf("expected FlowStyle")
	}
}

// --- Copy ---

func TestCopy(t *testing.T) {
	node := makeScalar("!!str", "hello")
	node.HeadComment = "comment"
	node.SetDocument(2)
	node.SetFilename("test.yaml")

	copied := node.Copy()
	if copied.Value != "hello" {
		t.Errorf("expected 'hello', got '%v'", copied.Value)
	}
	if copied.HeadComment != "comment" {
		t.Errorf("expected comment copied")
	}
	// Modify original shouldn't affect copy
	node.Value = "changed"
	if copied.Value != "hello" {
		t.Errorf("copy should be independent")
	}
}

// --- guessTagFromCustomType ---

func TestGuessTagFromCustomType_StandardTag(t *testing.T) {
	node := makeScalar("!!str", "hello")
	result := node.guessTagFromCustomType()
	if result != "!!str" {
		t.Errorf("expected !!str, got %v", result)
	}
}

func TestGuessTagFromCustomType_CustomTagWithValue(t *testing.T) {
	node := makeScalar("!custom", "42")
	result := node.guessTagFromCustomType()
	if result != "!!int" {
		t.Errorf("expected !!int, got %v", result)
	}
}

func TestGuessTagFromCustomType_CustomTagWithStringValue(t *testing.T) {
	node := makeScalar("!custom", "hello")
	result := node.guessTagFromCustomType()
	if result != "!!str" {
		t.Errorf("expected !!str, got %v", result)
	}
}

func TestGuessTagFromCustomType_EmptyValue(t *testing.T) {
	node := makeScalar("!custom", "")
	result := node.guessTagFromCustomType()
	if result != "!custom" {
		t.Errorf("expected !custom, got %v", result)
	}
}

// --- CopyWithoutContent ---

func TestCopyWithoutContent(t *testing.T) {
	node := &CandidateNode{
		Kind: MappingNode,
		Tag:  "!!map",
		Content: []*CandidateNode{
			makeScalar("!!str", "k"),
			makeScalar("!!str", "v"),
		},
	}
	copied := node.CopyWithoutContent()
	if copied.Kind != MappingNode {
		t.Errorf("expected MappingNode kind")
	}
	if len(copied.Content) != 0 {
		t.Errorf("expected no content, got %v", len(copied.Content))
	}
}
