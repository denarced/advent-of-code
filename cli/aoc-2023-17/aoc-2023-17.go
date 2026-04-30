package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2317"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-17"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Least heat loss:")
	fmt.Printf("    Normal: %d\n", aoc2317.DeriveLeastHeatLoss(lines, 1, 3))
	fmt.Printf("    Ultra: %d\n", aoc2317.DeriveLeastHeatLoss(lines, 4, 10))
	shared.Logger.Info("Done.")
}
