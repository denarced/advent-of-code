package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/denarced/advent-of-code/lib/aoc2408"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2024-08.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

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
