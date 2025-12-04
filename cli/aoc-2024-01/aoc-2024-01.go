// Package main implements the main CLI.
package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2401"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	shared.Logger.Info("Read data.")
	lines, err := shared.ReadLinesFromFile("data/2024-01.txt")
	shared.Die(err, "ReadLinesFromFile")

	leftStrs, rightStrs := shared.ToColumns(lines)
	left, err := shared.ToInts(leftStrs)
	shared.Die(err, "ToInts(leftStrs)")

	right, err := shared.ToInts(rightStrs)
	shared.Die(err, "ToInts(rightStrs)")

	fmt.Printf("Distance:   %d\n", aoc2401.Distance(left, right))
	fmt.Printf("Similarity: %d\n", aoc2401.Similarity(left, right))
	shared.Logger.Info("Done.")
}
