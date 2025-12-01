package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2501"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2025-01.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	fmt.Println("Password:")
	fmt.Printf("    Final zeroes: %d\n", aoc2501.SolvePassword(lines, false))
	fmt.Printf("    All zeroes:   %d\n", aoc2501.SolvePassword(lines, true))
	shared.Logger.Info("Done.")
}
