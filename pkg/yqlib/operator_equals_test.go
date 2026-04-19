package yqlib

import (
	"testing"
)

var equalsOperatorScenarios = []expressionScenario{
	{
		skipDoc:    true,
		expression: ".a == .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		expression: `(.k | length) == 0`,
		skipDoc:    true,
		expected: []string{
			"D0, P[k], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: cat`,
		expression: ".a == .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: cat`,
		expression: ".b == .a",
		expected: []string{
			"D0, P[b], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "cat",
		document2:  "dog",
		expression: "select(fi==0) == select(fi==1)",
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{}",
		expression: "(.a == .b) as $x | .",
		expected: []string{
			"D0, P[], (!!map)::{}\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{}",
		expression: ".a == .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{}",
		expression: "(.a != .b) as $x | .",
		expected: []string{
			"D0, P[], (!!map)::{}\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{}",
		expression: ".a != .b",
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{a: {b: 10}}",
		expression: "select(.c != null)",
		expected:   []string{},
	},
	{
		skipDoc:    true,
		document:   "{a: {b: 10}}",
		expression: "select(.d == .c)",
		expected: []string{
			"D0, P[], (!!map)::{a: {b: 10}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{a: {b: 10}}",
		expression: "select(null == .c)",
		expected: []string{
			"D0, P[], (!!map)::{a: {b: 10}}\n",
		},
	},
	{
		skipDoc:    true,
		document:   "{a: { b: {things: \"\"}, f: [1], g: [] }}",
		expression: ".. | select(. == \"\")",
		expected: []string{
			"D0, P[a b things], (!!str)::\n",
		},
	},
	{
		description: "Match string",
		document:    `[cat,goat,dog]`,
		expression:  `.[] | (. == "*at")`,
		expected: []string{
			"D0, P[0], (!!bool)::true\n",
			"D0, P[1], (!!bool)::true\n",
			"D0, P[2], (!!bool)::false\n",
		},
	},
	{
		description: "Don't match string",
		document:    `[cat,goat,dog]`,
		expression:  `.[] | (. != "*at")`,
		expected: []string{
			"D0, P[0], (!!bool)::false\n",
			"D0, P[1], (!!bool)::false\n",
			"D0, P[2], (!!bool)::true\n",
		},
	},
	{
		description: "Match number",
		document:    `[3, 4, 5]`,
		expression:  `.[] | (. == 4)`,
		expected: []string{
			"D0, P[0], (!!bool)::false\n",
			"D0, P[1], (!!bool)::true\n",
			"D0, P[2], (!!bool)::false\n",
		},
	},
	{
		description: "Don't match number",
		document:    `[3, 4, 5]`,
		expression:  `.[] | (. != 4)`,
		expected: []string{
			"D0, P[0], (!!bool)::true\n",
			"D0, P[1], (!!bool)::false\n",
			"D0, P[2], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: { cat: {b: apple, c: whatever}, pat: {b: banana} }`,
		expression: `.a | (.[].b == "apple")`,
		expected: []string{
			"D0, P[a cat b], (!!bool)::true\n",
			"D0, P[a pat b], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   ``,
		expression: `null == null`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description: "Match nulls",
		document:    ``,
		expression:  `null == ~`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		description: "Non existent key doesn't equal a value",
		document:    "a: frog",
		expression:  `select(.b != "thing")`,
		expected: []string{
			"D0, P[], (!!map)::a: frog\n",
		},
	},
	{
		description: "Two non existent keys are equal",
		document:    "a: frog",
		expression:  `select(.b == .c)`,
		expected: []string{
			"D0, P[], (!!map)::a: frog\n",
		},
	},
	// Test cases added by Mykhailo Isyp
	{
		skipDoc:    true,
		document:   `a: null`,
		expression: `.a == null`,
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: null`,
		expression: `.a != null`,
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: frog`,
		expression: `.a == null`,
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: frog`,
		expression: `.a != null`,
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {b: 1}}`,
		expression: `.a == null`,
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {b: 1}}`,
		expression: `.a != null`,
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: cat}`,
		expression: `.missing == null`,
		expected: []string{
			"D0, P[missing], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: cat}`,
		expression: `.missing != null`,
		expected: []string{
			"D0, P[], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[cat, cat, dog]`,
		expression: `.[] | (. == "cat")`,
		expected: []string{
			"D0, P[0], (!!bool)::true\n",
			"D0, P[1], (!!bool)::true\n",
			"D0, P[2], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `[cat, cat, dog]`,
		expression: `.[] | (. != "cat")`,
		expected: []string{
			"D0, P[0], (!!bool)::false\n",
			"D0, P[1], (!!bool)::false\n",
			"D0, P[2], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {nested: 1}, b: cat}`,
		expression: `.a == .b`,
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `{a: {nested: 1}, b: cat}`,
		expression: `.a != .b`,
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: null`,
		expression: `.a != .b`,
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: null`,
		expression: `.b != .a`,
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	// End of test cases added by Mykhailo Isyp
}

func TestEqualOperatorScenarios(t *testing.T) {
	for _, tt := range equalsOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "equals", equalsOperatorScenarios)
}

// Test cases added by Mykhailo Isyp
func TestIsNotEqualsScalarVsMapIsTrue(t *testing.T) {
	fn := isEquals(true)
	lhs := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "cat"}
	rhs := &CandidateNode{Kind: MappingNode, Tag: "!!map"}

	got, err := fn(nil, Context{}, lhs, rhs)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "true" {
		t.Fatalf("expected true bool candidate, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestIsEqualsScalarVsMapIsFalse(t *testing.T) {
	fn := isEquals(false)

	lhs := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "cat"}
	rhs := &CandidateNode{Kind: MappingNode, Tag: "!!map"}

	got, err := fn(nil, Context{}, lhs, rhs)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "false" {
		t.Fatalf("expected false, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestIsEqualsMapVsScalarIsFalse(t *testing.T) {
	fn := isEquals(false)

	lhs := &CandidateNode{Kind: MappingNode, Tag: "!!map"}
	rhs := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "cat"}

	got, err := fn(nil, Context{}, lhs, rhs)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "false" {
		t.Fatalf("expected false, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestIsNotEqualsMapVsScalarIsTrue(t *testing.T) {
	fn := isEquals(true)

	lhs := &CandidateNode{Kind: MappingNode, Tag: "!!map"}
	rhs := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "cat"}

	got, err := fn(nil, Context{}, lhs, rhs)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "true" {
		t.Fatalf("expected true, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestIsEqualsScalarVsSequenceIsFalse(t *testing.T) {
	fn := isEquals(false)

	lhs := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "cat"}
	rhs := &CandidateNode{Kind: SequenceNode, Tag: "!!seq"}

	got, err := fn(nil, Context{}, lhs, rhs)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "false" {
		t.Fatalf("expected false, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestIsNotEqualsScalarVsSequenceIsTrue(t *testing.T) {
	fn := isEquals(true)

	lhs := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "cat"}
	rhs := &CandidateNode{Kind: SequenceNode, Tag: "!!seq"}

	got, err := fn(nil, Context{}, lhs, rhs)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "true" {
		t.Fatalf("expected true, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestIsEqualsScalarVsNullNodeIsFalse(t *testing.T) {
	fn := isEquals(false)

	lhs := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "cat"}
	rhs := &CandidateNode{Kind: ScalarNode, Tag: "!!null", Value: "null"}

	got, err := fn(nil, Context{}, lhs, rhs)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "false" {
		t.Fatalf("expected false, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestIsEqualsSequenceVsScalarIsFalse(t *testing.T) {
	fn := isEquals(false)

	lhs := &CandidateNode{Kind: SequenceNode, Tag: "!!seq"}
	rhs := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "cat"}

	got, err := fn(nil, Context{}, lhs, rhs)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "false" {
		t.Fatalf("expected false, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestIsNotEqualsSequenceVsScalarIsTrue(t *testing.T) {
	fn := isEquals(true)

	lhs := &CandidateNode{Kind: SequenceNode, Tag: "!!seq"}
	rhs := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "cat"}

	got, err := fn(nil, Context{}, lhs, rhs)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "true" {
		t.Fatalf("expected true, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestIsNotEqualsScalarVsNullScalarIsTrue(t *testing.T) {
	fn := isEquals(true)

	lhs := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "cat"}
	rhs := &CandidateNode{Kind: ScalarNode, Tag: "!!null", Value: "null"}

	got, err := fn(nil, Context{}, lhs, rhs)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "true" {
		t.Fatalf("expected true, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestIsEqualsScalarVsMapWithSameValueStillFalse(t *testing.T) {
	fn := isEquals(false)

	lhs := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "cat"}
	rhs := &CandidateNode{Kind: MappingNode, Tag: "!!map", Value: "cat"}

	got, err := fn(nil, Context{}, lhs, rhs)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "false" {
		t.Fatalf("expected false, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestIsNotEqualsScalarVsMapWithSameValueStillTrue(t *testing.T) {
	fn := isEquals(true)

	lhs := &CandidateNode{Kind: ScalarNode, Tag: "!!str", Value: "cat"}
	rhs := &CandidateNode{Kind: MappingNode, Tag: "!!map", Value: "cat"}

	got, err := fn(nil, Context{}, lhs, rhs)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "true" {
		t.Fatalf("expected true, got tag=%v value=%v", got.Tag, got.Value)
	}
}

// End of test cases added by Mykhailo Isyp
