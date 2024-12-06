// Package main implements the main CLI.
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

	shared.Logger.Info("Read data.")
	file, err := os.Open("data/2024-01.txt")
	shared.Die(err, "open file")
	defer file.Close()
	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	leftStrs, rightStrs := shared.ToColumns(lines)
	left, err := aoc2024.ToInts(leftStrs)
	shared.Die(err, "ToInts(leftStrs)")

	right, err := aoc2024.ToInts(rightStrs)
	shared.Die(err, "ToInts(rightStrs)")

	fmt.Println("Distance:", aoc2024.Advent01Distance(left, right))
	fmt.Println("Similarity:", aoc2024.Advent01Similarity(left, right))
	shared.Logger.Info("Done.")
}
