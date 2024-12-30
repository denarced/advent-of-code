package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/denarced/advent-of-code/lib/aoc2411"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file := shared.OrPanic2(os.Open("data/2024-11.txt"))("open file")
	defer file.Close()
	lines := shared.OrPanic2(shared.ReadLines(file))("ReadLines")
	stones := shared.OrPanic2(shared.ToInts(strings.Fields(lines[0])))("ToInts")
	fmt.Println("Stone count:")
	for i := 25; i < 100; i += 50 {
		count := aoc2411.CountStones(stones, i)
		fmt.Printf("    %02d blinks: %d\n", i, count)
	}

	shared.Logger.Info("Done.")
}
