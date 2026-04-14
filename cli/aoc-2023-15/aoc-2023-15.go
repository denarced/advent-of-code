package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2315"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-15"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Hash sum:       %d\n", aoc2315.SumHashes(lines))
	fmt.Printf("Focusing power: %d\n", aoc2315.DeriveFocusingPower(lines))
	shared.Logger.Info("Done.")
}
