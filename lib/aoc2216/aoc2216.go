package aoc2216

import (
	"fmt"
	"slices"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

const firstChamberName = "AA"

func DeriveMaximumPressure(lines []string) int {
	valves := parseLines(lines)
	valveNameToValve := mapValves(valves)
	nameToChamber := mapChambers(valves)
	primeTail := shared.AddLink(nil, firstChamberName)
	travelers := []*traveler{
		{
			chamber:   nameToChamber[firstChamberName],
			valveTail: primeTail,
			fullTail:  primeTail,
		},
	}
	shared.Logger.Info("Start travelling.")
	for clockSeconds := range 30 {
		var addedTravelers []*traveler
		for _, each := range travelers {
			currValve := valveNameToValve[each.chamber.name]
			if currValve.flowRate > 0 && !each.hasTurnedValve(each.chamber.name) {
				shared.Logger.Debug("New traveler to turn valve.", "name", currValve.name)
				addedTravelers = append(
					addedTravelers,
					each.copyWithRate(currValve.flowRate, 29-clockSeconds),
				)
			}
			var firstTaken bool
			for _, name := range currValve.connectedNames {
				if isSteppingOnTail(each.valveTail, name) {
					continue
				}
				if firstTaken {
					nextTraveler := each.copyWithChamber(nameToChamber[name])
					addedTravelers = append(addedTravelers, nextTraveler)
					continue
				}
				each.chamber = nameToChamber[name]
				each.valveTail = shared.AddLink(each.valveTail, name)
				each.fullTail = shared.AddLink(each.fullTail, name)
				firstTaken = true
			}
		}
		countBefore := len(travelers)
		travelers = append(travelers, addedTravelers...)
		countAfter := len(travelers)
		if countBefore < countAfter {
			shared.Logger.Debug(
				"Traveler count increased.",
				"before", countBefore,
				"after", countAfter)
		}
		var droppedTravelerIndexes []int
		for i, each := range travelers {
			if each.rate > 0 {
				if !each.chamber.updateSupreme(clockSeconds, each.rate) {
					droppedTravelerIndexes = append(droppedTravelerIndexes, i)
				}
			}
		}
		for i := len(droppedTravelerIndexes) - 1; i >= 0; i-- {
			index := droppedTravelerIndexes[i]
			shared.Logger.Debug("Drop traveler.", "rate", travelers[index].rate)
			travelers = append(travelers[:index], travelers[index+1:]...)
		}
	}
	shared.Logger.Info("Travelling done.", "traveller count", len(travelers))
	var maximus int
	for i, each := range travelers {
		if each.rate > 0 {
			shared.Logger.Info(
				"Traveller.",
				"i", i,
				"rate", each.rate,
				"valve tail", renderTail(each.valveTail),
				"full tail", renderTail(each.fullTail),
				"turned valves", each.rateValves)
		}
		maximus = max(maximus, each.rate)
	}
	shared.Logger.Info("Maximum pressure derived.", "maximum", maximus)
	return maximus
}

func parseLines(lines []string) []*valve {
	valves := make([]*valve, 0, len(lines))
	for _, each := range lines {
		trimmed := strings.TrimSpace(each)
		if trimmed != "" {
			valves = append(valves, parseLine(each))
		}
	}
	return valves
}

type valve struct {
	name           string
	flowRate       int
	connectedNames []string
}

func parseLine(line string) *valve {
	fields := strings.Fields(strings.ReplaceAll(line, ", ", " "))
	name := fields[1]
	var flowRate int
	parseCount, err := fmt.Sscanf(fields[4], "rate=%d", &flowRate)
	if err != nil || parseCount != 1 {
		shared.Logger.Error(
			"Failed to scan flow rate.",
			"line", line,
			"err", err,
			"count", parseCount)
		panic(err)
	}
	connectedNames := fields[9:]
	return &valve{
		name:           name,
		flowRate:       flowRate,
		connectedNames: connectedNames,
	}
}

func mapValves(valves []*valve) map[string]*valve {
	mapped := make(map[string]*valve, len(valves))
	for _, each := range valves {
		mapped[each.name] = each
	}
	return mapped
}

type chamber struct {
	name      string
	timeToMax map[int]int
}

func mapChambers(valves []*valve) map[string]*chamber {
	chambers := map[string]*chamber{}
	for _, each := range valves {
		aChamber := &chamber{
			name:      each.name,
			timeToMax: map[int]int{},
		}
		chambers[each.name] = aChamber
	}
	return chambers
}

type traveler struct {
	chamber    *chamber
	rate       int
	rateValves []string
	valveTail  *shared.Link[string]
	fullTail   *shared.Link[string]
}

func (v *traveler) copyWithChamber(aChamber *chamber) *traveler {
	return &traveler{
		chamber:    aChamber,
		rate:       v.rate,
		rateValves: v.rateValves,
		valveTail:  shared.AddLink(v.valveTail, aChamber.name),
		fullTail:   shared.AddLink(v.fullTail, aChamber.name),
	}
}

func (v *traveler) copyWithRate(rate, remainingSeconds int) *traveler {
	rateValves := append([]string(nil), v.rateValves...)
	rateValves = append(rateValves, v.chamber.name)
	return &traveler{
		chamber:    v.chamber,
		rate:       v.rate + remainingSeconds*rate,
		rateValves: rateValves,
		valveTail:  shared.AddLink(nil, v.chamber.name),
		fullTail:   shared.AddLink(v.fullTail, v.chamber.name),
	}
}

func (v *traveler) hasTurnedValve(name string) bool {
	return slices.Contains(v.rateValves, name)
}

// Update supreme rate. Return true if it was.
func (v *chamber) updateSupreme(seconds, rate int) bool {
	current := v.timeToMax[seconds]
	updated := current < rate
	if updated {
		shared.Logger.Debug("Update chamber max.", "name", v.name, "seconds", seconds, "rate", rate)
		v.timeToMax[seconds] = rate
	}
	return updated
}

func renderTail(root *shared.Link[string]) string {
	var joints []string
	for j := root; j != nil; j = j.Parent {
		joints = append(joints, j.Item)
	}
	slices.Reverse(joints)
	return strings.Join(joints, "-")
}

func isSteppingOnTail(tail *shared.Link[string], name string) bool {
	for l := tail; l != nil; l = l.Parent {
		if l.Item == name {
			return true
		}
	}
	return false
}
