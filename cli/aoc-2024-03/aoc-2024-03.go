package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/denarced/advent-of-code/lib/aoc2024"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")
	file, err := os.Open("data/2024-03.txt")
	shared.Die(err, "open file")
	defer file.Close()
	lines, err := aoc2024.ReadLines(file)
	shared.Die(err, "ReadLines")
	text := strings.Join(lines, "")
	fmt.Println("Without do/don't:", aoc2024.Multiply(text, false))
	fmt.Println("With do/don't:", aoc2024.Multiply(text, true))
	shared.Logger.Info("Done.")
}
