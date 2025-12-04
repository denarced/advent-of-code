package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2504"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := shared.ReadLinesFromFile("data/2025-04.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Movable rolls:")
	fmt.Printf("    One try: %d\n", aoc2504.CountRolls(lines, 1))
	fmt.Printf("    ~ tries: %d\n", aoc2504.CountRolls(lines, -1))
	shared.Logger.Info("Done.")
}
