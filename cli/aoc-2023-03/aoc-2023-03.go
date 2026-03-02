package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2303"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-03"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Sum of")
	fmt.Printf("    Part numbers: %d\n", aoc2303.SumPartNumbers(lines))
	fmt.Printf("    Gear ratios:  %d\n", aoc2303.SumGearRatios(lines))
	shared.Logger.Info("Done.")
}
