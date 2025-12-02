package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2502"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2025-02.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	fmt.Println("Invalid ID sums:")
	fmt.Printf("    Twice: %d\n", aoc2502.SumInvalidIDs(lines[0], true))
	fmt.Printf("    More:  %d\n", aoc2502.SumInvalidIDs(lines[0], false))
	shared.Logger.Info("Done.")
}
