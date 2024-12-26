package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2404"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2024-04.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	fmt.Printf("XMAS count:    %d\n", aoc2404.CountInTable(lines, "XMAS"))
	fmt.Printf("MAX-MAX count: %d\n", aoc2404.CountWordCrosses(lines, "MAS"))

	shared.Logger.Info("Done.")
}
