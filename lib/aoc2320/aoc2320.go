package aoc2320

import (
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

const (
	low pulse = iota
	high

	initName = "broadcaster"
)

func CountSignalProductFromLines(lines []string) int {
	shared.Logger.Info("Start counting signals.", "line count", len(lines))
	processors := parseLines(lines)
	shared.Logger.Info("Lines parsed.", "processor count", len(processors))
	raysetLength := len(processors)
	raysets := make([]rayset, raysetLength)
	var lowCount, highCount int
	for range 1_000 {
		// It's a bit silly but use a circular buffer to dramatically reduce allocations during
		// execution. Compared to simple deque (pop from beginning, append to end) benchmark
		// improved from 265µs to 144µs (-46%). Allocations dropped from 7k to 1k and memory usage
		// from 370KiB to 18KiB. That's with puzzle example so with real input the numbers would be
		// even more dramatic. Duration with real input: 10ms to 3ms. I.e. the optimization is
		// pointless but it was fun to try so why not.
		raysets[0] = rayset{name: "button", pul: low, targets: []string{initName}}
		rayCount := 1
		popIndex := 0
		pushIndex := popIndex + rayCount
		for rayCount > 0 {
			nextPopIndex := pushIndex
			var nextRayCount int
			for range rayCount {
				rs := raysets[popIndex%raysetLength]
				popIndex++
				for _, target := range rs.targets {
					if rs.pul == low {
						lowCount++
					} else {
						highCount++
					}
					if processor := processors[target]; processor != nil {
						if next := processor(rs.name, rs.pul); len(next.targets) > 0 {
							raysets[pushIndex%raysetLength] = next
							pushIndex++
							nextRayCount++
						}
					}
				}
			}
			popIndex = nextPopIndex
			rayCount = nextRayCount
			pushIndex = popIndex + rayCount
		}
	}
	total := lowCount * highCount
	shared.Logger.Info("Signals counted.", "total", total, "low", lowCount, "high", highCount)
	return total
}

type pulse int

type rayset struct {
	// Name of the sender.
	name string
	pul  pulse
	// Names of the targets where pulse is sent to.
	targets []string
}

type processor func(string, pulse) rayset

func parseLines(lines []string) map[string]processor {
	components := map[string]processor{}
	conjunctionRegistrars := map[string]func([]string){}
	componentsToCallers := map[string][]string{}
	addComponent := func(component, caller string) {
		callers := componentsToCallers[component]
		if callers == nil {
			callers = []string{}
		}
		callers = append(callers, caller)
		componentsToCallers[component] = callers
	}
	for _, each := range lines {
		pieces := gent.Map(strings.Split(each, "->"), strings.TrimSpace)
		targets := gent.Map(strings.Split(pieces[1], ","), strings.TrimSpace)
		if pieces[0] == initName {
			components[pieces[0]] = func(_ string, pul pulse) rayset {
				return rayset{
					name:    pieces[0],
					pul:     pul,
					targets: targets,
				}
			}
			for _, target := range targets {
				addComponent(target, pieces[0])
			}
			continue
		}
		prefix, name := pieces[0][0], pieces[0][1:]
		for _, target := range targets {
			addComponent(target, name)
		}
		switch prefix {
		case '%':
			components[name] = createFlipFlop(name, targets)
		case '&':
			var reg func([]string)
			components[name], reg = createConjunction(name, targets)
			conjunctionRegistrars[name] = reg
		default:
			panic("unknown type: " + string(prefix))
		}
	}
	for name, registrar := range conjunctionRegistrars {
		registrar(componentsToCallers[name])
	}
	return components
}

func createFlipFlop(name string, targets []string) processor {
	on := false
	return func(_ string, pul pulse) (ray rayset) {
		ray.name = name
		if pul == high {
			return
		}
		ray.targets = targets
		ray.pul = gent.Tri(on, low, high)
		on = !on
		return
	}
}

func createConjunction(name string, targets []string) (processor, func([]string)) {
	var shooters map[string]pulse
	return func(shooter string, pul pulse) (ray rayset) {
			shooters[shooter] = pul
			ray.name = name
			ray.targets = targets
			for _, each := range shooters {
				if each == low {
					ray.pul = high
					return
				}
			}
			return
		},
		func(shooterNames []string) {
			shooters = map[string]pulse{}
			for _, each := range shooterNames {
				shooters[each] = low
			}
		}
}
