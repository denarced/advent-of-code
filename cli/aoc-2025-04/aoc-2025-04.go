package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2504"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2025-04.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	fmt.Println("Movable rolls:")
	fmt.Printf("    One try: %d\n", aoc2504.CountRolls(lines, 1))
	fmt.Printf("    ~ tries: %d\n", aoc2504.CountRolls(lines, -1))
	shared.Logger.Info("Done.")
}
