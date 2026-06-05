package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2321"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-21"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Range: %d\n", aoc2321.CountRangeFromLines(lines, 64))
	shared.Logger.Info("Done.")
}
