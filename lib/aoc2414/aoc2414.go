package aoc2414

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

func DeriveSafetyFactor(lines []string, width, height, steps int) int {
	ints := parseLines(lines)
	quadrants := []int{0, 0, 0, 0, -999_999_999_999}
	for _, each := range ints {
		x, y := deriveCoordinates(each, width, height, steps)
		quadrants[deriveQuadrant(x, y, width, height)]++
	}
	shared.Logger.Info("Quadrants derived.", "quadrants", quadrants)
	return multiply(quadrants[:4])
}

func FindChristmasTree(lines []string, width, height int) int {
	ints := parseLines(lines)
	for i := 1; i < 1_000_000; i++ {
		coords := [][]int{}
		for _, each := range ints {
			x, y := deriveCoordinates(each, width, height, i)
			coords = append(coords, []int{x, y})
		}

		count := countNeighbours(coords)
		if count >= 20 {
			// printBoard(coords, width, height)
			return i
		}
	}
	return -1
}

func multiply(values []int) int {
	result := values[0]
	for _, each := range values[1:] {
		result *= each
	}
	return result
}

func deriveCoordinates(specs []int, width, height, steps int) (x int, y int) {
	local := append([]int{}, specs...)
	local[0] += steps * local[2]
	local[1] += steps * local[3]
	x = (local[0]%width + width) % width
	y = (local[1]%height + height) % height
	shared.Logger.Debug(
		"Outcome.",
		"local", local,
		"width", width,
		"height", height,
		"steps", steps,
		"x", x,
		"y", y)
	return
}

func deriveQuadrant(x, y, width, height int) int {
	midX := width / 2
	midY := height / 2
	if y < midY {
		if x < midX {
			return 0
		} else if x > midX {
			return 1
		}
	} else if y > midY {
		if x < midX {
			return 2
		} else if x > midX {
			return 3
		}
	}
	return 4
}

func parseLines(lines []string) [][]int {
	result := [][]int{}
	for _, each := range lines {
		p, v := splitPair(each, " ")
		if p == "" {
			continue
		}
		posX, posY := toIntPair(splitPair(strings.TrimPrefix(p, "p="), ","))
		dx, dy := toIntPair(splitPair(strings.TrimPrefix(v, "v="), ","))
		quad := []int{posX, posY, dx, dy}
		result = append(result, quad)
	}
	return result
}

func toIntPair(a, b string) (first int, second int) {
	var err error
	first, err = strconv.Atoi(a)
	shared.Die(err, "toIntPair, first")
	second, err = strconv.Atoi(b)
	shared.Die(err, "toIntPair, second")
	return
}

func splitPair(s, sep string) (first string, second string) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return "", ""
	}
	pieces := strings.Split(trimmed, sep)
	if len(pieces) != 2 {
		panic(
			fmt.Sprintf(
				"Unexpected result with split. Length: %d. Value: %s.",
				len(pieces),
				trimmed,
			),
		)
	}
	first = strings.TrimSpace(pieces[0])
	second = strings.TrimSpace(pieces[1])
	return
}

func countNeighbours(line [][]int) int {
	m := map[int][]int{}
	for _, each := range line {
		if xSet, exists := m[each[1]]; exists {
			xSet = append(xSet, each[0])
			m[each[1]] = xSet
		} else {
			m[each[1]] = []int{each[0]}
		}
	}

	highest := 0
	for _, xSet := range m {
		sort.Ints(xSet)
		previous := -1
		count := 0
		for _, each := range xSet {
			if previous == -1 {
				previous = each
				continue
			}
			if previous+1 == each {
				count++
			}
			previous = each
		}
		highest = shared.Max(highest, count)
	}
	return highest
}

func printBoard(coords [][]int, width, height int) {
	board := []string{}
	for range height {
		board = append(board, strings.Repeat(" ", width))
	}
	for _, each := range coords {
		line := board[each[1]]
		x := each[0]
		line = line[:x] + "." + line[x+1:]
		board[each[1]] = line
	}
	fmt.Println(strings.Join(board, "\n") + "\n")
}
