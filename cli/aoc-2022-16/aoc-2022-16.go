package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2216"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-16"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Maximum achievable pressure:")
	fmt.Printf("    Alone:         %d\n", aoc2216.DeriveMaximumPressure(lines, 30, true))
	fmt.Printf("    With elephant: %d\n", aoc2216.DeriveMaximumPressure(lines, 26, false))
	shared.Logger.Info("Done.")
}
