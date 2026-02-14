package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "    123, counting  ",
			expected: []string{"123,", "counting"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		expectedLen := len(c.expected)
		if len(actual) != expectedLen {
			t.Errorf("string slices lengths do not match")

		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("strings words do not match")
			}

		}
	}
}
