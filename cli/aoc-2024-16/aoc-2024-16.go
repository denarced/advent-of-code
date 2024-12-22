package main

import (
	"fmt"
	"os"
	"time"

	"github.com/denarced/advent-of-code/lib/aoc2416"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2024-16.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	alpha := time.Now()
	score, seatCount := aoc2416.CountLowestScore(lines, false)
	fmt.Printf("Result (%v):\n", time.Since(alpha))
	fmt.Printf("    Score:      %d\n", score)
	fmt.Printf("    Seat count: %d\n", seatCount)

	shared.Logger.Info("Done.")
}
