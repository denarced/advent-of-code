package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2414"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines := gent.OrPanic2(shared.ReadLinesFromFile("data/2024-14.txt"))("ReadLinesFromFile")

	fmt.Println("Safety factor:", aoc2414.DeriveSafetyFactor(lines, 101, 103, 100))
	fmt.Println("Steps to find Christmas tree:", aoc2414.FindChristmasTree(lines, 101, 103))

	shared.Logger.Info("Done.")
}
