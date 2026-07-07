package main

import (
	"fmt"
	"os"

	"github.com/denarced/advent-of-code/lib/aoc2322"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-22"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	os.Setenv("OMG", "777")
	defer os.Unsetenv("OMG")

	count := aoc2322.CountBricksFromLines(lines)
	suffix := ""
	if count == 895 || count == 578 {
		suffix = " (wrong answer, too high)"
	}
	fmt.Printf("Brick count: %d%s\n", count, suffix)
	shared.Logger.Info("Done.")
}
