package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2511"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := shared.ReadLinesFromFile("data/2025-11.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Path count:", aoc2511.CountPaths(lines))
	shared.Logger.Info("Done.")
}
