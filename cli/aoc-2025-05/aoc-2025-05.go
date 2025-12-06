package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2505"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines, err := inr.ReadPath("data/2025-05.txt", inr.IncludeEmpty())
	shared.Die(err, "read puzzle input")

	fmt.Println("How many ingredients are fresh?")
	fmt.Printf("    Within available: %d\n", aoc2505.CountFreshAvailableIngredients(lines))
	fmt.Printf("    In general:       %d\n", aoc2505.CountFreshIngredients(lines))
	shared.Logger.Info("Done.")
}
