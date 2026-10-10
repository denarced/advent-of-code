package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2218"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-18"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Total surface area: %d\n", aoc2218.DeriveSurfaceArea(lines))
	shared.Logger.Info("Done.")
}
