package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2413"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := shared.ReadLinesFromFile("data/2024-13.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Fewest tokens to have it all:")
	fmt.Printf("    Without unit conversion fix: %d.\n", aoc2413.DeriveFewestTokens(lines, false))
	fmt.Printf("    With unit conversion fix:    %d.\n", aoc2413.DeriveFewestTokens(lines, true))
	shared.Logger.Info("Done.")
}
