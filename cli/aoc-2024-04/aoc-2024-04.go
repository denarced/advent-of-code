package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2404"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := shared.ReadLinesFromFile("data/2024-04.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("XMAS count:    %d\n", aoc2404.CountInTable(lines, "XMAS"))
	fmt.Printf("MAX-MAX count: %d\n", aoc2404.CountWordCrosses(lines, "MAS"))

	shared.Logger.Info("Done.")
}
