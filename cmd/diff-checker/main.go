package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/lornasong/diff-checker/src/compare"
)

const pathToInputFiles = ""

func main() {
	a, err := ioutil.ReadFile("a.txt")
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile("b.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := compare.Match(string(a), string(b))
	if len(lines) == 0 {
		fmt.Println("No differences between files")
		return
	}

	lineNumB := 0
	lineNumA := 0
	fmt.Println("Check Diff-----")
	for _, line := range lines {
		if line.Same() {
			lineNumA++
			lineNumB++
			continue
		}

		if line.OnlyInA() {
			lineNumA++
			fmt.Printf("  Line %d: +A %s\n", lineNumA, line.A())
		}
		if line.OnlyInB() {
			lineNumB++
			fmt.Printf("  Line %d: +B %s\n", lineNumB, line.B())
		}
	}
	fmt.Println("End Diff-------")
}
