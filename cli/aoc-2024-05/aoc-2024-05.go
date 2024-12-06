package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2024"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2024-05.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := aoc2024.ReadLines(file)
	shared.Die(err, "ReadLines")

	fmt.Println("Correct middle page number sum:", aoc2024.SumCorrectMiddlePageNumbers(lines))
	fmt.Println("Incorrect middle page number sum:", aoc2024.SumIncorrectMiddlePageNumbers(lines))

	shared.Logger.Info("Done.")
}
