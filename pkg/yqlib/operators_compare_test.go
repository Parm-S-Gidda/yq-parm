package yqlib

import (
	"testing"
	"time"
)

var compareOperatorScenarios = []expressionScenario{
	// ints, not equal
	{
		description: "Compare numbers (>)",
		document:    "a: 5\nb: 4",
		expression:  ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		expression: "(.k | length) >= 0",
		expected: []string{
			"D0, P[k], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		expression: `"2022-01-30T15:53:09Z" > "2020-01-30T15:53:09Z"`,
		expected: []string{
			"D0, P[], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5\nb: 4",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:     true,
		description: "Compare integers (>=)",
		document:    "a: 5\nb: 4",
		expression:  ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5\nb: 4",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},

	// ints, equal
	{
		skipDoc:     true,
		description: "Compare equal numbers (>)",
		document:    "a: 5\nb: 5",
		expression:  ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5\nb: 5",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		description: "Compare equal numbers (>=)",
		document:    "a: 5\nb: 5",
		expression:  ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5\nb: 5",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},

	// floats, not equal
	{
		skipDoc:    true,
		document:   "a: 5.2\nb: 4.1",
		expression: ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5.2\nb: 4.1",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5.2\nb: 4.1",
		expression: ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5.5\nb: 4.1",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},

	// floats, equal
	{
		skipDoc:    true,
		document:   "a: 5.5\nb: 5.5",
		expression: ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5.5\nb: 5.5",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5.1\nb: 5.1",
		expression: ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5.1\nb: 5.1",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},

	// strings, not equal
	{
		description:    "Compare strings",
		subdescription: "Compares strings by their bytecode.",
		document:       "a: zoo\nb: apple",
		expression:     ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: zoo\nb: apple",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: zoo\nb: apple",
		expression: ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: zoo\nb: apple",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},

	// strings, equal
	{
		skipDoc:    true,
		document:   "a: cat\nb: cat",
		expression: ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: cat\nb: cat",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: cat\nb: cat",
		expression: ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: cat\nb: cat",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},

	// datetime, not equal
	{
		description:    "Compare date times",
		subdescription: "You can compare date times. Assumes RFC3339 date time format, see [date-time operators](https://mikefarah.gitbook.io/yq/operators/date-time-operators) for more information.",
		document:       "a: 2021-01-01T03:10:00Z\nb: 2020-01-01T03:10:00Z",
		expression:     ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 2021-01-01T03:10:00Z\nb: 2020-01-01T03:10:00Z",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 2021-01-01T03:10:00Z\nb: 2020-01-01T03:10:00Z",
		expression: ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 2021-01-01T03:10:00Z\nb: 2020-01-01T03:10:00Z",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},

	// datetime, equal
	{
		skipDoc:    true,
		document:   "a: 2021-01-01T03:10:00Z\nb: 2021-01-01T03:10:00Z",
		expression: ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 2021-01-01T03:10:00Z\nb: 2021-01-01T03:10:00Z",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 2021-01-01T03:10:00Z\nb: 2021-01-01T03:10:00Z",
		expression: ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 2021-01-01T03:10:00Z\nb: 2021-01-01T03:10:00Z",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	// both null
	{
		description: "Both sides are null: > is false",
		expression:  ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		description: "Both sides are null: >= is true",
		expression:  ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},

	// one null
	{
		skipDoc:     true,
		description: "One side is null: > is false",
		document:    `a: 5`,
		expression:  ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: 5`,
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:     true,
		description: "One side is null: >= is false",
		document:    `a: 5`,
		expression:  ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: 5`,
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: 5`,
		expression: ".b <= .a",
		expected: []string{
			"D0, P[b], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: 5`,
		expression: ".b < .a",
		expected: []string{
			"D0, P[b], (!!bool)::false\n",
		},
	},
	// Test cases added by Mykhailo Isyp
	{
		skipDoc:    true,
		document:   "a: 5\nb: 5.0",
		expression: ".a >= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5\nb: 5.0",
		expression: ".a <= .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5\nb: 4.5",
		expression: ".a > .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 4.5\nb: 5",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:       true,
		document:      "a: 2021-01-01T00:00:00Z\nb: bad-date",
		expression:    ".a > .b",
		expectedError: `parsing time "bad-date" as "2006-01-02": cannot parse "bad-date" as "2006"`,
	},
	{
		skipDoc:       true,
		document:      "a: 5\nb: cat",
		expression:    ".a > .b",
		expectedError: "!!int not yet supported for comparison",
	},
	{
		skipDoc:       true,
		document:      "a: cat\nb: 5",
		expression:    ".a < .b",
		expectedError: "!!str not yet supported for comparison",
	},
	{
		skipDoc:    true,
		document:   `a: null`,
		expression: `.a >= .b`,
		expected: []string{
			"D0, P[a], (!!bool)::true\n",
		},
	},
	{
		skipDoc:    true,
		document:   `a: null`,
		expression: `.a > .b`,
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:       true,
		document:      "a: not-an-int\nb: 5",
		expression:    ".a > .b",
		expectedError: "!!str not yet supported for comparison",
	},
	{
		skipDoc:       true,
		document:      "a: 5\nb: not-an-int",
		expression:    ".a < .b",
		expectedError: "!!int not yet supported for comparison",
	},
	{
		skipDoc:       true,
		document:      "a: 5.1\nb: bad-float",
		expression:    ".a > .b",
		expectedError: "!!float not yet supported for comparison",
	},
	{
		skipDoc:    true,
		document:   "a: null\nb: 5",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	{
		skipDoc:    true,
		document:   "a: 5\nb: null",
		expression: ".a < .b",
		expected: []string{
			"D0, P[a], (!!bool)::false\n",
		},
	},
	// End of test cases added by Mykhailo Isyp
}

func TestCompareOperatorScenarios(t *testing.T) {
	for _, tt := range compareOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "compare", compareOperatorScenarios)
}

var minOperatorScenarios = []expressionScenario{
	{
		description: "Minimum int",
		document:    "[99, 16, 12, 6, 66]\n",
		expression:  `min`,
		expected: []string{
			"D0, P[3], (!!int)::6\n",
		},
	},
	{
		description: "Minimum string",
		document:    "[foo, bar, baz]\n",
		expression:  `min`,
		expected: []string{
			"D0, P[1], (!!str)::bar\n",
		},
	},
	{
		description: "Minimum of empty",
		document:    "[]\n",
		expression:  `min`,
		expected:    []string{},
	},
}

func TestMinOperatorScenarios(t *testing.T) {
	for _, tt := range minOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "min", minOperatorScenarios)
}

var maxOperatorScenarios = []expressionScenario{
	{
		description: "Maximum int",
		document:    "[99, 16, 12, 6, 66]\n",
		expression:  `max`,
		expected: []string{
			"D0, P[0], (!!int)::99\n",
		},
	},
	{
		description: "Maximum string",
		document:    "[foo, bar, baz]\n",
		expression:  `max`,
		expected: []string{
			"D0, P[0], (!!str)::foo\n",
		},
	},
	{
		description: "Maximum of empty",
		document:    "[]\n",
		expression:  `max`,
		expected:    []string{},
	},
}

func TestMaxOperatorScenarios(t *testing.T) {
	for _, tt := range maxOperatorScenarios {
		testScenario(t, &tt)
	}
	documentOperatorScenarios(t, "max", maxOperatorScenarios)
}

// Test cases added by Mykhailo Isyp
func compareNode(kind Kind, tag, value string) *CandidateNode {
	return &CandidateNode{
		Kind:  kind,
		Tag:   tag,
		Value: value,
	}
}

func TestCompareNilLhsReturnsFalse(t *testing.T) {
	fn := compare(compareTypePref{OrEqual: false})
	rhs := compareNode(ScalarNode, "!!int", "5")

	got, err := fn(nil, Context{}, nil, rhs)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "false" {
		t.Fatalf("expected false bool candidate, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestCompareNilRhsReturnsFalse(t *testing.T) {
	fn := compare(compareTypePref{OrEqual: false})
	lhs := compareNode(ScalarNode, "!!int", "5")

	got, err := fn(nil, Context{}, lhs, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "false" {
		t.Fatalf("expected false bool candidate, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestCompareBothNilOrEqualTrue(t *testing.T) {
	fn := compare(compareTypePref{OrEqual: true})

	got, err := fn(nil, Context{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "true" {
		t.Fatalf("expected true bool candidate, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestCompareReturnsErrorForScalarVsMap(t *testing.T) {
	fn := compare(compareTypePref{OrEqual: false})
	lhs := compareNode(ScalarNode, "!!int", "5")
	rhs := compareNode(MappingNode, "!!map", "")

	_, err := fn(nil, Context{}, lhs, rhs)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCompareDateTimeInvalidLHS(t *testing.T) {
	lhs := compareNode(ScalarNode, "!!timestamp", "bad-date")
	rhs := compareNode(ScalarNode, "!!timestamp", "2021-01-01T00:00:00Z")

	_, err := compareDateTime(time.RFC3339, compareTypePref{OrEqual: false}, lhs, rhs)
	if err == nil {
		t.Fatal("expected datetime parse error")
	}
}

func TestCompareDateTimeInvalidRHS(t *testing.T) {
	lhs := compareNode(ScalarNode, "!!timestamp", "2021-01-01T00:00:00Z")
	rhs := compareNode(ScalarNode, "!!timestamp", "bad-date")

	_, err := compareDateTime(time.RFC3339, compareTypePref{OrEqual: false}, lhs, rhs)
	if err == nil {
		t.Fatal("expected datetime parse error")
	}
}

func TestCompareScalarsInvalidIntLHS(t *testing.T) {
	lhs := compareNode(ScalarNode, "!!int", "bad")
	rhs := compareNode(ScalarNode, "!!int", "5")

	_, err := compareScalars(Context{}, compareTypePref{OrEqual: false}, lhs, rhs)
	if err == nil {
		t.Fatal("expected int parse error")
	}
}

func TestCompareScalarsInvalidIntRHS(t *testing.T) {
	lhs := compareNode(ScalarNode, "!!int", "5")
	rhs := compareNode(ScalarNode, "!!int", "bad")

	_, err := compareScalars(Context{}, compareTypePref{OrEqual: false}, lhs, rhs)
	if err == nil {
		t.Fatal("expected int parse error")
	}
}

func TestCompareScalarsInvalidFloatLHS(t *testing.T) {
	lhs := compareNode(ScalarNode, "!!float", "bad")
	rhs := compareNode(ScalarNode, "!!float", "5.2")

	_, err := compareScalars(Context{}, compareTypePref{OrEqual: false}, lhs, rhs)
	if err == nil {
		t.Fatal("expected float parse error")
	}
}

func TestCompareScalarsInvalidFloatRHS(t *testing.T) {
	lhs := compareNode(ScalarNode, "!!float", "5.2")
	rhs := compareNode(ScalarNode, "!!float", "bad")

	_, err := compareScalars(Context{}, compareTypePref{OrEqual: false}, lhs, rhs)
	if err == nil {
		t.Fatal("expected float parse error")
	}
}

func TestCompareBothNilOrEqualFalse(t *testing.T) {
	fn := compare(compareTypePref{OrEqual: false})

	got, err := fn(nil, Context{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "false" {
		t.Fatalf("expected false, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestCompareScalarVsMapReturnsError(t *testing.T) {
	fn := compare(compareTypePref{OrEqual: false})

	lhs := &CandidateNode{Kind: ScalarNode, Tag: "!!int", Value: "5"}
	rhs := &CandidateNode{Kind: MappingNode, Tag: "!!map"}

	_, err := fn(nil, Context{}, lhs, rhs)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCompareBothNilReturnsTrueWhenOrEqualTrue(t *testing.T) {
	fn := compare(compareTypePref{OrEqual: true})

	got, err := fn(nil, Context{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tag != "!!bool" || got.Value != "true" {
		t.Fatalf("expected true, got tag=%v value=%v", got.Tag, got.Value)
	}
}

func TestCompareScalarVsSequenceReturnsError(t *testing.T) {
	fn := compare(compareTypePref{OrEqual: false})

	lhs := &CandidateNode{Kind: ScalarNode, Tag: "!!int", Value: "5"}
	rhs := &CandidateNode{Kind: SequenceNode, Tag: "!!seq"}

	_, err := fn(nil, Context{}, lhs, rhs)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCompareSequenceVsScalarReturnsError(t *testing.T) {
	fn := compare(compareTypePref{OrEqual: false})

	lhs := &CandidateNode{Kind: SequenceNode, Tag: "!!seq"}
	rhs := &CandidateNode{Kind: ScalarNode, Tag: "!!int", Value: "5"}

	_, err := fn(nil, Context{}, lhs, rhs)
	if err == nil {
		t.Fatal("expected error")
	}
}

// End of test cases added by Mykhailo Isyp
