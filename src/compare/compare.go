package compare

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	newLineDelimiter    = "\n"
	spaceDelimiter      = ' '
	dashDelimiter       = '-'
	underScoreDelimiter = '_'
	commaDelimiter      = ','
	colonDelimiter      = ':'
	quotationDelimiter  = '"'
)

// Matcher contains text from string a that has been matched to string
// b up until newline. No match is represented by an empty string
type Matcher struct {
	a            string
	b            string
	childMatches []*Matcher
}

// NewMatcher returns a new Matcher
func NewMatcher(a, b string, opts ...func(*Matcher)) *Matcher {
	m := &Matcher{
		a: a,
		b: b,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// WithChildMatches sets the child matches. Without this option, the matcher
// has no child matches
func WithChildMatches(child []*Matcher) func(*Matcher) {
	return func(m *Matcher) { m.childMatches = child }
}

// A returns the value of string a
func (p *Matcher) A() string {
	return p.a
}

// B returns the value of string b
func (p *Matcher) B() string {
	return p.b
}

// Children returns the array of the children
func (p *Matcher) Children() []*Matcher {
	return p.childMatches
}

// Same returns whether or not the contents of a and b are the same
// for this line
func (p *Matcher) Same() bool {
	return p.a == p.b
}

// Similar returns whether the line is considered similar between a and b
func (p *Matcher) Similar() bool {
	return len(p.childMatches) > 0
}

// OnlyInA returns whether or not this line only occurs in a
func (p *Matcher) OnlyInA() bool {
	return len(p.a) > 0 && len(p.b) == 0
}

// OnlyInB returns whether or not this line only occurs in b
func (p *Matcher) OnlyInB() bool {
	return len(p.b) > 0 && len(p.a) == 0
}

// MatchLine returns a list of Matchers that represent the matches between string
// a and string b if any.
func MatchLine(a, b string) []*Matcher {
	aPieces := strings.Split(a, "\n")
	bPieces := strings.Split(b, "\n")

	matches := make([]*Matcher, 0)
	unmatchedB := make([]*Matcher, 0)
	matched := false

	for _, a := range aPieces {
		matched = false
		unmatchedB = nil

		// look for matches to b
		for ixb, b := range bPieces {

			children := sameSimilarLine(a, b)
			if children != nil {
				if len(unmatchedB) > 0 {
					matches = append(matches, unmatchedB...)
					unmatchedB = nil
				}

				matched = true
				matches = append(matches, NewMatcher(a, b, WithChildMatches(children)))

				// trim off added stuff
				bPieces = bPieces[ixb+1:]
				break
			}

			unmatchedB = append(unmatchedB, NewMatcher("", b))

		}
		// there was no match to b
		if !matched {
			matches = append(matches, NewMatcher(a, ""))
		}
	}

	// incomplete b
	for _, b := range bPieces {
		matches = append(matches, NewMatcher("", b))
	}

	return matches
}

func sameSimilarLine(a, b string) []*Matcher {
	fmt.Println("matching lines")
	fmt.Println("a: ", a)
	fmt.Println("b: ", b)

	words := matchWords(a, b)
	if len(words) == 0 {
		fmt.Println("not similar (1)")
		return nil
	}

	if similar(words, 0.5) {
		return words
	}
	return nil
}

func similar(matches []*Matcher, threshold float64) bool {
	i := 0.0
	for _, childMatch := range matches {
		if childMatch.Same() {
			i++
		}
	}
	percent := i / float64(len(matches))
	if percent > threshold {
		fmt.Println("similar", i, len(matches), percent)
		return true
	}
	fmt.Println("not similar (2)", i, len(matches), percent)
	return false
}

func splitLine(s string) []string {
	delims := make(map[rune]bool)
	delims[spaceDelimiter] = true
	delims[dashDelimiter] = true
	delims[underScoreDelimiter] = true
	delims[commaDelimiter] = true
	delims[colonDelimiter] = true
	delims[quotationDelimiter] = true

	var buf bytes.Buffer
	pieces := make([]string, 0)

	for _, r := range s {
		if _, ok := delims[r]; !ok {
			buf.WriteRune(r)
			continue
		}

		pieces = append(pieces, buf.String())
		buf.Reset()
		pieces = append(pieces, string(r))
	}
	if len(buf.String()) > 0 {
		pieces = append(pieces, buf.String())
	}

	return pieces
}

func matchWords(a, b string) []*Matcher {

	aPieces := splitLine(a)
	bPieces := splitLine(b)

	matches := make([]*Matcher, 0)
	unmatchedB := make([]*Matcher, 0)
	matched := false

	for _, a := range aPieces {
		matched = false
		unmatchedB = nil

		// look for matches to b
		for ixb, b := range bPieces {

			if a == b {
				if len(unmatchedB) > 0 {
					matches = append(matches, unmatchedB...)
					unmatchedB = nil
				}

				matched = true
				matches = append(matches, NewMatcher(a, b))

				// trim off added stuff
				bPieces = bPieces[ixb+1:]
				break
			}

			unmatchedB = append(unmatchedB, NewMatcher("", b))
		}

		// there was no match to b
		if !matched {
			matches = append(matches, NewMatcher(a, ""))
		}
	}

	// incomplete b
	for _, b := range bPieces {
		matches = append(matches, NewMatcher("", b))
	}

	return matches
}
