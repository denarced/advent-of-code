package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2204"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-04"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Assignment pairs: %d\n", aoc2204.CountAssignmentPairs(lines))
	shared.Logger.Info("Done.")
}
