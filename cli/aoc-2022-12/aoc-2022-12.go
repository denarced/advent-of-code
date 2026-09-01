package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2212"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-12"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Fewest steps:")
	fmt.Printf("    Start to end: %d\n", aoc2212.DeriveFewestSteps(lines, aoc2212.StartToEnd))
	fmt.Printf("    Low to end:   %d\n", aoc2212.DeriveFewestSteps(lines, aoc2212.LowToEnd))
	shared.Logger.Info("Done.")
}
