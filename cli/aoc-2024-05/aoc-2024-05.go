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

	fmt.Println("Correct middle page number sum:", aoc2405.SumCorrectMiddlePageNumbers(lines))
	fmt.Println("Incorrect middle page number sum:", aoc2405.SumIncorrectMiddlePageNumbers(lines))

	shared.Logger.Info("Done.")
}
