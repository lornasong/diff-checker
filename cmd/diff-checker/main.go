package main

import (
	"io/ioutil"
	"log"
	_ "net/http/pprof"

	"github.com/lornasong/diff-checker/src/compare"
	"github.com/lornasong/diff-checker/src/console"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.MemProfile).Stop()

	a, err := ioutil.ReadFile("a.txt")
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile("b.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := compare.MatchLine(string(a), string(b))
	console.NewPrinter(lines).Diff()
}
