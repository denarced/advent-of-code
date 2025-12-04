package main

import (
	"fmt"
	"strings"

	"github.com/denarced/advent-of-code/lib/aoc2411"
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	lines := gent.OrPanic2(shared.ReadLinesFromFile("data/2024-11.txt"))("ReadLinesFromFile")
	stones := gent.OrPanic2(shared.ToInts(strings.Fields(lines[0])))("ToInts")
	fmt.Println("Stone count:")
	for i := 25; i < 100; i += 50 {
		count := aoc2411.CountStones(stones, i)
		fmt.Printf("    %02d blinks: %d\n", i, count)
	}

	shared.Logger.Info("Done.")
}
