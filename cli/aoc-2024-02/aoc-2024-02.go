package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2402"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")
	file, err := os.Open("data/2024-02.txt")
	shared.Die(err, "open file")
	defer file.Close()
	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")
	table := shared.ToIntTable(lines)
	fmt.Printf("Safe count without dampener: %d\n", aoc2402.CountSafe(table, false))
	fmt.Printf("Safe count with dampener:    %d\n", aoc2402.CountSafe(table, true))
	shared.Logger.Info("Done.")
}
