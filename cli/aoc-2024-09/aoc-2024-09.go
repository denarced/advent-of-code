package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2409"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2024-09.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	fmt.Println("Checksum:       ", aoc2409.CountChecksum(lines[0]))
	fmt.Println("Defrag checksum:", aoc2409.CountDefragmentedChecksum(lines[0]))

	shared.Logger.Info("Done.")
}
