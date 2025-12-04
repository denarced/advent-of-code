package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2406"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines := gent.OrPanic2(shared.ReadLinesFromFile("data/2024-06.txt"))("ReadLinesFromFile")

	fmt.Println("Distinct positions:", aoc2406.CountDistinctPositions(lines))
	fmt.Println(
		"Count of blocks resulting in indefinite loops:",
		aoc2406.CountBlocksForIndefiniteLoops(lines).Count(),
	)

	shared.Logger.Info("Done.")
}
