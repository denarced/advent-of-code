package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2209"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-09"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Visited locations:")
	fmt.Printf("    Short tail(1): %d\n", aoc2209.CountVisitedPositions(lines, 1))
	fmt.Printf("    Long tail(9):  %d\n", aoc2209.CountVisitedPositions(lines, 9))
	shared.Logger.Info("Done.")
}
