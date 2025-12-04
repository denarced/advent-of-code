package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2412"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := shared.ReadLinesFromFile("data/2024-12.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Price:")
	fmt.Printf("    Without discount: %d.\n", aoc2412.DeriveTotalPrice(lines, false))
	fmt.Printf("    With discount:    %d.\n", aoc2412.DeriveTotalPrice(lines, true))

	shared.Logger.Info("Done.")
}
