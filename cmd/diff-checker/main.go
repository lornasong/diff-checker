package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	_ "net/http/pprof"

	"github.com/fatih/color"
	"github.com/lornasong/diff-checker/src/compare"
	"github.com/lornasong/diff-checker/src/console"
	"github.com/pkg/profile"
)

func main() {
	fmt.Println("change one")
	fmt.Println("change two")
	defer profile.Start(profile.MemProfile).Stop()

	// have to do some research here. this doesn't seem to be the right way in combination with make
	aColorStr := flag.String("before-color", "cyan", "color to display string 'a'")

	colors := make(map[string]color.Attribute)
	colors["black"] = color.FgHiBlack
	colors["red"] = color.FgHiRed
	colors["green"] = color.FgHiGreen
	colors["yellow"] = color.FgHiYellow
	colors["blue"] = color.FgHiBlue
	colors["magenta"] = color.FgHiMagenta
	colors["cyan"] = color.FgHiCyan
	colors["white"] = color.FgHiWhite

	aColor := color.FgHiMagenta
	if c, ok := colors[*aColorStr]; ok {
		aColor = c
	} else {
		log.Fatalln("warning. invalid color")
	}

	a, err := ioutil.ReadFile("a.txt")
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile("b.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := compare.MatchLine(string(a), string(b))
	console.NewPrinter(lines, console.WithAColorAttribute(aColor)).Diff()
}
