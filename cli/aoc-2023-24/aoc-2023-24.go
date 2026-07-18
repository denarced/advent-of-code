package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2324"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-24"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	minimum := int64(200_000_000_000_000)
	maximum := int64(400_000_000_000_000)
	fmt.Printf("Intersections: %d\n", aoc2324.CountIntersections(lines, minimum, maximum))
	shared.Logger.Info("Done.")
}
