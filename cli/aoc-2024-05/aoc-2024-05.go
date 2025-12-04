package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2405"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := shared.ReadLinesFromFile("data/2024-05.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Middle page number sum:")
	fmt.Printf("    Correct:   %d\n", aoc2405.SumCorrectMiddlePageNumbers(lines))
	fmt.Printf("    Incorrect: %d\n", aoc2405.SumIncorrectMiddlePageNumbers(lines))

	shared.Logger.Info("Done.")
}
