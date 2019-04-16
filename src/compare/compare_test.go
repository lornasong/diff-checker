package compare

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchLine_Same(t *testing.T) {
	cases := []struct {
		scenario string
		str      string
		expected []*LineMatch
	}{
		{
			scenario: "Same, No Newline",
			str:      "abcd",
			expected: []*LineMatch{
				NewLineMatch("abcd", "abcd"),
			},
		},
		{
			scenario: "Same, With Single Newline, Basic",
			str:      "abcd\n1234",
			expected: []*LineMatch{
				NewLineMatch("abcd", "abcd"),
				NewLineMatch("1234", "1234"),
			},
		},
		{
			scenario: "Same, With Multi Newline",
			str:      "abcd\n1234\nqwerty",
			expected: []*LineMatch{
				NewLineMatch("abcd", "abcd"),
				NewLineMatch("1234", "1234"),
				NewLineMatch("qwerty", "qwerty"),
			},
		},
		{
			scenario: "Same, With Single Newline, Consecutive Newlines",
			str:      "abcd\n\n1234",
			expected: []*LineMatch{
				NewLineMatch("abcd", "abcd"),
				NewLineMatch("", ""),
				NewLineMatch("1234", "1234"),
			},
		},
		{
			scenario: "Same, With Multi Newline, Complex",
			str:      "{\n\"key\": \"value\",\n\"stuff\": \"more stuff\"\n}",
			expected: []*LineMatch{
				NewLineMatch("{", "{"),
				NewLineMatch(`"key": "value",`, `"key": "value",`),
				NewLineMatch(`"stuff": "more stuff"`, `"stuff": "more stuff"`),
				NewLineMatch("}", "}"),
			},
		},
	}

	for _, tc := range cases {
		fmt.Println("Running test for scenario: ", tc.scenario)
		actual := MatchLine(tc.str, tc.str)
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
		expected []*LineMatch
	}{
		{
			scenario: "Different, No Newline, Completely",
			a:        "abcd",
			b:        "1234",
			expected: []*LineMatch{
				NewLineMatch("abcd", ""),
				NewLineMatch("", "1234"),
			},
		},
		{
			scenario: "Different, No Newline, Partially",
			a:        "abcd",
			b:        "abc4",
			expected: []*LineMatch{
				NewLineMatch("abcd", ""),
				NewLineMatch("", "abc4"),
			},
		},
		{
			scenario: "Different, With Newline, First Line",
			a:        "1234\nxyz",
			b:        "abcd\nxyz",
			expected: []*LineMatch{
				NewLineMatch("1234", ""),
				NewLineMatch("", "abcd"),
				NewLineMatch("xyz", "xyz"),
			},
		},
		{
			scenario: "Different, With Newline, Mid Line",
			a:        "abcd\n1234\nxyz",
			b:        "abcd\n5678\nxyz",
			expected: []*LineMatch{
				NewLineMatch("abcd", "abcd"),
				NewLineMatch("1234", ""),
				NewLineMatch("", "5678"),
				NewLineMatch("xyz", "xyz"),
			},
		},
		{
			scenario: "Different, With Newline, Last Line",
			a:        "abcd\n1234",
			b:        "abcd\nxyz",
			expected: []*LineMatch{
				NewLineMatch("abcd", "abcd"),
				NewLineMatch("1234", ""),
				NewLineMatch("", "xyz"),
			},
		},
		{
			scenario: "Different, With Newline, Partial Criss-Cross Match",
			a:        "abcd\n1234",
			b:        "1234\nxyz",
			expected: []*LineMatch{
				NewLineMatch("abcd", ""),
				NewLineMatch("1234", "1234"),
				NewLineMatch("", "xyz"),
			},
		},
		{
			scenario: "Different, With Newline, Complete Criss-Cross Match",
			a:        "abcd\n1234",
			b:        "1234\nabcd",
			expected: []*LineMatch{
				NewLineMatch("", "1234"),
				NewLineMatch("abcd", "abcd"),
				NewLineMatch("1234", ""),
			},
		},
	}

	for _, tc := range cases {
		fmt.Println("Running test for scenario: ", tc.scenario)
		actual := MatchLine(tc.a, tc.b)
		require.Equal(t, len(tc.expected), len(actual), "Failed: "+tc.scenario)

		for ix, p := range tc.expected {
			assert.Equal(t, p.A(), actual[ix].A(), "Failed: "+tc.scenario)
			assert.Equal(t, p.B(), actual[ix].B(), "Failed: "+tc.scenario)
		}
	}
}

func TestSplitLine(t *testing.T) {
	cases := []struct {
		scenario string
		s        string
		expected []string
	}{
		{
			scenario: "Empty String",
			s:        "",
			expected: []string{},
		},
		{
			scenario: "No split",
			s:        "config",
			expected: []string{"config"},
		},
		{
			scenario: "Split by all word delimiters",
			s:        "space dash-underscore_comma,colon:quotation\"",
			expected: []string{"space", " ", "dash", "-", "underscore", "_", "comma", ",", "colon", ":", "quotation", "\""},
		},
		{
			scenario: "End with a delimiter",
			s:        "key,value,",
			expected: []string{"key", ",", "value", ","},
		},
	}

	for _, tc := range cases {
		actual := splitLine(tc.s)
		require.Equal(t, len(tc.expected), len(actual), "Failed: "+tc.scenario)

		for ix, e := range tc.expected {
			assert.Equal(t, e, actual[ix], "Failed: "+tc.scenario)
		}
	}
}

func TestSimilar(t *testing.T) {
	cases := []struct {
		scenario  string
		matches   []*LineMatch
		threshold float64
		expected  bool
	}{
		{
			scenario: "Complete Same",
			matches: []*LineMatch{
				NewLineMatch("same", "same"),
				NewLineMatch("same2", "same2"),
			},
			threshold: 0.5,
			expected:  true,
		},
		{
			scenario: "Complete Different",
			matches: []*LineMatch{
				NewLineMatch("diff1", ""),
				NewLineMatch("", "diff2"),
			},
			threshold: 0.5,
			expected:  false,
		},
		{
			scenario: "Percent = Threshold",
			matches: []*LineMatch{
				NewLineMatch("same", "same"),
				NewLineMatch("diff2", "diff1"),
			},
			threshold: 0.5,
			expected:  false,
		},
		{
			scenario: "Percent >= Threshold",
			matches: []*LineMatch{
				NewLineMatch("same", "same"),
				NewLineMatch("diff2", "diff1"),
			},
			threshold: 0.49,
			expected:  true,
		},
	}

	for _, tc := range cases {
		actual := similar(tc.matches, tc.threshold)
		assert.Equal(t, tc.expected, actual, "Failed: "+tc.scenario)
	}
}

func TestMatchWord(t *testing.T) {
	cases := []struct {
		scenario string
		a        string
		b        string
		expected []*LineMatch
	}{
		{
			scenario: "Same line",
			a:        "environ:dev",
			b:        "environ:dev",
			expected: []*LineMatch{
				NewLineMatch("environ", "environ"),
				NewLineMatch(":", ":"),
				NewLineMatch("dev", "dev"),
			},
		},
		{
			scenario: "a match",
			a:        "environ:dev",
			b:        "environ:staging",
			expected: []*LineMatch{
				NewLineMatch("environ", "environ"),
				NewLineMatch(":", ":"),
				NewLineMatch("dev", ""),
				NewLineMatch("", "staging"),
			},
		},
		{
			scenario: "some matching",
			a:        "environ-dev",
			b:        "environ:staging",
			expected: []*LineMatch{
				NewLineMatch("environ", "environ"),
				NewLineMatch("-", ""),
				NewLineMatch("dev", ""),
				NewLineMatch("", ":"),
				NewLineMatch("", "staging"),
			},
		},
	}

	for _, tc := range cases {
		fmt.Println("Running test for scenario: ", tc.scenario)
		actual := matchWords(tc.a, tc.b)
		require.Equal(t, len(tc.expected), len(actual), "Failed: "+tc.scenario)

		for ix, p := range tc.expected {
			assert.Equal(t, p.A(), actual[ix].A(), "Failed: "+tc.scenario)
			assert.Equal(t, p.B(), actual[ix].B(), "Failed: "+tc.scenario)
		}
	}
}
