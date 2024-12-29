package main

import (
	"fmt"
	"os"
	"strings"
	"time"

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
	for _, each := range []int{25, 75} {
		alpha := time.Now()
		count := aoc2411.CountStones(stones, each)
		elapsed := time.Since(alpha)
		fmt.Printf("    %02d blinks: %d (%v)\n", each, count, elapsed)
	}

	shared.Logger.Info("Done.")
}
