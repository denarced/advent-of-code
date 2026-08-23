package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2205"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-05"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := inr.ReadPath(fmt.Sprintf("data/%s.txt", id), inr.IncludeEmpty(), inr.NoTrim())
	shared.Die(err, "ReadLinesFromFile")

	fmt.Println("Top crates:")
	fmt.Printf("    One-by-one: %s\n", aoc2205.DeriveTopCrates(lines, true))
	fmt.Printf("    At once:    %s\n", aoc2205.DeriveTopCrates(lines, false))
	shared.Logger.Info("Done.")
}
