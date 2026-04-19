package yqlib

import (
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func getExpressionParser() ExpressionParserInterface {
	InitExpressionParser()
	return ExpressionParser
}

func TestParserCreateMapColonOnItsOwn(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(":")
	test.AssertResultComplex(t, "':' expects 2 args but there is 0", err.Error())
}

func TestParserNoMatchingCloseBracket(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(".cat | with(.;.bob")
	test.AssertResultComplex(t, "bad expression - probably missing close bracket on WITH", err.Error())
}

func TestParserNoMatchingCloseCollect(t *testing.T) {
	_, err := getExpressionParser().ParseExpression("[1,2")
	test.AssertResultComplex(t, "bad expression, could not find matching `]`", err.Error())
}
func TestParserNoMatchingCloseObjectInCollect(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`[{"b": "c"]`)
	test.AssertResultComplex(t, "bad expression, could not find matching `}`", err.Error())
}

func TestParserNoMatchingCloseInCollect(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`[(.a]`)
	test.AssertResultComplex(t, "bad expression, could not find matching `)`", err.Error())
}

func TestParserNoMatchingCloseCollectObject(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`{"a": "b"`)
	test.AssertResultComplex(t, "bad expression, could not find matching `}`", err.Error())
}

func TestParserNoMatchingCloseCollectInCollectObject(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`{"b": [1}`)
	test.AssertResultComplex(t, "bad expression, could not find matching `]`", err.Error())
}

func TestParserNoMatchingCloseBracketInCollectObject(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`{"b": (1}`)
	test.AssertResultComplex(t, "bad expression, could not find matching `)`", err.Error())
}

func TestParserNoArgsForTwoArgOp(t *testing.T) {
	_, err := getExpressionParser().ParseExpression("=")
	test.AssertResultComplex(t, "'=' expects 2 args but there is 0", err.Error())
}

func TestParserOneLhsArgsForTwoArgOp(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(".a =")
	test.AssertResultComplex(t, "'=' expects 2 args but there is 1", err.Error())
}

func TestParserOneRhsArgsForTwoArgOp(t *testing.T) {
	_, err := getExpressionParser().ParseExpression("= .a")
	test.AssertResultComplex(t, "'=' expects 2 args but there is 1", err.Error())
}

func TestParserTwoArgsForTwoArgOp(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(".a = .b")
	test.AssertResultComplex(t, nil, err)
}

func TestParserNoArgsForOneArgOp(t *testing.T) {
	_, err := getExpressionParser().ParseExpression("explode")
	test.AssertResultComplex(t, "'explode' expects 1 arg but received none", err.Error())
}

func TestParserOneArgForOneArgOp(t *testing.T) {
	_, err := getExpressionParser().ParseExpression("explode(.)")
	test.AssertResultComplex(t, nil, err)
}

func TestParserExtraArgs(t *testing.T) {
	_, err := getExpressionParser().ParseExpression("sortKeys(.) explode(.)")
	test.AssertResultComplex(t, "bad expression, please check expression syntax", err.Error())
}

func TestParserEmptyExpression(t *testing.T) {
	_, err := getExpressionParser().ParseExpression("")
	test.AssertResultComplex(t, nil, err)
}

func TestParserSingleOperation(t *testing.T) {
	result, err := getExpressionParser().ParseExpression(".")
	test.AssertResultComplex(t, nil, err)
	if result == nil {
		t.Fatal("Expected non-nil result for single operation")
		return
	}
	if result.Operation == nil {
		t.Fatal("Expected operation to be set")
	}
}

func TestParserFirstOpWithZeroArgs(t *testing.T) {
	// Test the special case where firstOpType can accept zero args
	result, err := getExpressionParser().ParseExpression("first")
	test.AssertResultComplex(t, nil, err)
	if result == nil {
		t.Fatal("Expected non-nil result for first operation with zero args")
	}
}

func TestParserInvalidExpressionTree(t *testing.T) {
	// This tests the createExpressionTree function with malformed postfix
	parser := getExpressionParser().(*expressionParserImpl)

	// Create invalid postfix operations that would leave more than one item on stack
	invalidOps := []*Operation{
		{OperationType: &operationType{NumArgs: 0}},
		{OperationType: &operationType{NumArgs: 0}},
	}

	_, err := parser.createExpressionTree(invalidOps)
	test.AssertResultComplex(t, "bad expression, please check expression syntax", err.Error())
}

// Test cases added by Mykhailo Isyp
func TestParserBinaryOperatorCreatesChildren(t *testing.T) {
	result, err := getExpressionParser().ParseExpression(".a == .b")
	test.AssertResultComplex(t, nil, err)

	if result == nil || result.LHS == nil || result.RHS == nil {
		t.Fatal("expected root with both lhs and rhs")
	}
}

func TestParserBinaryOperatorSetsParentPointers(t *testing.T) {
	result, err := getExpressionParser().ParseExpression(".a == .b")
	test.AssertResultComplex(t, nil, err)

	if result.LHS == nil || result.RHS == nil {
		t.Fatal("expected lhs and rhs to be set")
	}
	if result.LHS.Parent != result {
		t.Fatal("expected lhs parent pointer to be set")
	}
	if result.RHS.Parent != result {
		t.Fatal("expected rhs parent pointer to be set")
	}
}

func TestParserTokeniseErrorUnclosedString(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`"unterminated`)
	if err == nil {
		t.Fatal("expected tokeniser error")
	}
}

func TestParserZeroArgOperationHasNoChildren(t *testing.T) {
	result, err := getExpressionParser().ParseExpression("first")
	test.AssertResultComplex(t, nil, err)

	if result == nil {
		t.Fatal("expected result")
	}
	if result.LHS != nil || result.RHS != nil {
		t.Fatal("expected zero-arg op to have no children")
	}
}

func TestParserCloseCollectWithoutOpen(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`]`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParserCloseBracketWithoutOpen(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`)`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParserMismatchedNestedCollects(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`[{"a": 1)]`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParserDeepMismatchedNestedCollects(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`[{"a": [.b)}]`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParserTraverseArrayCollectOnIndexExpression(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`.[.a]`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParserTraverseArrayCollectOnPipeExpression(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`(.a | .b)[.c]`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParserArrayCloseWithoutOpenAfterPipe(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`.a | ]`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParserRoundCloseWithoutOpenAfterPipe(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`.a | )`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParserNestedIndexExpression(t *testing.T) {
	_, err := getExpressionParser().ParseExpression(`.[.a][.b]`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParserNestedBinaryOperatorSetsInnerParents(t *testing.T) {
	result, err := getExpressionParser().ParseExpression("(.a == .b) and (.c == .d)")
	test.AssertResultComplex(t, nil, err)

	if result == nil {
		t.Fatal("expected result")
	}
	if result.LHS == nil || result.RHS == nil {
		t.Fatal("expected lhs and rhs on root")
	}
	if result.LHS.Parent != result {
		t.Fatal("expected lhs parent pointer on root")
	}
	if result.RHS.Parent != result {
		t.Fatal("expected rhs parent pointer on root")
	}

	if result.LHS.LHS == nil || result.LHS.RHS == nil {
		t.Fatal("expected children on left subtree")
	}
	if result.LHS.LHS.Parent != result.LHS {
		t.Fatal("expected left subtree lhs parent to be set")
	}
	if result.LHS.RHS.Parent != result.LHS {
		t.Fatal("expected left subtree rhs parent to be set")
	}

	if result.RHS.LHS == nil || result.RHS.RHS == nil {
		t.Fatal("expected children on right subtree")
	}
	if result.RHS.LHS.Parent != result.RHS {
		t.Fatal("expected right subtree lhs parent to be set")
	}
	if result.RHS.RHS.Parent != result.RHS {
		t.Fatal("expected right subtree rhs parent to be set")
	}
}

// End of test cases added by Mykhailo Isyp
