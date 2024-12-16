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

	file, err := os.Open("data/2024-11.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	stones, err := shared.ToInts(strings.Fields(lines[0]))
	shared.Die(err, "ToInts")
	fmt.Println("Stone count:")
	for _, each := range []int{25, 75} {
		alpha := time.Now()
		count := aoc2411.CountStones(each, stones)
		elapsed := time.Since(alpha)
		fmt.Printf("    %02d blinks: %d (%v)\n", each, count, elapsed)
	}

	shared.Logger.Info("Done.")
}
