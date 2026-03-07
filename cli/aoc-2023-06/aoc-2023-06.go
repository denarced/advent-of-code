package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2306"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-06"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Ways to win:")
	fmt.Printf("    Multiple races: %d\n", aoc2306.MultiplyCounts(lines, true))
	fmt.Printf("    One race:       %d\n", aoc2306.MultiplyCounts(lines, false))
	shared.Logger.Info("Done.")
}
