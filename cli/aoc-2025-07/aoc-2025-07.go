package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2507"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := shared.ReadLinesFromFile("data/2025-07.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Count:")
	fmt.Printf("    Splits: %d\n", aoc2507.CountSplits(lines))
	shared.Logger.Info("Done.")
}
