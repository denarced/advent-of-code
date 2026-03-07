package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2307"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-07"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Winnings:")
	fmt.Printf("    Total: %d\n", aoc2307.CountWinnings(lines))
	shared.Logger.Info("Done.")
}
