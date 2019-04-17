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
		expected []*Matcher
	}{
		{
			scenario: "Same, No Newline",
			str:      "abcd",
			expected: []*Matcher{
				NewMatcher("abcd", "abcd"),
			},
		},
		{
			scenario: "Same, With Single Newline, Basic",
			str:      "abcd\n1234",
			expected: []*Matcher{
				NewMatcher("abcd", "abcd"),
				NewMatcher("1234", "1234"),
			},
		},
		{
			scenario: "Same, With Multi Newline",
			str:      "abcd\n1234\nqwerty",
			expected: []*Matcher{
				NewMatcher("abcd", "abcd"),
				NewMatcher("1234", "1234"),
				NewMatcher("qwerty", "qwerty"),
			},
		},
		{
			scenario: "Same, With Single Newline, Consecutive Newlines",
			str:      "abcd\n\n1234",
			expected: []*Matcher{
				NewMatcher("abcd", "abcd"),
				NewMatcher("", ""),
				NewMatcher("1234", "1234"),
			},
		},
		{
			scenario: "Same, With Multi Newline, Complex",
			str:      "{\n\"key\": \"value\",\n\"stuff\": \"more stuff\"\n}",
			expected: []*Matcher{
				NewMatcher("{", "{"),
				NewMatcher(`"key": "value",`, `"key": "value",`),
				NewMatcher(`"stuff": "more stuff"`, `"stuff": "more stuff"`),
				NewMatcher("}", "}"),
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
		expected []*Matcher
	}{
		{
			scenario: "Different, No Newline, Completely",
			a:        "abcd",
			b:        "1234",
			expected: []*Matcher{
				NewMatcher("abcd", ""),
				NewMatcher("", "1234"),
			},
		},
		{
			scenario: "Different, No Newline, Partially",
			a:        "abcd",
			b:        "abc4",
			expected: []*Matcher{
				NewMatcher("abcd", ""),
				NewMatcher("", "abc4"),
			},
		},
		{
			scenario: "Different, With Newline, First Line",
			a:        "1234\nxyz",
			b:        "abcd\nxyz",
			expected: []*Matcher{
				NewMatcher("1234", ""),
				NewMatcher("", "abcd"),
				NewMatcher("xyz", "xyz"),
			},
		},
		{
			scenario: "Different, With Newline, Mid Line",
			a:        "abcd\n1234\nxyz",
			b:        "abcd\n5678\nxyz",
			expected: []*Matcher{
				NewMatcher("abcd", "abcd"),
				NewMatcher("1234", ""),
				NewMatcher("", "5678"),
				NewMatcher("xyz", "xyz"),
			},
		},
		{
			scenario: "Different, With Newline, Last Line",
			a:        "abcd\n1234",
			b:        "abcd\nxyz",
			expected: []*Matcher{
				NewMatcher("abcd", "abcd"),
				NewMatcher("1234", ""),
				NewMatcher("", "xyz"),
			},
		},
		{
			scenario: "Different, With Newline, Partial Criss-Cross Match",
			a:        "abcd\n1234",
			b:        "1234\nxyz",
			expected: []*Matcher{
				NewMatcher("abcd", ""),
				NewMatcher("1234", "1234"),
				NewMatcher("", "xyz"),
			},
		},
		{
			scenario: "Different, With Newline, Complete Criss-Cross Match",
			a:        "abcd\n1234",
			b:        "1234\nabcd",
			expected: []*Matcher{
				NewMatcher("", "1234"),
				NewMatcher("abcd", "abcd"),
				NewMatcher("1234", ""),
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
		matches   []*Matcher
		threshold float64
		expected  bool
	}{
		{
			scenario: "Complete Same",
			matches: []*Matcher{
				NewMatcher("same", "same"),
				NewMatcher("same2", "same2"),
			},
			threshold: 0.5,
			expected:  true,
		},
		{
			scenario: "Complete Different",
			matches: []*Matcher{
				NewMatcher("diff1", ""),
				NewMatcher("", "diff2"),
			},
			threshold: 0.5,
			expected:  false,
		},
		{
			scenario: "Percent = Threshold",
			matches: []*Matcher{
				NewMatcher("same", "same"),
				NewMatcher("diff2", "diff1"),
			},
			threshold: 0.5,
			expected:  false,
		},
		{
			scenario: "Percent >= Threshold",
			matches: []*Matcher{
				NewMatcher("same", "same"),
				NewMatcher("diff2", "diff1"),
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
		expected []*Matcher
	}{
		{
			scenario: "Same line",
			a:        "environ:dev",
			b:        "environ:dev",
			expected: []*Matcher{
				NewMatcher("environ", "environ"),
				NewMatcher(":", ":"),
				NewMatcher("dev", "dev"),
			},
		},
		{
			scenario: "a match",
			a:        "environ:dev",
			b:        "environ:staging",
			expected: []*Matcher{
				NewMatcher("environ", "environ"),
				NewMatcher(":", ":"),
				NewMatcher("dev", ""),
				NewMatcher("", "staging"),
			},
		},
		{
			scenario: "some matching",
			a:        "environ-dev",
			b:        "environ:staging",
			expected: []*Matcher{
				NewMatcher("environ", "environ"),
				NewMatcher("-", ""),
				NewMatcher("dev", ""),
				NewMatcher("", ":"),
				NewMatcher("", "staging"),
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
