package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2310"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-10"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Counts:")
	fmt.Printf("    Steps:           %d\n", aoc2310.CountSteps(lines))
	fmt.Printf("    Squeezed blocks: %d\n", aoc2310.FindCrackCount(lines))
	shared.Logger.Info("Done.")
}
