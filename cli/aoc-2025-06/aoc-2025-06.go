package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2506"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := inr.ReadPath("data/2025-06.txt", inr.NoTrim())
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Sum:")
	fmt.Printf("    Standard math:   %d\n", aoc2506.Calculate(lines, true))
	fmt.Printf("    Cephalopod math: %d\n", aoc2506.Calculate(lines, false))
	shared.Logger.Info("Done.")
}
