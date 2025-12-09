package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2507"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	defer shared.SetupCPUProfiling("profile2507")()
	lines, err := shared.ReadLinesFromFile("data/2025-07.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Count:")
	fmt.Printf("    Splits:    %d\n", aoc2507.CountSplits(lines))
	fmt.Printf("    Timelines: %d\n", aoc2507.CountTimelines(lines))
	shared.Logger.Info("Done.")
}
