package aoc2320

import (
	"sort"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

const (
	Low Pulse = iota
	High

	initName = "broadcaster"
)

type Pulse int

func (v Pulse) Flip() Pulse {
	if v == Low {
		return High
	}
	return Low
}

type FiringSquad struct {
	SignalCb         func(aPulse Pulse)
	RoundCb          func(int, map[string]*Processor) bool
	processors       map[string]*Processor
	ComponentCallers map[string][]string
}

func NewFiringSquad(lines []string) *FiringSquad {
	shared.Logger.Info("Create FiringSquad.", "line count", len(lines))
	processors, componentsToCallers := parseLines(lines)
	shared.Logger.Info("Lines parsed.", "processor count", len(processors))
	return &FiringSquad{processors: processors, ComponentCallers: componentsToCallers}
}

func (v *FiringSquad) Fire() {
	shared.Logger.Info("Fire.", "monitor processors", v.RoundCb != nil)
	raysetLength := len(v.processors)
	raysets := make([]rayset, raysetLength)
	var i int
	for i = 0; ; i++ {
		// It's a bit silly but use a circular buffer to dramatically reduce allocations during
		// execution. Compared to simple deque (pop from beginning, append to end) benchmark
		// improved from 265µs to 144µs (-46%). Allocations dropped from 7k to 1k and memory usage
		// from 370KiB to 18KiB. That's with puzzle example so with real input the numbers would be
		// even more dramatic. Duration with real input: 10ms to 3ms. I.e. the optimization is
		// pointless but it was fun to try so why not.
		raysets[0] = rayset{name: "button", pulse: Low, targets: []string{initName}}
		rayCount := 1
		popIndex := 0
		pushIndex := popIndex + rayCount
		var halt bool
		for rayCount > 0 {
			nextPopIndex := pushIndex
			var nextRayCount int
			for range rayCount {
				rs := raysets[popIndex%raysetLength]
				popIndex++
				for _, target := range rs.targets {
					if v.SignalCb != nil {
						v.SignalCb(rs.pulse)
					}
					if processor := v.processors[target]; processor != nil {
						if next := processor.process(rs.name, rs.pulse); len(next.targets) > 0 {
							raysets[pushIndex%raysetLength] = next
							pushIndex++
							nextRayCount++
						}
					}
				}
			}
			if v.RoundCb != nil && !v.RoundCb(i+1, v.processors) {
				shared.Logger.Info("RoundCb asked to halt.")
				halt = true
				break
			}
			popIndex = nextPopIndex
			rayCount = nextRayCount
			pushIndex = popIndex + rayCount
		}
		if v.RoundCb == nil && i >= 999 {
			break
		}

		if halt {
			break
		}
	}
	shared.Logger.Info("Fired done.", "button presses", i+1)
}

func (v Pulse) String() string {
	if v == Low {
		return "low"
	}
	return "high"
}

type rayset struct {
	// Name of the sender.
	name  string
	pulse Pulse
	// Names of the targets where pulse is sent to.
	targets []string
}

type Processor struct {
	All     func(Pulse) bool
	process func(string, Pulse) rayset
}

func parseLines(lines []string) (map[string]*Processor, map[string][]string) {
	components := map[string]*Processor{}
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
			components[pieces[0]] = &Processor{
				process: func(_ string, pul Pulse) rayset {
					return rayset{
						name:    pieces[0],
						pulse:   pul,
						targets: targets,
					}
				},
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
	return components, componentsToCallers
}

func createFlipFlop(name string, targets []string) *Processor {
	on := false
	return &Processor{
		process: func(_ string, pul Pulse) (ray rayset) {
			ray.name = name
			if pul == High {
				return
			}
			ray.targets = targets
			ray.pulse = gent.Tri(on, Low, High)
			on = !on
			return
		},
		All: func(_ Pulse) bool {
			panic("no one should ever call flip-flop All")
		},
	}
}

func createConjunction(name string, targets []string) (*Processor, func([]string)) {
	var shooters map[string]Pulse
	var keys []string
	proc := Processor{
		process: func(shooter string, pul Pulse) (ray rayset) {
			shooters[shooter] = pul
			ray.name = name
			ray.targets = targets
			for _, each := range shooters {
				if each == Low {
					ray.pulse = High
					return
				}
			}
			return
		},
		All: func(pulse Pulse) bool {
			for _, each := range shooters {
				if each != pulse {
					return false
				}
			}
			return true
		},
	}
	return &proc, func(shooterNames []string) {
		shooters = map[string]Pulse{}
		keys = []string{}
		for _, each := range shooterNames {
			shooters[each] = Low
			keys = append(keys, each)
		}
		sort.Strings(keys)
	}
}

func FindTracked(components map[string][]string, name string, pulse Pulse) ([]string, Pulse) {
	for {
		callers := components[name]
		if len(callers) == 0 {
			return nil, Low
		}
		pulse = pulse.Flip()
		if len(callers) == 1 {
			name = callers[0]
			continue
		}
		return callers, pulse
	}
}

type SignalTracker struct {
	LowCount  int
	HighCount int
}

func (v *SignalTracker) Add(pulse Pulse) {
	if pulse == Low {
		v.LowCount++
	} else {
		v.HighCount++
	}
}

type RoundMonitor struct {
	components         []string
	expectedPulse      Pulse
	counts             map[string][]int
	resolvedComponents *gent.Set[string]
	Frequencies        []int
}

func NewRoundMonitor(components []string, expectedPulse Pulse) *RoundMonitor {
	return &RoundMonitor{
		components:         components,
		expectedPulse:      expectedPulse,
		counts:             make(map[string][]int, len(components)),
		resolvedComponents: gent.NewSet[string](),
		Frequencies:        make([]int, len(components)),
	}
}

func (v *RoundMonitor) Monitor(clickCount int, processors map[string]*Processor) bool {
	// Component could be initialized to the expected state but that's ignored because such an
	// obvious accident clearly isn't the sought answer.
	if clickCount <= 1 {
		return true
	}
	for i, each := range v.components {
		if v.resolvedComponents.Has(each) {
			continue
		}
		if !processors[each].All(v.expectedPulse) {
			continue
		}
		last, ok := getLast(v.counts[each])
		if ok && last == clickCount {
			continue
		}
		v.counts[each] = append(v.counts[each], clickCount)
		if len(v.counts[each]) > 10 {
			freq, ok := resolveFrequency(v.counts[each])
			if ok {
				v.Frequencies[i] = freq
				v.resolvedComponents.Add(each)
				shared.Logger.Info(
					"Component frequency found.",
					"name", each,
					"frequency", freq)
			}
		}
	}
	if v.resolvedComponents.Count() >= len(v.components) {
		shared.Logger.Info("All component frequencies found, quitting.")
		return false
	}
	return true
}

func getLast(values []int) (int, bool) {
	if len(values) == 0 {
		return 0, false
	}
	return values[len(values)-1], true
}

func resolveFrequency(values []int) (int, bool) {
	size := len(values)
	if size < 10 {
		return 0, false
	}
	last := -1
	for i := size - 1; i >= size-5; i-- {
		prev, next := values[i-1], values[i]
		freq := next - prev
		if last < 0 {
			last = freq
			continue
		}
		if freq != last {
			return 0, false
		}
	}
	return last, true
}
