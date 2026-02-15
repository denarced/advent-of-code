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

	fmt.Println("Path count:")
	for _, each := range []string{"you", "svr"} {
		fmt.Printf("    %s -> out: %d\n", each, aoc2511.CountPaths(lines, each))
	}
	shared.Logger.Info("Done.")
}
