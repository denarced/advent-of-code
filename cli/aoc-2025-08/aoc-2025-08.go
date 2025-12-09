package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2508"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	defer shared.SetupCPUProfiling("profile2508")()
	lines, err := shared.ReadLinesFromFile("data/2025-08.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Circuit math:")
	fmt.Printf("    : %d\n", aoc2508.CountCircuits(lines, 1000))
	shared.Logger.Info("Done.")
}
