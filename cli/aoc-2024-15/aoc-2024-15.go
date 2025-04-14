package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2415"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file := gent.OrPanic2(os.Open("data/2024-15.txt"))("open input file")
	defer file.Close()
	lines := gent.OrPanic2(shared.ReadLines(file))("ReadLines")

	fmt.Printf("Sum of GPS coordinates:\n")
	fmt.Printf("    Single: %d\n", aoc2415.CountCoordinateSum(lines, false))
	fmt.Printf("    Double: %d\n", aoc2415.CountCoordinateSum(lines, true))

	shared.Logger.Info("Done.")
}
