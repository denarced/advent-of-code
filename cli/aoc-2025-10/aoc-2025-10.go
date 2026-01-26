package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2510"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")
	defer shared.Logger.Info("Done.")

	filep := "data/2025-10.txt"
	if len(os.Args) > 1 {
		filep = os.Args[1]
	}
	defer shared.SetupCPUProfiling("2025-10.profile")()
	lines, err := shared.ReadLinesFromFile(filep)
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Fewest clicks:")
	fmt.Printf(
		"    Indicator lights: %d\n",
		aoc2510.DeriveFewestClicks(lines, true))
	fmt.Printf(
		"    Joltage levels:   %d\n",
		aoc2510.DeriveFewestClicks(lines, false))
}
