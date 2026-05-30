package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2320"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-20"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Total signal count: %d\n", aoc2320.CountSignalProductFromLines(lines))
	shared.Logger.Info("Done.")
}
