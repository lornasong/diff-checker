package compare_test

import (
	"fmt"
	"testing"

	comp "github.com/lornasong/diff-checker/src/compare"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchLine_Same(t *testing.T) {
	cases := []struct {
		scenario string
		str      string
		expected []*comp.LineMatch
	}{
		{
			scenario: "Same, No Newline",
			str:      "abcd",
			expected: []*comp.LineMatch{
				comp.NewLineMatch("abcd", "abcd"),
			},
		},
		{
			scenario: "Same, With Single Newline, Basic",
			str:      "abcd\n1234",
			expected: []*comp.LineMatch{
				comp.NewLineMatch("abcd", "abcd"),
				comp.NewLineMatch("1234", "1234"),
			},
		},
		{
			scenario: "Same, With Multi Newline",
			str:      "abcd\n1234\nqwerty",
			expected: []*comp.LineMatch{
				comp.NewLineMatch("abcd", "abcd"),
				comp.NewLineMatch("1234", "1234"),
				comp.NewLineMatch("qwerty", "qwerty"),
			},
		},
		{
			scenario: "Same, With Single Newline, Consecutive Newlines",
			str:      "abcd\n\n1234",
			expected: []*comp.LineMatch{
				comp.NewLineMatch("abcd", "abcd"),
				comp.NewLineMatch("", ""),
				comp.NewLineMatch("1234", "1234"),
			},
		},
		{
			scenario: "Same, With Multi Newline, Complex",
			str:      "{\n\"key\": \"value\",\n\"stuff\": \"more stuff\"\n}",
			expected: []*comp.LineMatch{
				comp.NewLineMatch("{", "{"),
				comp.NewLineMatch(`"key": "value",`, `"key": "value",`),
				comp.NewLineMatch(`"stuff": "more stuff"`, `"stuff": "more stuff"`),
				comp.NewLineMatch("}", "}"),
			},
		},
	}

	for _, tc := range cases {
		fmt.Println("Running test for scenario: ", tc.scenario)
		actual := comp.MatchLine(tc.str, tc.str)
		require.Equal(t, len(tc.expected), len(actual), "Failed: "+tc.scenario)

		for ix, p := range tc.expected {
			assert.Equal(t, p.A(), actual[ix].A(), "Failed: "+tc.scenario)
			assert.Equal(t, p.B(), actual[ix].B(), "Failed: "+tc.scenario)
		}
	}
}

func TestMatchLine_Diff(t *testing.T) {
	cases := []struct {
		scenario string
		a        string
		b        string
		expected []*comp.LineMatch
	}{
		{
			scenario: "Different, No Newline, Completely",
			a:        "abcd",
			b:        "1234",
			expected: []*comp.LineMatch{
				comp.NewLineMatch("abcd", ""),
				comp.NewLineMatch("", "1234"),
			},
		},
		{
			scenario: "Different, No Newline, Partially",
			a:        "abcd",
			b:        "abc4",
			expected: []*comp.LineMatch{
				comp.NewLineMatch("abcd", ""),
				comp.NewLineMatch("", "abc4"),
			},
		},
		{
			scenario: "Different, With Newline, First Line",
			a:        "1234\nxyz",
			b:        "abcd\nxyz",
			expected: []*comp.LineMatch{
				comp.NewLineMatch("1234", ""),
				comp.NewLineMatch("", "abcd"),
				comp.NewLineMatch("xyz", "xyz"),
			},
		},
		{
			scenario: "Different, With Newline, Mid Line",
			a:        "abcd\n1234\nxyz",
			b:        "abcd\n5678\nxyz",
			expected: []*comp.LineMatch{
				comp.NewLineMatch("abcd", "abcd"),
				comp.NewLineMatch("1234", ""),
				comp.NewLineMatch("", "5678"),
				comp.NewLineMatch("xyz", "xyz"),
			},
		},
		{
			scenario: "Different, With Newline, Last Line",
			a:        "abcd\n1234",
			b:        "abcd\nxyz",
			expected: []*comp.LineMatch{
				comp.NewLineMatch("abcd", "abcd"),
				comp.NewLineMatch("1234", ""),
				comp.NewLineMatch("", "xyz"),
			},
		},
		{
			scenario: "Different, With Newline, Partial Criss-Cross Match",
			a:        "abcd\n1234",
			b:        "1234\nxyz",
			expected: []*comp.LineMatch{
				comp.NewLineMatch("abcd", ""),
				comp.NewLineMatch("1234", "1234"),
				comp.NewLineMatch("", "xyz"),
			},
		},
		{
			scenario: "Different, With Newline, Complete Criss-Cross Match",
			a:        "abcd\n1234",
			b:        "1234\nabcd",
			expected: []*comp.LineMatch{
				comp.NewLineMatch("", "1234"),
				comp.NewLineMatch("abcd", "abcd"),
				comp.NewLineMatch("1234", ""),
			},
		},
	}

	for _, tc := range cases {
		fmt.Println("Running test for scenario: ", tc.scenario)
		actual := comp.MatchLine(tc.a, tc.b)
		require.Equal(t, len(tc.expected), len(actual), "Failed: "+tc.scenario)

		for ix, p := range tc.expected {
			assert.Equal(t, p.A(), actual[ix].A(), "Failed: "+tc.scenario)
			assert.Equal(t, p.B(), actual[ix].B(), "Failed: "+tc.scenario)
		}
	}
}
