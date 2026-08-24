package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2207"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-07"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

    fmt.Println("Size:")
    total, toDelete := aoc2207.SumRecursiveDirSize(lines)
	fmt.Printf("    Total:     %d\n", total)
	fmt.Printf("    To delete: %d\n", toDelete)
	shared.Logger.Info("Done.")
}
