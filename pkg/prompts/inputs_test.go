package prompts

import (
	"bytes"
	"fmt"
	"testing"
)

type testCase struct {
	label    string
	expected string
}

func TestStringPrompt(t *testing.T) {
	cases := []testCase{
		{"baz", "baz"},
		{" baz ", "baz"},
		{"baz   ", "baz"},
	}

	var stdin bytes.Buffer

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Label: %q", tc.label), func(t *testing.T) {
			stdin.Write([]byte(tc.label))

			got := StringPrompt("Foobar", &stdin)

			if got != tc.expected {
				t.Errorf("Expected 'baz' but got '%s'", got)
			}
		})
	}
}
