package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2417"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := shared.ReadLinesFromFile("data/2024-17.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Output:", aoc2417.DeriveOutput(lines))

	shared.Logger.Info("Done.")
}
