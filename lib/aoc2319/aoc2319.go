package aoc2319

import (
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

func SumRatings(lines []string) int {
	shared.Logger.Info("Sum ratings.", "line count", len(lines))
	nameToFlow, parts := parseLines(lines)
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Got situation.", "worksflows", nameToFlow, "parts", parts)
	}
	steps := make([]step, len(parts))
	for i, each := range parts {
		steps[i] = step{workflowName: "in", aPart: each}
	}
	var sum int
	var aStep step
	var count int
	for len(steps) > 0 {
		count++
		aStep, steps = steps[0], steps[1:]
		nextStep := processStep(aStep, nameToFlow[aStep.workflowName])
		switch nextStep.workflowName {
		case "R":
			continue
		case "A":
			sub := nextStep.aPart.sum()
			if shared.IsDebugEnabled() {
				shared.Logger.Debug("Add to the sum.", "sub", sub)
			}
			sum += sub
		default:
			steps = append(steps, nextStep)
		}
	}
	shared.Logger.Info("Got ratings sum.", "sum", sum, "loop count", count)
	return sum
}

func parseLines(lines []string) (map[string]workflow, []part) {
	blocks := shared.SplitToBlocks(lines)
	if len(blocks) != 2 {
		shared.Logger.Error("Invalid block count.", "count", len(blocks))
		panic("invalid data - bad block count")
	}
	workflows := parseWorkflows(blocks[0])
	parts := parseParts(blocks[1])
	return workflows, parts
}

type step struct {
	workflowName string
	aPart        part
}

type workflow struct {
	name  string
	specs []spec
}

type spec struct {
	attr        string
	dest        string
	value       int
	less        bool
	endComplete bool
}

type part map[string]int

func (v part) sum() int {
	var sum int
	for _, v := range v {
		sum += v
	}
	return sum
}

func (v part) pick(name string) int {
	val, ok := v[name]
	if !ok {
		panic("no such part attribute: " + name)
	}
	return val
}

func parseWorkflows(lines []string) map[string]workflow {
	m := map[string]workflow{}
	for _, each := range lines {
		flow := parseWorkflow(each)
		m[flow.name] = flow
	}
	return m
}

func parseWorkflow(s string) workflow {
	openIndex := strings.Index(s, "{")
	closeIndex := strings.Index(s, "}")
	name := s[:openIndex]
	flow := workflow{name: name}
	for _, each := range strings.Split(s[openIndex+1:closeIndex], ",") {
		flow.specs = append(flow.specs, parseWorkflowSpec(each))
	}
	return flow
}

func parseWorkflowSpec(s string) spec {
	pieces := strings.Split(s, ":")
	if len(pieces) == 1 {
		return spec{dest: s, endComplete: true}
	}
	opPieces := strings.Split(pieces[0], "<")
	less := true
	if len(opPieces) == 1 {
		opPieces = strings.Split(pieces[0], ">")
		less = false
	}
	val := parseIntOrDie(opPieces[1])
	return spec{
		attr:  opPieces[0],
		dest:  pieces[1],
		value: val,
		less:  less,
	}
}

func parseParts(lines []string) []part {
	var parts []part
	for _, each := range lines {
		parts = append(parts, parsePart(each))
	}
	return parts
}

func parsePart(s string) part {
	openIndex := strings.Index(s, "{")
	closeIndex := strings.Index(s, "}")
	content := s[openIndex+1 : closeIndex]
	aPart := map[string]int{}
	for _, each := range strings.Split(content, ",") {
		pieces := strings.Split(each, "=")
		num := pieces[1]
		aPart[pieces[0]] = parseIntOrDie(num)
	}

	return aPart
}

func parseIntOrDie(num string) int {
	val, err := strconv.Atoi(num)
	if err != nil {
		shared.Logger.Error("Failed to parse int.", "s", num, "err", err)
		panic(err)
	}
	return val
}

func processStep(aStep step, flow workflow) step {
	for _, each := range flow.specs {
		if each.endComplete {
			return step{workflowName: each.dest, aPart: aStep.aPart}
		}
		value := aStep.aPart.pick(each.attr)
		if each.less {
			if value < each.value {
				return step{workflowName: each.dest, aPart: aStep.aPart}
			}
		} else if value > each.value {
			return step{workflowName: each.dest, aPart: aStep.aPart}
		}
	}
	panic("my ass")
}
