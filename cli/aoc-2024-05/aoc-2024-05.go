package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2405"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2024-05.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	fmt.Println("Middle page number sum:")
	fmt.Printf("    Correct:   %d\n", aoc2405.SumCorrectMiddlePageNumbers(lines))
	fmt.Printf("    Incorrect: %d\n", aoc2405.SumIncorrectMiddlePageNumbers(lines))

	shared.Logger.Info("Done.")
}
