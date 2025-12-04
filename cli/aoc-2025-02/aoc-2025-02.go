package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2502"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := shared.ReadLinesFromFile("data/2025-02.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Invalid ID sums:")
	fmt.Printf("    Twice: %d\n", aoc2502.SumInvalidIDs(lines[0], true))
	fmt.Printf("    More:  %d\n", aoc2502.SumInvalidIDs(lines[0], false))
	shared.Logger.Info("Done.")
}
