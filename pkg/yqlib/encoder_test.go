//go:build !yq_nojson

package yqlib

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

func yamlToJSON(t *testing.T, sampleYaml string, indent int) string {
	t.Helper()
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	prefs := ConfiguredJSONPreferences.Copy()
	prefs.Indent = indent
	prefs.UnwrapScalar = false
	var jsonEncoder = NewJSONEncoder(prefs)
	inputs, err := readDocuments(strings.NewReader(sampleYaml), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}
	node := inputs.Front().Value.(*CandidateNode)
	log.Debugf("%v", NodeToString(node))
	// log.Debugf("Content[0] %v", NodeToString(node.Content[0]))

	err = jsonEncoder.Encode(writer, node)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	return strings.TrimSuffix(output.String(), "\n")
}

func TestJSONEncoderPreservesObjectOrder(t *testing.T) {
	var sampleYaml = `zabbix: winner
apple: great
banana:
- {cobra: kai, angus: bob}
`
	var expectedJSON = `{
  "zabbix": "winner",
  "apple": "great",
  "banana": [
    {
      "cobra": "kai",
      "angus": "bob"
    }
  ]
}`
	var actualJSON = yamlToJSON(t, sampleYaml, 2)
	test.AssertResult(t, expectedJSON, actualJSON)
}

func TestJsonNullInArray(t *testing.T) {
	var sampleYaml = `[null]`
	var actualJSON = yamlToJSON(t, sampleYaml, 0)
	test.AssertResult(t, sampleYaml, actualJSON)
}

func TestJsonNull(t *testing.T) {
	var sampleYaml = `null`
	var actualJSON = yamlToJSON(t, sampleYaml, 0)
	test.AssertResult(t, sampleYaml, actualJSON)
}

func TestJsonNullInObject(t *testing.T) {
	var sampleYaml = `{x: null}`
	var actualJSON = yamlToJSON(t, sampleYaml, 0)
	test.AssertResult(t, `{"x":null}`, actualJSON)
}

func TestJsonEncoderDoesNotEscapeHTMLChars(t *testing.T) {
	var sampleYaml = `build: "( ./lint && ./format && ./compile ) < src.code"`
	var expectedJSON = `{"build":"( ./lint && ./format && ./compile ) < src.code"}`
	var actualJSON = yamlToJSON(t, sampleYaml, 0)
	test.AssertResult(t, expectedJSON, actualJSON)
}

// Test cases added by Mykhailo Isyp
func yamlToJSONWithPrefs(t *testing.T, sampleYaml string, prefs JsonPreferences) string {
	t.Helper()
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	var jsonEncoder = NewJSONEncoder(prefs)
	inputs, err := readDocuments(strings.NewReader(sampleYaml), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}
	node := inputs.Front().Value.(*CandidateNode)

	err = jsonEncoder.Encode(writer, node)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	return strings.TrimSuffix(output.String(), "\n")
}

func yamlToYAMLWithPrefs(t *testing.T, sampleYaml string, prefs YamlPreferences) string {
	t.Helper()
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)

	var yamlEncoder = NewYamlEncoder(prefs)
	inputs, err := readDocuments(strings.NewReader(sampleYaml), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		panic(err)
	}
	node := inputs.Front().Value.(*CandidateNode)

	err = yamlEncoder.Encode(writer, node)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	return strings.TrimSuffix(output.String(), "\n")
}

// JSON tests by Mykhailo Isyp
func TestJsonEncoderUnwrapScalarTrue(t *testing.T) {
	prefs := ConfiguredJSONPreferences.Copy()
	prefs.UnwrapScalar = true
	prefs.Indent = 2

	actual := yamlToJSONWithPrefs(t, "cat", prefs)
	test.AssertResult(t, "cat", actual)
}

func TestJsonEncoderUnwrapScalarFalse(t *testing.T) {
	prefs := ConfiguredJSONPreferences.Copy()
	prefs.UnwrapScalar = false
	prefs.Indent = 2

	actual := yamlToJSONWithPrefs(t, "cat", prefs)
	test.AssertResult(t, `"cat"`, actual)
}

func TestJsonEncoderIndent4(t *testing.T) {
	prefs := ConfiguredJSONPreferences.Copy()
	prefs.UnwrapScalar = false
	prefs.Indent = 4

	actual := yamlToJSONWithPrefs(t, "a: 1\nb: 2", prefs)
	expected := "{\n    \"a\": 1,\n    \"b\": 2\n}"
	test.AssertResult(t, expected, actual)
}

type failingWriter struct{}

func (f failingWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write failed")
}

type failOnSecondWriteWriter struct {
	writes int
	buf    bytes.Buffer
}

func (w *failOnSecondWriteWriter) Write(p []byte) (int, error) {
	w.writes++
	if w.writes >= 2 {
		return 0, errors.New("second write failed")
	}
	return w.buf.Write(p)
}

func TestJsonEncoderReturnsErrorOnWriteFailure(t *testing.T) {
	prefs := ConfiguredJSONPreferences.Copy()
	prefs.UnwrapScalar = false

	enc := NewJSONEncoder(prefs)
	inputs, err := readDocuments(strings.NewReader("a: 1"), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		t.Fatal(err)
	}
	node := inputs.Front().Value.(*CandidateNode)

	err = enc.Encode(failingWriter{}, node)
	if err == nil {
		t.Fatal("expected write error")
	}
	test.AssertResultComplex(t, "write failed", err.Error())
}

func TestJsonEncoderColorsEnabledIncludesANSI(t *testing.T) {
	prefs := ConfiguredJSONPreferences.Copy()
	prefs.ColorsEnabled = true
	prefs.UnwrapScalar = false
	prefs.Indent = 0

	actual := yamlToJSONWithPrefs(t, "a: 1", prefs)

	if !strings.Contains(actual, "\x1b[") {
		t.Fatalf("expected ANSI color codes, got %q", actual)
	}
	if !strings.Contains(actual, `"a"`) {
		t.Fatalf("expected encoded key in output, got %q", actual)
	}
}

// YAML tests by Mykhailo Isyp
func TestYamlEncoderUnwrapScalarTrue(t *testing.T) {
	prefs := ConfiguredYamlPreferences.Copy()
	prefs.UnwrapScalar = true

	actual := yamlToYAMLWithPrefs(t, "cat", prefs)
	test.AssertResult(t, "cat", actual)
}

func TestYamlEncoderUnwrapScalarFalse(t *testing.T) {
	prefs := ConfiguredYamlPreferences.Copy()
	prefs.UnwrapScalar = false

	actual := yamlToYAMLWithPrefs(t, "cat", prefs)
	test.AssertResult(t, "cat", actual)
}

func TestYamlEncoderCompactSequenceIndent(t *testing.T) {
	prefs := ConfiguredYamlPreferences.Copy()
	prefs.CompactSequenceIndent = true
	prefs.Indent = 2

	actual := yamlToYAMLWithPrefs(t, "a:\n  - x\n  - y", prefs)
	test.AssertResult(t, "a:\n- x\n- y", actual)
}

func TestYamlEncoderColorsEnabledIncludesANSI(t *testing.T) {
	prefs := ConfiguredYamlPreferences.Copy()
	prefs.ColorsEnabled = true
	prefs.Indent = 2

	actual := yamlToYAMLWithPrefs(t, "a: 1", prefs)

	if !strings.Contains(actual, "\x1b[") {
		t.Fatalf("expected ANSI color codes, got %q", actual)
	}
	if !strings.Contains(actual, "a") {
		t.Fatalf("expected key content in output, got %q", actual)
	}
	if !strings.Contains(actual, "1") {
		t.Fatalf("expected value content in output, got %q", actual)
	}
}

func TestYamlEncoderEmptyStringNotUnwrapped(t *testing.T) {
	prefs := ConfiguredYamlPreferences.Copy()
	prefs.UnwrapScalar = false

	actual := yamlToYAMLWithPrefs(t, `""`, prefs)
	test.AssertResult(t, `""`, actual)
}

func TestYamlEncoderReturnsErrorOnEncodeWriteFailure(t *testing.T) {
	prefs := ConfiguredYamlPreferences.Copy()
	enc := NewYamlEncoder(prefs)

	inputs, err := readDocuments(strings.NewReader("a: 1"), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		t.Fatal(err)
	}
	node := inputs.Front().Value.(*CandidateNode)

	err = enc.Encode(failingWriter{}, node)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestYamlEncoderReturnsErrorOnTrailingContentWriteFailure(t *testing.T) {
	prefs := ConfiguredYamlPreferences.Copy()
	enc := NewYamlEncoder(prefs)

	// comment after document to force trailing/leading content printing
	input := "a: 1\n# trailing comment\n"
	inputs, err := readDocuments(strings.NewReader(input), "sample.yml", 0, NewYamlDecoder(ConfiguredYamlPreferences))
	if err != nil {
		t.Fatal(err)
	}
	node := inputs.Front().Value.(*CandidateNode)

	w := &failOnSecondWriteWriter{}
	err = enc.Encode(w, node)
	if err == nil {
		t.Fatal("expected error")
	}
}

// End of test cases added by Mykhailo Isyp
