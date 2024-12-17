package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2413"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2024-13.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	fmt.Println("Fewest tokens to have it all:")
	fmt.Printf("    Without unit conversion fix: %d.\n", aoc2413.DeriveFewestTokens(lines, false))
	fmt.Printf("    With unit conversion fix:    %d.\n", aoc2413.DeriveFewestTokens(lines, true))
	shared.Logger.Info("Done.")
}
