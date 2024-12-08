package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/denarced/advent-of-code/lib/aoc2407"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2024-07.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	fmt.Println("Calibration sums:")
	tab := strings.Repeat(" ", 3)
	fmt.Println(tab, "Without concat:", aoc2407.DeriveCalibrationSum(lines, false))
	fmt.Println(tab, "With concat:   ", aoc2407.DeriveCalibrationSum(lines, true))

	shared.Logger.Info("Done.")
}
