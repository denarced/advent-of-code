package main

import (
	"fmt"

	"github.com/denarced/advent-of-code/lib/aoc2320"
	"github.com/denarced/advent-of-code/shared"
)

func main() {
	shared.InitLogging()
	shared.Logger.Info("Start.")

	id := "2023-20"
	//revive:disable-next-line:defer
	defer shared.SetupCPUProfiling(fmt.Sprintf("%s.profile", id))()
	lines, err := shared.ReadLinesFromFile(fmt.Sprintf("data/%s.txt", id))
	shared.Die(err, "ReadLinesFromFile")

	squad := aoc2320.NewFiringSquad(lines)
	tracker := new(aoc2320.SignalTracker)
	squad.SignalCb = tracker.Add
	squad.Fire()
	fmt.Println("Total count:")
	fmt.Printf("    Signals:                %d\n", tracker.LowCount*tracker.HighCount)

	squad = aoc2320.NewFiringSquad(lines)
	trackedComponents, expectedPulse := aoc2320.FindTracked(
		squad.ComponentCallers,
		"rx",
		aoc2320.Low,
	)
	if len(trackedComponents) < 2 {
		panic("expected more components to track")
	}
	monitor := aoc2320.NewRoundMonitor(trackedComponents, expectedPulse)
	squad.RoundCb = monitor.Monitor
	squad.Fire()
	fmt.Printf(
		"    For rx to receive low:  %d\n",
		shared.DeriveLeastCommonMultiple(
			monitor.Frequencies[0],
			monitor.Frequencies[1],
			monitor.Frequencies[2:]...),
	)
	shared.Logger.Info("Done.")
}
