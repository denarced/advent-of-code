package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2403"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")
	file, err := os.Open("data/2024-03.txt")
	shared.Die(err, "open file")
	defer file.Close()
	text, err := shared.ReadAll(file)
	shared.Die(err, "ReadLines")
	fmt.Printf("Without do/don't: %d\n", aoc2403.Multiply(text, false))
	fmt.Printf("With do/don't:    %d\n", aoc2403.Multiply(text, true))
	shared.Logger.Info("Done.")
}
