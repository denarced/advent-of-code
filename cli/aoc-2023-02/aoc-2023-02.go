package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2302"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-02"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	gameCountSum := aoc2302.DeriveGameCountSum(
		lines,
		map[aoc2302.Kind]int{
			aoc2302.KindRed:   12,
			aoc2302.KindGreen: 13,
			aoc2302.KindBlue:  14,
		},
	)
	fmt.Println("Sums:")
	fmt.Printf("    Sum of feasible IDs: %d\n", gameCountSum)
	fmt.Printf("    Sum of powers      : %d\n", aoc2302.DerivePowerSum(lines))
	shared.Logger.Info("Done.")
}
