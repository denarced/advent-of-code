package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2323"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-23"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Longest path:")
	fmt.Printf("    Downhill: %d\n", aoc2323.FindLongestPath(lines))
	fmt.Printf("    +Uphill:  %d\n", aoc2323.FindLongestPathWithGraph(lines))
	shared.Logger.Info("Done.")
}
