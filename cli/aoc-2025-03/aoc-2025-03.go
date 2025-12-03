package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2503"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2025-03.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	fmt.Println("Maximum joltage sum:")
	fmt.Printf("     2 batteries: %d\n", aoc2503.DeriveMaxJoltageSum(lines, 2))
	fmt.Printf("    12 batteries: %d\n", aoc2503.DeriveMaxJoltageSum(lines, 12))
	shared.Logger.Info("Done.")
}
