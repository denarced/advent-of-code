package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2512"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := inr.ReadPath("data/2025-12.txt", inr.IncludeEmpty())
	shared.Die(err, "ReadLinesFromFile")

	count, err := aoc2512.CountFittingRegions(lines)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(2)
	}
	fmt.Printf("Region count: %d\n", count)
	shared.Logger.Info("Done.")
}
