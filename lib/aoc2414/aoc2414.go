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
	minimumNeighbourCount := 20
	for i := 1; i < 1_000_000; i++ {
		yToCoords := make(map[int][]int, height)
		skip := true
		for j, each := range ints {
			x, y := deriveCoordinates(each, width, height, i)
			xs := yToCoords[y]
			if xs == nil {
				xs = make([]int, 0, 20)
			}
			xs = append(xs, x)
			yToCoords[y] = xs
			if j >= minimumNeighbourCount && minimumNeighbourCount <= len(xs) {
				skip = false
			}
		}

		if skip {
			continue
		}

		count := countNeighbours(yToCoords, minimumNeighbourCount)
		if count >= minimumNeighbourCount {
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
	a := specs[0] + steps*specs[2]
	b := specs[1] + steps*specs[3]
	x = (a%width + width) % width
	y = (b%height + height) % height
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
	result := make([][]int, 0, len(lines))
	for _, each := range lines {
		// E.g. "p=98,97 v=25,80" -> "p=98,97" and "v=25,80".
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

func countNeighbours(m map[int][]int, minimum int) int {
	highest := 0
	for _, xSet := range m {
		if len(xSet) < minimum {
			continue
		}
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
		if count >= minimum {
			return count
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
