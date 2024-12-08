package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2024"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2024-06.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	fmt.Println("Distinct positions:", aoc2024.CountDistinctPositions(lines))
	fmt.Println(
		"Count of blocks resulting in indefinite loops:",
		aoc2024.CountBlocksForIndefiniteLoops(lines).Count(),
	)

	shared.Logger.Info("Done.")
}
