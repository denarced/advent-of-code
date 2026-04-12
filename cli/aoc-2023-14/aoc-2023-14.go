package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2314"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-14"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Total load: %d\n", aoc2314.CountTotalLoad(lines))
	shared.Logger.Info("Done.")
}
