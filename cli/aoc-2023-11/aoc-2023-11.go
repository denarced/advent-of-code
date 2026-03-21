package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2311"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-11"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Sum of distances: %d\n", aoc2311.SumDistances(lines))
	shared.Logger.Info("Done.")
}
