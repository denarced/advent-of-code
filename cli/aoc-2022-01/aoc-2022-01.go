package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2201"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-01"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := inr.ReadPath(fmt.Sprintf("data/%s.txt", id), inr.IncludeEmpty())
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Heaviest: %d\n", aoc2201.DeriveHeaviest(lines))
	shared.Logger.Info("Done.")
}
