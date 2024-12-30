package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2414"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file := shared.OrPanic2(os.Open("data/2024-14.txt"))("open file")
	defer file.Close()
	lines := shared.OrPanic2(shared.ReadLines(file))("ReadLines")

	fmt.Println("Safety factor:", aoc2414.DeriveSafetyFactor(lines, 101, 103, 100))
	fmt.Println("Steps to find Christmas tree:", aoc2414.FindChristmasTree(lines, 101, 103))

	shared.Logger.Info("Done.")
}
