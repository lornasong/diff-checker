package main

import (
	"fmt"

	"github.com/lornasong/diff-checker/src/compare"
)

func main() {

	a := getA()
	b := getB()

	lines := compare.Match(a, b)

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
			fmt.Printf("Line %d: +A %s\n", lineNumA, line.A())
		}
		if line.OnlyInB() {
			lineNumB++
			fmt.Printf("Line %d: +B %s\n", lineNumB, line.B())
		}
	}
	fmt.Println("End Diff-------")
}

func getA() string {
	return `
	`
}

func getB() string {
	return `
	`
}
