package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2402"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")
	lines, err := shared.ReadLinesFromFile("data/2024-02.txt")
	shared.Die(err, "ReadLinesFromFile")
	table := shared.ToIntTable(lines)
	fmt.Printf("Safe count without dampener: %d\n", aoc2402.CountSafe(table, false))
	fmt.Printf("Safe count with dampener:    %d\n", aoc2402.CountSafe(table, true))
	shared.Logger.Info("Done.")
}
