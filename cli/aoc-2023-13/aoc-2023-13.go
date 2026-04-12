package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2313"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-13"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := inr.ReadPath(fmt.Sprintf("data/%s.txt", id), inr.IncludeEmpty())
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Sum of reflections:")
	fmt.Printf("    With smudges:    %d\n", aoc2313.SumReflections(lines, false))
	fmt.Printf("    Without smudges: %d\n", aoc2313.SumReflections(lines, true))
	shared.Logger.Info("Done.")
}
