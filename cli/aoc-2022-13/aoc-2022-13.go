package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2213"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-13"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := inr.ReadPath(fmt.Sprintf("data/%s.txt", id), inr.IncludeEmpty())
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Result:")
	fmt.Printf("    Sum of sorted indices:        %d\n", aoc2213.SumSortedIndexes(lines))
	fmt.Printf("    Product of divider positions: %d\n", aoc2213.SortAll(lines))
	shared.Logger.Info("Done.")
}
