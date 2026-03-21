package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2311"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-11"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Sum of distances:")
	fmt.Printf("    Young galaxies: %d\n", aoc2311.SumDistances(lines, 2))
	fmt.Printf("    Old galaxies:   %d\n", aoc2311.SumDistances(lines, 1_000_000))
	shared.Logger.Info("Done.")
}
