package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2206"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-06"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Characters until package: %d\n", aoc2206.DeriveStartOfPackage(lines[0]))
	shared.Logger.Info("Done.")
}
