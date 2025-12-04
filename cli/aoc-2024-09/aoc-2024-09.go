package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2409"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := shared.ReadLinesFromFile("data/2024-09.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Checksum:       ", aoc2409.CountChecksum(lines[0]))
	fmt.Println("Defrag checksum:", aoc2409.CountDefragmentedChecksum(lines[0]))

	shared.Logger.Info("Done.")
}
