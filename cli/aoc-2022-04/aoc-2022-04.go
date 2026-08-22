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

	fmt.Println("Overlap count:")
	fmt.Printf("    Complete: %d\n", aoc2204.CountOverlaps(lines, true))
	fmt.Printf("    Partial:  %d\n", aoc2204.CountOverlaps(lines, false))
	shared.Logger.Info("Done.")
}
