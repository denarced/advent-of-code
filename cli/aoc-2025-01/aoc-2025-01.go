package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2501"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := shared.ReadLinesFromFile("data/2025-01.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Password:")
	fmt.Printf("    Final zeroes: %d\n", aoc2501.SolvePassword(lines, false))
	fmt.Printf("    All zeroes:   %d\n", aoc2501.SolvePassword(lines, true))
	shared.Logger.Info("Done.")
}
