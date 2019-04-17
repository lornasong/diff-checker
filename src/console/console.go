package console

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lornasong/diff-checker/src/compare"
)

type Printer struct {
	data      []*compare.Matcher
	plusColor func(a ...interface{}) string
	diffColor func(a ...interface{}) string
	aColor    func(a ...interface{}) string
	bColor    func(a ...interface{}) string
}

func NewPrinter(data []*compare.Matcher, opts ...func(*Printer)) *Printer {
	p := &Printer{
		data:      data,
		plusColor: color.New(color.FgGreen).SprintFunc(),
		diffColor: color.New(color.FgHiYellow).SprintFunc(),
		aColor:    color.New(color.FgCyan).SprintFunc(),
		bColor:    color.New(color.FgHiMagenta).SprintFunc(),
	}

	return p
}

func (p *Printer) Diff() {
	p.printHeader()
	defer p.printFooter()

	if len(p.data) == 0 {
		p.printNoDiffMsg()
		return
	}

	lineNumB := 0
	lineNumA := 0
	for _, line := range p.data {
		if line.Same() {
			lineNumA++
			lineNumB++
			p.printSameLine(line, lineNumA, lineNumB)
			continue
		}

		if line.Similar() {
			lineNumA++
			lineNumB++
			p.printSimilarLine(line, lineNumA, lineNumB)
			continue
		}

		if line.OnlyInA() {
			lineNumA++
			p.printOnlyALine(line, lineNumA)
		}

		if line.OnlyInB() {
			lineNumB++
			p.printOnlyBLine(line, lineNumB)
		}
	}
}

func (p *Printer) printSameLine(line *compare.Matcher, lineNumA, lineNumB int) {
	fmt.Printf("L%d/%d A/B\t:%s\n", lineNumA, lineNumB, line.A())
}

func (p *Printer) printSimilarLine(line *compare.Matcher, lineNumA, lineNumB int) {
	fmt.Printf("L%d/%d %sA/B\t:", lineNumA, lineNumB, p.diffColor("~"))

	for _, w := range line.Children() {
		if w.Same() {
			fmt.Printf(w.A())
		}
		if w.OnlyInA() {
			fmt.Printf("%s", p.aColor(w.A()))
		}
		if w.OnlyInB() {
			fmt.Printf("%s", p.bColor(w.B()))
		}
	}
	fmt.Println()
}

func (p *Printer) printOnlyALine(line *compare.Matcher, lineNumA int) {
	fmt.Printf("L%d %sA\t\t:%s\n", lineNumA, p.plusColor("+"), p.aColor(line.A()))
}

func (p *Printer) printOnlyBLine(line *compare.Matcher, lineNumB int) {
	fmt.Printf("L%d %sB\t\t:%s\n", lineNumB, p.plusColor("+"), p.bColor(line.B()))
}

func (p *Printer) printHeader() {
	printLine()
	fmt.Printf("Diff %s & %s:\n", p.aColor("A"), p.bColor("B"))
	printLine()
}

func (p *Printer) printFooter() {
	printLine()
	fmt.Println("")
	printLine()
}

func printLine() {
	fmt.Println("---------------------------------------------------------------------")
}

func (p *Printer) printNoDiffMsg() {
	fmt.Println("No differences between files")
}
