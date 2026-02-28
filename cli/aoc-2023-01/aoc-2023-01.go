package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2301"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := inr.ReadPath("data/2023-01.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Sum:")
	fmt.Printf("    Digits   : %d\n", aoc2301.SumCalibrationValues(lines, true))
	fmt.Printf("    And words: %d\n", aoc2301.SumCalibrationValues(lines, false))
}
