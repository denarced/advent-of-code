package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2322"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-22"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	count := aoc2322.CountBricksFromLines(lines)
	fmt.Printf("Brick count:     %d\n", count)
	fmt.Printf("Total fallout: %d\n", aoc2322.KillBricks(lines))
	shared.Logger.Info("Done.")
}
