package aoc2308

import (
	"fmt"
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
