package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2217"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-17"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Tower height: %d\n", aoc2217.MeasureTowerHeight(lines))
	shared.Logger.Info("Done.")
}
