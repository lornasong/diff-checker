package main

import (
	"fmt"
	"io/ioutil"
	"log"
	_ "net/http/pprof"

	"github.com/fatih/color"
	"github.com/lornasong/diff-checker/src/compare"
	"github.com/pkg/profile"
)

const pathToInputFiles = ""

func main() {
	defer profile.Start(profile.MemProfile).Stop()

	plusColor := color.New(color.FgGreen).SprintFunc()
	diffColor := color.New(color.FgHiYellow).SprintFunc()

	aColor := color.New(color.FgCyan).SprintFunc()
	bColor := color.New(color.FgHiMagenta).SprintFunc()

	a, err := ioutil.ReadFile("a.txt")
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile("b.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := compare.MatchLine(string(a), string(b))
	if len(lines) == 0 {
		fmt.Println("No differences between files")
		return
	}

	lineNumB := 0
	lineNumA := 0
	fmt.Println("---------------------------------------------------------------------")
	fmt.Printf("Diff %s & %s:\n", aColor("A"), bColor("B"))
	fmt.Println("---------------------------------------------------------------------")
	for _, line := range lines {
		if line.Same() {
			lineNumA++
			lineNumB++
			fmt.Printf("L%d/%d A/B\t:%s\n", lineNumA, lineNumB, line.A())
			continue
		}

		if line.Similar() {
			lineNumA++
			lineNumB++
			fmt.Printf("L%d/%d %sA/B\t:", lineNumA, lineNumB, diffColor("~"))
			printSimilarLine(line, aColor, bColor)
			continue
		}

		if line.OnlyInA() {
			lineNumA++
			fmt.Printf("L%d %sA\t\t:%s\n", lineNumA, plusColor("+"), aColor(line.A()))
		}
		if line.OnlyInB() {
			lineNumB++
			fmt.Printf("L%d %sB\t\t:%s\n", lineNumB, plusColor("+"), bColor(line.B()))
		}
	}
	fmt.Println("---------------------------------------------------------------------")
	fmt.Println("")
	fmt.Println("---------------------------------------------------------------------")
}

func printSimilarLine(line *compare.LineMatch, aColor, bColor func(a ...interface{}) string) {
	for _, w := range line.Children() {
		if w.Same() {
			fmt.Printf(w.A())
		}
		if w.OnlyInA() {
			fmt.Printf("%s", aColor(w.A()))
		}
		if w.OnlyInB() {
			fmt.Printf("%s", bColor(w.B()))
		}
	}
	fmt.Println()
}
