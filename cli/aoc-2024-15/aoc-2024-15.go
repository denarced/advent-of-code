package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2415"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2024-15.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	fmt.Printf("Sum of GPS coordinates:\n")
	fmt.Printf("    Single: %d\n", aoc2415.CountCoordinateSum(lines, false))
	fmt.Printf("    Double: %d\n", aoc2415.CountCoordinateSum(lines, true))

	shared.Logger.Info("Done.")
}
