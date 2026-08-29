package main

import (
	"fmt"
	"strings"

	"github.com/denarced/advent-of-code/lib/aoc2210"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2022-10"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	sum, crtLines := aoc2210.SumSignalStrengths(lines)
	fmt.Printf("Sum of signal strengths: %d\n", sum)
	fmt.Println(strings.Join(crtLines, "\n"))
	shared.Logger.Info("Done.")
}
