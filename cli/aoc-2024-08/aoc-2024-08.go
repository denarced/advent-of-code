package main

import (
	"fmt"
	"strings"

	"github.com/denarced/advent-of-code/lib/aoc2408"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := shared.ReadLinesFromFile("data/2024-08.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Unique antinode location count:")
	tab := strings.Repeat(" ", 3)
	fmt.Println(
		tab,
		"Without resonant harmonics:",
		aoc2408.CountUniqueAntinodeLocations(lines, false),
	)
	fmt.Println(
		tab,
		"With resonant harmonics:   ",
		aoc2408.CountUniqueAntinodeLocations(lines, true),
	)

	shared.Logger.Info("Done.")
}
