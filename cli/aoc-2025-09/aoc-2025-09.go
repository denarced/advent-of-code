package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2509"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")
	defer shared.Logger.Info("Done.")
	defer shared.SetupCPUProfiling("profile-2025-09")()

	lines, err := shared.ReadLinesFromFile("data/2025-09.txt")
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Biggest rectangle:")
	fmt.Printf("    All:       %d\n", aoc2509.DeriveBiggestRectangle(lines, false))
	fmt.Printf("    Red/green: %d\n", aoc2509.DeriveBiggestRectangle(lines, true))
}
