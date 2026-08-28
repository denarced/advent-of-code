package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2208"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-08"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Visible tree count: %d\n", aoc2208.DeriveVisibleTreeCount(lines))
	shared.Logger.Info("Done.")
}
