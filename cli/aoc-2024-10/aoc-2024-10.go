package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2410"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := shared.ReadLinesFromFile("data/2024-10.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Sum of trailhead scores: ", aoc2410.DeriveSumOfTrailheadScores(lines, false))
	fmt.Println("Sum of trailhead ratings:", aoc2410.DeriveSumOfTrailheadScores(lines, true))

	shared.Logger.Info("Done.")
}
