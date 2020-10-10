package treeops

import (
	"strings"
	"testing"

	"github.com/mikefarah/yq/v3/test"
	yaml "gopkg.in/yaml.v3"
)

var treeNavigator = NewDataTreeNavigator(NavigationPrefs{})
var treeCreator = NewPathTreeCreator()

func readDoc(t *testing.T, content string) []*CandidateNode {
	decoder := yaml.NewDecoder(strings.NewReader(content))
	var dataBucket yaml.Node
	err := decoder.Decode(&dataBucket)
	if err != nil {
		t.Error(err)
	}
	return []*CandidateNode{&CandidateNode{Node: &dataBucket, Document: 0}}
}

func resultsToString(results []*CandidateNode) string {
	var pretty string = ""
	for _, n := range results {
		pretty = pretty + "\n" + NodeToString(n)
	}
	return pretty
}

func TestDataTreeNavigatorSimple(t *testing.T) {

	nodes := readDoc(t, `a: 
  b: apple`)

	path, errPath := treeCreator.ParsePath("a")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a]
  Tag: !!map, Kind: MappingNode, Anchor: 
  b: apple
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorSubtractSimple(t *testing.T) {

	nodes := readDoc(t, `a: 
  b: apple
  c: camel`)

	path, errPath := treeCreator.ParsePath("a .- b")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a]
  Tag: !!map, Kind: MappingNode, Anchor: 
  c: camel
`
	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorSubtractTwice(t *testing.T) {

	nodes := readDoc(t, `a: 
  b: apple
  c: camel
  d: dingo`)

	path, errPath := treeCreator.ParsePath("a .- b OR a .- c")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a]
  Tag: !!map, Kind: MappingNode, Anchor: 
  d: dingo
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorSubtractWithUnion(t *testing.T) {

	nodes := readDoc(t, `a: 
  b: apple
  c: camel
  d: dingo`)

	path, errPath := treeCreator.ParsePath("a .- (b OR c)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a]
  Tag: !!map, Kind: MappingNode, Anchor: 
  d: dingo
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorSubtractArray(t *testing.T) {

	nodes := readDoc(t, `a: 
  - b: apple
  - b: sdfsd
  - b: apple`)

	path, errPath := treeCreator.ParsePath("a .- (b == a*)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a]
  Tag: !!seq, Kind: SequenceNode, Anchor: 
  - b: sdfsd
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorArraySimple(t *testing.T) {

	nodes := readDoc(t, `- b: apple`)

	path, errPath := treeCreator.ParsePath("[0]")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [0]
  Tag: !!map, Kind: MappingNode, Anchor: 
  b: apple
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorSimpleAssign(t *testing.T) {

	nodes := readDoc(t, `a: 
  b: apple`)

	path, errPath := treeCreator.ParsePath("a.b := frog")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a b]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  frog
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorArraySplat(t *testing.T) {

	nodes := readDoc(t, `- b: apple
- c: banana`)

	path, errPath := treeCreator.ParsePath("[*]")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [0]
  Tag: !!map, Kind: MappingNode, Anchor: 
  b: apple

-- Node --
  Document 0, path: [1]
  Tag: !!map, Kind: MappingNode, Anchor: 
  c: banana
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorSimpleDeep(t *testing.T) {

	nodes := readDoc(t, `a: 
  b: apple`)

	path, errPath := treeCreator.ParsePath("a.b")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a b]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  apple
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorSimpleMismatch(t *testing.T) {

	nodes := readDoc(t, `a: 
  c: apple`)

	path, errPath := treeCreator.ParsePath("a.b")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := ``

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorWild(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: apple
  mad: things`)

	path, errPath := treeCreator.ParsePath("a.*a*")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a cat]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  apple

-- Node --
  Document 0, path: [a mad]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  things
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorWildDeepish(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: {b: 3}
  mad: {b: 4}
  fad: {c: t}`)

	path, errPath := treeCreator.ParsePath("a.*a*.b")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a cat b]
  Tag: !!int, Kind: ScalarNode, Anchor: 
  3

-- Node --
  Document 0, path: [a mad b]
  Tag: !!int, Kind: ScalarNode, Anchor: 
  4
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorOrSimple(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: apple
  mad: things`)

	path, errPath := treeCreator.ParsePath("a.(cat or mad)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a cat]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  apple

-- Node --
  Document 0, path: [a mad]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  things
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorOrSimpleWithDepth(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: {b: 3}
  mad: {b: 4}
  fad: {c: t}`)

	path, errPath := treeCreator.ParsePath("a.(cat.* or fad.*)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a cat b]
  Tag: !!int, Kind: ScalarNode, Anchor: 
  3

-- Node --
  Document 0, path: [a fad c]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  t
`
	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorOrDeDupes(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: apple
  mad: things`)

	path, errPath := treeCreator.ParsePath("a.(cat or cat)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a cat]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  apple
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorAnd(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: apple
  pat: apple
  cow: apple
  mad: things`)

	path, errPath := treeCreator.ParsePath("a.(*t and c*)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a cat]
  Tag: !!str, Kind: ScalarNode, Anchor: 
  apple
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorEquals(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: {b: apple, c: yes}
  pat: {b: banana}
`)

	path, errPath := treeCreator.ParsePath("a.(b == apple)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a cat]
  Tag: !!map, Kind: MappingNode, Anchor: 
  {b: apple, c: yes}
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorArrayEquals(t *testing.T) {

	nodes := readDoc(t, `- { b: apple, animal: rabbit }
- { b: banana, animal: cat }
- { b: corn, animal: dog }`)

	path, errPath := treeCreator.ParsePath("(b == apple or animal == dog)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [0]
  Tag: !!map, Kind: MappingNode, Anchor: 
  {b: apple, animal: rabbit}

-- Node --
  Document 0, path: [2]
  Tag: !!map, Kind: MappingNode, Anchor: 
  {b: corn, animal: dog}
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorArrayEqualsDeep(t *testing.T) {

	nodes := readDoc(t, `apples:
  - { b: apple, animal: {legs: 2} }
  - { b: banana, animal: {legs: 4} }
  - { b: corn, animal: {legs: 6} }
`)

	path, errPath := treeCreator.ParsePath("apples(animal.legs == 4)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [apples 1]
  Tag: !!map, Kind: MappingNode, Anchor: 
  {b: banana, animal: {legs: 4}}
`

	test.AssertResult(t, expected, resultsToString(results))
}

func TestDataTreeNavigatorEqualsTrickey(t *testing.T) {

	nodes := readDoc(t, `a: 
  cat: {b: apso, c: {d : yes}}
  pat: {b: apple, c: {d : no}}
  sat: {b: apsy, c: {d : yes}}
  fat: {b: apple}
`)

	path, errPath := treeCreator.ParsePath("a.(b == ap* and c.d == yes)")
	if errPath != nil {
		t.Error(errPath)
	}
	results, errNav := treeNavigator.GetMatchingNodes(nodes, path)

	if errNav != nil {
		t.Error(errNav)
	}

	expected := `
-- Node --
  Document 0, path: [a cat]
  Tag: !!map, Kind: MappingNode, Anchor: 
  {b: apso, c: {d: yes}}

-- Node --
  Document 0, path: [a sat]
  Tag: !!map, Kind: MappingNode, Anchor: 
  {b: apsy, c: {d: yes}}
`

	test.AssertResult(t, expected, resultsToString(results))
}
