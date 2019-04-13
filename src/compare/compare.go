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

// LineMatch contains text from string a that has been matched to string
// b up until newline. No match is represented by an empty string
type LineMatch struct {
	a            string
	b            string
	childMatches []*LineMatch
}

// NewLineMatch returns a new LineMatch
func NewLineMatch(a, b string, opts ...func(*LineMatch)) *LineMatch {
	m := &LineMatch{
		a: a,
		b: b,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func WithChildMatches(child []*LineMatch) func(*LineMatch) {
	return func(m *LineMatch) { m.childMatches = child }
}

// A returns the value of string a
func (p *LineMatch) A() string {
	return p.a
}

// B returns the value of string b
func (p *LineMatch) B() string {
	return p.b
}

// A returns the value of string a
func (p *LineMatch) Children() []*LineMatch {
	return p.childMatches
}

// Same returns whether or not the contents of a and b are the same
// for this line
func (p *LineMatch) Same() bool {
	return p.a == p.b
}

// TODO:
func (p *LineMatch) Similar() bool {
	return len(p.childMatches) > 0
}

// OnlyInA returns whether or not this line only occurs in a
func (p *LineMatch) OnlyInA() bool {
	return len(p.a) > 0 && len(p.b) == 0
}

// OnlyInB returns whether or not this line only occurs in b
func (p *LineMatch) OnlyInB() bool {
	return len(p.b) > 0 && len(p.a) == 0
}

// MatchLine returns a list of LineMatchs that represent the matches between string
// a and string b if any.
func MatchLine(a, b string) []*LineMatch {
	aPieces := strings.Split(a, "\n")
	bPieces := strings.Split(b, "\n")

	matches := make([]*LineMatch, 0)
	unmatchedB := make([]*LineMatch, 0)
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
				matches = append(matches, NewLineMatch(a, b, WithChildMatches(children)))

				// trim off added stuff
				bPieces = bPieces[ixb+1:]
				break
			}

			unmatchedB = append(unmatchedB, NewLineMatch("", b))

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

func sameSimilarLine(a, b string) []*LineMatch {
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

func similar(matches []*LineMatch, threshold float64) bool {
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

func matchWords(a, b string) []*LineMatch {

	aPieces := splitLine(a)
	bPieces := splitLine(b)

	matches := make([]*LineMatch, 0)
	unmatchedB := make([]*LineMatch, 0)
	matched := false

	// TODO: refactor. think through any optimizations
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
				matches = append(matches, NewLineMatch(a, b))

				// trim off added stuff
				bPieces = bPieces[ixb+1:]
				break
			}

			unmatchedB = append(unmatchedB, NewLineMatch("", b))
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
