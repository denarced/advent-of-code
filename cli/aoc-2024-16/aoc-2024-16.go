package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2416"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file := gent.OrPanic2(os.Open("data/2024-16.txt"))("open file")
	defer file.Close()
	lines := gent.OrPanic2(shared.ReadLines(file))("ReadLines")

	score, seatCount := aoc2416.CountLowestScore(lines, false)
	fmt.Println("Result:")
	fmt.Printf("    Score:      %d\n", score)
	fmt.Printf("    Seat count: %d\n", seatCount)

	shared.Logger.Info("Done.")
}
