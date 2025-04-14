package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/denarced/advent-of-code/lib/aoc2407"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file := gent.OrPanic2(os.Open("data/2024-07.txt"))("open input file")
	defer file.Close()
	lines := gent.OrPanic2(shared.ReadLines(file))("ReadLines")

	fmt.Println("Calibration sums:")
	tab := strings.Repeat(" ", 3)
	fmt.Println(tab, "Without concat:", aoc2407.DeriveCalibrationSum(lines, false))
	fmt.Println(tab, "With concat:   ", aoc2407.DeriveCalibrationSum(lines, true))

	shared.Logger.Info("Done.")
}
