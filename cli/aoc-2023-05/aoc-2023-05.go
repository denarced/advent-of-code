package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2305"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-05"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := inr.ReadPath(fmt.Sprintf("data/%s.txt", id), inr.IncludeEmpty())
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Lowest location: %d\n", aoc2305.DeriveLowestLocation(lines))
	shared.Logger.Info("Done.")
}
