package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2215"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-15"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	fmt.Printf("Count of no-go for beacons: %d\n", aoc2215.FindNoGoForBeacons(lines, 2_000_000))
	fmt.Printf("Beacon frequency:           %d\n", aoc2215.FindBeaconFrequency(lines, 4_000_000))
	shared.Logger.Info("Done.")
}
