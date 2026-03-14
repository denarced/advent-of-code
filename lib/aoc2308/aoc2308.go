package aoc2308

import (
	"fmt"
	"math"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

func getNext(nod *node, r rune) *node {
	switch r {
	case 'L':
		return nod.left
	case 'R':
		return nod.right
	default:
		shared.Logger.Error("Invalid step on path.", "step", string(r))
		panic("invalid step on path, should be L or R")
	}
}

func CountSteps(lines []string) int {
	path, nodes := parseLines(lines)
	shared.Logger.Info("Count steps.", "path length", len(path), "node count", len(nodes))
	current := nodes["AAA"]
	steps := []rune(path)
	var i int
	for {
		each := steps[i%len(steps)]
		i++
		next := getNext(current, each)
		if next.name == "ZZZ" {
			shared.Logger.Info(
				"ZZZ found.",
				"current", current.name,
				"left", current.left.name,
				"right", current.right.name)
			break
		}
		shared.Logger.Info("Step forward.", "from", current.name, "to", next.name)
		current = next
	}
	shared.Logger.Info("Steps counted.", "count", i)
	return i
}

type pathSpec struct {
	firstCount  int
	repeatCount int
}

func CountStepsInSync(lines []string) int {
	path, nodes := parseLines(lines)
	var starters []string
	for key := range nodes {
		if key[len(key)-1] == 'A' {
			starters = append(starters, key)
		}
	}
	shared.Logger.Info(
		"Count steps in sync.",
		"path length", len(path),
		"node count", len(nodes),
		"starter count", len(starters),
	)

	pathSpecs := make([]pathSpec, len(starters))
	for i, each := range starters {
		pathSpecs[i] = findPathSpecs(path, nodes[each])
	}
	shared.Logger.Info("Path specs derived.", "specs", pathSpecs)
	findLowest := func() int {
		low := math.MaxInt
		var index int
		for i, each := range pathSpecs {
			if each.firstCount < low {
				index = i
				low = each.firstCount
			}
		}
		return index
	}
	areAllEqual := func() bool {
		value := -1
		for _, each := range pathSpecs {
			if value < 0 {
				value = each.firstCount
			} else if value != each.firstCount {
				return false
			}
		}
		return true
	}
	for {
		index := findLowest()
		pathSpecs[index].firstCount += pathSpecs[index].repeatCount
		shared.Logger.Debug("Path specs.", "specs", pathSpecs)
		if areAllEqual() {
			count := pathSpecs[0].firstCount
			shared.Logger.Info("Steps counted.", "count", count)
			return count
		}
	}
}

type node struct {
	left, right *node
	name        string
}

func parseLines(lines []string) (string, map[string]*node) {
	shared.Logger.Debug("Parse lines.", "line count", len(lines))
	blocks := shared.SplitToBlocks(lines)
	if len(blocks[0]) != 1 {
		shared.Logger.Error("First block should have one line for RL steps.", "block", blocks[0])
		panic("invalid first block")
	}
	if len(blocks) != 2 {
		shared.Logger.Error("Should have exactly two blocks.", "length", len(blocks))
		panic("invalid block count")
	}
	nodes := map[string]*node{}
	for _, line := range blocks[1] {
		pieces := strings.SplitN(line, "=", 2)
		assertLength(
			len(pieces),
			2,
			"line",
			"line", line)
		name := strings.TrimSpace(pieces[0])
		assertLength(len(name), 3, "name", "name", name)
		nod, ok := nodes[name]
		if !ok {
			nod = &node{name: name}
			nodes[name] = nod
		}

		leftName, rightName := parseTargetNodes(strings.TrimSpace(pieces[1]))
		leftNode, ok := nodes[leftName]
		if !ok {
			leftNode = &node{name: leftName}
			nodes[leftName] = leftNode
		}
		rightNode, ok := nodes[rightName]
		if !ok {
			rightNode = &node{name: rightName}
			nodes[rightName] = rightNode
		}

		nod.left = leftNode
		nod.right = rightNode
	}

	return blocks[0][0], nodes
}

func assertLength(actual, expected int, name string, logArgs ...any) {
	if actual == expected {
		return
	}
	logArgs = append(
		[]any{
			"name", name,
			"expected length", expected,
			"actual length", actual,
		},
		logArgs...)
	shared.Logger.Error("Invalid values found, panicking.", logArgs...)
	panic(fmt.Sprintf("invalid %s found", name))
}

func parseTargetNodes(s string) (left, right string) {
	after, found := strings.CutPrefix(s, "(")
	if !found {
		shared.Logger.Error("Target nodes string doesn't have parenthesis prefix.", "s", s)
		panic("required parenthesis prefix not found")
	}
	main, found := strings.CutSuffix(after, ")")
	if !found {
		shared.Logger.Error("Target nodes string doesn't have parenthesis suffix.", "s", s)
		panic("required parenthesis suffix not found")
	}
	pieces := strings.FieldsFunc(main, func(r rune) bool { return r == ',' || r == ' ' })
	assertLength(
		len(pieces),
		2,
		"target names",
		"names", s)
	left = pieces[0]
	right = pieces[1]
	return
}

func findPathSpecs(path string, nod *node) (spec pathSpec) {
	steps := []rune(path)
	var i int
	for {
		each := steps[i%len(steps)]
		i++
		nod = getNext(nod, each)
		if nod.name[len(nod.name)-1] == 'Z' {
			spec.firstCount = i
			break
		}
	}
	count := 0
	for {
		each := steps[i%len(steps)]
		i++
		count++
		nod = getNext(nod, each)
		if nod.name[len(nod.name)-1] == 'Z' {
			spec.repeatCount = count
			break
		}
	}
	return
}
