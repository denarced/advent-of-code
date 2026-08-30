package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2211"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-11"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Level of monkey business:")
	fmt.Printf(
		"    20 rounds without worries: %d\n",
		aoc2211.DeriveMonkeyBusiness(lines, 20, false),
	)
	fmt.Printf(
		"    10k rounds with worries:   %d\n",
		aoc2211.DeriveMonkeyBusiness(lines, 10_000, true),
	)
	shared.Logger.Info("Done.")
}
