package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2406"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file := shared.OrPanic2(os.Open("data/2024-06.txt"))("open input file")
	defer file.Close()
	lines := shared.OrPanic2(shared.ReadLines(file))("read lines")

	fmt.Println("Distinct positions:", aoc2406.CountDistinctPositions(lines))
	fmt.Println(
		"Count of blocks resulting in indefinite loops:",
		aoc2406.CountBlocksForIndefiniteLoops(lines).Count(),
	)

	shared.Logger.Info("Done.")
}
