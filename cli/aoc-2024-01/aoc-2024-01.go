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
	die(err, "open file")
	defer file.Close()
	lines, err := aoc2024.ReadLines(file)
	die(err, "ReadLines")

	leftStrs, rightStrs := aoc2024.ToColumns(lines)
	left, err := aoc2024.ToInts(leftStrs)
	die(err, "ToInts(leftStrs)")

	right, err := aoc2024.ToInts(rightStrs)
	die(err, "ToInts(rightStrs)")

	fmt.Println("Distance:", aoc2024.Advent01Distance(left, right))
	fmt.Println("Similarity:", aoc2024.Advent01Similarity(left, right))
	shared.Logger.Info("Done.")
}

func die(err error, message string) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Message: \"%s\". Error: %s.\n", message, err)
	//revive:disable-next-line:deep-exit
	os.Exit(2)
}
