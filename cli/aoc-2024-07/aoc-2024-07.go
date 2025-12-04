package main

import (
	"fmt"
	"strings"

	"github.com/denarced/advent-of-code/lib/aoc2407"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines := gent.OrPanic2(shared.ReadLinesFromFile("data/2024-07.txt"))("ReadLinesFromFile")

	fmt.Println("Calibration sums:")
	tab := strings.Repeat(" ", 3)
	fmt.Println(tab, "Without concat:", aoc2407.DeriveCalibrationSum(lines, false))
	fmt.Println(tab, "With concat:   ", aoc2407.DeriveCalibrationSum(lines, true))

	shared.Logger.Info("Done.")
}
