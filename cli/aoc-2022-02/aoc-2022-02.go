package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2202"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-02"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Total score:")
	fmt.Printf("    what to play: %d\n", aoc2202.DeriveTotalScore(lines, false))
	fmt.Printf("    how to end:   %d\n", aoc2202.DeriveTotalScore(lines, true))
	shared.Logger.Info("Done.")
}
