package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2417"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	file, err := os.Open("data/2024-17.txt")
	shared.Die(err, "open file")
	defer file.Close()

	lines, err := shared.ReadLines(file)
	shared.Die(err, "ReadLines")

	fmt.Println("Output:", aoc2417.DeriveOutput(lines))

	shared.Logger.Info("Done.")
}
