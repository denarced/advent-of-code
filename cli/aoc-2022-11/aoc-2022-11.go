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

	fmt.Printf("Level of monkey business: %d\n", aoc2211.DeriveMonkeyBusiness(lines))
	shared.Logger.Info("Done.")
}
