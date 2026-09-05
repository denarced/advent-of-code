package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2214"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-14"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Sand grains:")
	fmt.Printf("    Until feeding abyss: %d\n", aoc2214.MeasureSand(lines, false))
	fmt.Printf("    With infinite floor: %d\n", aoc2214.MeasureSand(lines, true))
	shared.Logger.Info("Done.")
}
