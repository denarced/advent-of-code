package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2309"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-09"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Sum of extrapolated values:")
	fmt.Printf("    Right: %d\n", aoc2309.SumExtrapolatedValues(lines, true))
	fmt.Printf("    Left : %d\n", aoc2309.SumExtrapolatedValues(lines, false))
	shared.Logger.Info("Done.")
}
