package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2308"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-08"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := inr.ReadPath(fmt.Sprintf("data/%s.txt", id), inr.IncludeEmpty())
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Step count: %d\n", aoc2308.CountSteps(lines))
	shared.Logger.Info("Done.")
}
