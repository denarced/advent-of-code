package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2412"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2024-12.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	fmt.Println("Price:")
	fmt.Printf("    Without discount: %d.\n", aoc2412.DeriveTotalPrice(lines, false))
	fmt.Printf("    With discount:    %d.\n", aoc2412.DeriveTotalPrice(lines, true))

	shared.Logger.Info("Done.")
}
