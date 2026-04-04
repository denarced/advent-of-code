package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2312"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-12"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Sum of arrangement permutations: %d\n", aoc2312.SumPermutations(lines))
	shared.Logger.Info("Done.")
}
