package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2209"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-09"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Visited positions: %d\n", aoc2209.CountVisitedPositions(lines))
	shared.Logger.Info("Done.")
}
