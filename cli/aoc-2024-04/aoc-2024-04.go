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
	file, err := os.Open("data/2024-04.txt")
	shared.Die(err, "open file")
	defer file.Close()
	table, err := aoc2024.ReadLines(file)
	shared.Die(err, "ReadLines")
	fmt.Println("XMAS count:", aoc2024.CountInTable(table, "XMAS"))
	fmt.Println("MAX-MAX count:", aoc2024.CountWordCrosses(table, "MAS"))
	shared.Logger.Info("Done.")
}
