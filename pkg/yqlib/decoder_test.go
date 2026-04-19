package yqlib

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
	"github.com/zclconf/go-cty/cty"
)

type formatScenario struct {
	input          string
	indent         int
	expression     string
	expected       string
	description    string
	subdescription string
	skipDoc        bool
	scenarioType   string
	expectedError  string
}

func processFormatScenario(s formatScenario, decoder Decoder, encoder Encoder) (string, error) {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	if decoder == nil {
		decoder = NewYamlDecoder(ConfiguredYamlPreferences)
	}

	log.Debugf("reading docs")
	inputs, err := readDocuments(strings.NewReader(s.input), "sample.yml", 0, decoder)
	if err != nil {
		return "", err
	}

	log.Debugf("done reading the documents")

	expression := s.expression
	if expression == "" {
		expression = "."
	}

	exp, err := getExpressionParser().ParseExpression(expression)

	if err != nil {
		return "", err
	}

	context, err := NewDataTreeNavigator().GetMatchingNodes(Context{MatchingNodes: inputs}, exp)

	log.Debugf("Going to print: %v", NodesToString(context.MatchingNodes))

	if err != nil {
		return "", err
	}

	printer := NewPrinter(encoder, NewSinglePrinterWriter(writer))
	err = printer.PrintResults(context.MatchingNodes)
	if err != nil {
		return "", err
	}
	writer.Flush()

	return output.String(), nil
}

func mustProcessFormatScenario(s formatScenario, decoder Decoder, encoder Encoder) string {

	result, err := processFormatScenario(s, decoder, encoder)
	if err != nil {
		log.Errorf("Bad scenario %v: %v", s.description, err)
		return fmt.Sprintf("Bad scenario %v: %v", s.description, err.Error())
	}
	return result

}

//----------------------------------------------------------------------
//Parm

func TestConvertCtyValueToNode(t *testing.T) {
	
	n := convertCtyValueToNode(cty.NullVal(cty.String))
	if n.Tag != "!!null" {
		t.Errorf("expected null tag, got %s", n.Tag)
	}

	// string
	n = convertCtyValueToNode(cty.StringVal("hello"))
	if n.Value != "hello" {
		t.Errorf("expected hello, got %v", n.Value)
	}

	// bool
	n = convertCtyValueToNode(cty.BoolVal(true))
	if n.Value != "true" {
		t.Errorf("expected true, got %v", n.Value)
	}

	// int number
	n = convertCtyValueToNode(cty.NumberIntVal(5))
	if n.Value != "5" {
		t.Errorf("expected 5, got %v", n.Value)
	}

	// float number
	n = convertCtyValueToNode(cty.MustParseNumberVal("3.14"))
	if n.Value != "3.14" {
		t.Errorf("expected 3.14, got %v", n.Value)
	}

	// list / sequence
	n = convertCtyValueToNode(cty.ListVal([]cty.Value{
		cty.StringVal("a"),
		cty.StringVal("b"),
	}))
	if n.Kind != SequenceNode {
		t.Errorf("expected sequence node")
	}

	// map / object
	n = convertCtyValueToNode(cty.MapVal(map[string]cty.Value{
		"name": cty.StringVal("parm"),
	}))
	if n.Kind != MappingNode {
		t.Errorf("expected mapping node")
	}
}


