package compare

import (
	"strings"
)

// LineMatch contains text from string a that has been matched to string
// b up until newline. No match is represented by an empty string
type LineMatch struct {
	a string
	b string
}

// NewLineMatch returns a new LineMatch
func NewLineMatch(a, b string) *LineMatch {
	return &LineMatch{
		a: a,
		b: b,
	}
}

// A returns the value of string a
func (p *LineMatch) A() string {
	return p.a
}

// B returns the value of string b
func (p *LineMatch) B() string {
	return p.b
}

// Same returns whether or not the contents of a and b are the same
// for this line
func (p *LineMatch) Same() bool {
	return p.a == p.b
}

// OnlyInA returns whether or not this line only occurs in a
func (p *LineMatch) OnlyInA() bool {
	return len(p.a) > 0 && len(p.b) == 0
}

// OnlyInB returns whether or not this line only occurs in b
func (p *LineMatch) OnlyInB() bool {
	return len(p.b) > 0 && len(p.a) == 0
}

// Match returns a list of LineMatchs that represent the matches between string
// a and string b if any.
func Match(a, b string) []*LineMatch {

	aPieces := strings.Split(a, "\n")
	bPieces := strings.Split(b, "\n")

	matches := make([]*LineMatch, 0)
	unmatchedB := make([]*LineMatch, 0)
	matched := false

	// TODO: refactor. think through any optimizations
	for _, a := range aPieces {
		matched = false
		unmatchedB = nil

		// look for matches to b
		for ixb, b := range bPieces {

			if !sameLine(a, b) {
				unmatchedB = append(unmatchedB, NewLineMatch("", b))
				continue
			}

			if len(unmatchedB) > 0 {
				matches = append(matches, unmatchedB...)
				unmatchedB = nil
			}

			matched = true
			matches = append(matches, NewLineMatch(a, b))

			// trim off added stuff
			bPieces = bPieces[ixb+1:]
			break
		}
		// there was no match to b
		if !matched {
			matches = append(matches, NewLineMatch(a, ""))
		}
	}

	// incomplete b
	for _, b := range bPieces {
		matches = append(matches, NewLineMatch("", b))
	}

	return matches
}

func sameLine(a, b string) bool {
	// grow in complexity to be "similarLine"
	return a == b
}
