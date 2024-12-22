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

	fmt.Printf("Lowest score:\n")
	alpha := time.Now()
	fmt.Printf("    %d (%v)\n", aoc2416.CountLowestScore(lines, false), time.Since(alpha))

	shared.Logger.Info("Done.")
}
