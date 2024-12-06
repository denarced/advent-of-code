// Package aoc2024 contains implementation for 2024 Advent of Code solutions.
package aoc2024

import (
	"bytes"
	"io"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

var (
	mulPattern = regexp.MustCompile(`mul\(\d+,\d+\)|don't\(\)|do\(\)`)
	directions = []direction{{1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}}
)

func Advent01Distance(left, right []int) int {
	shared.Logger.Info(
		"Advent 01: derive distance.",
		"left length",
		len(left),
		"right length",
		len(right),
	)
	slices.Sort(left)
	slices.Sort(right)
	distance := 0
	for i := range len(left) {
		distance += shared.Abs(left[i] - right[i])
		if distance < 0 {
			panic("Distance is <0 (int overflow).")
		}
	}
	shared.Logger.Info("Distance derived.", "distance", distance)
	return distance
}

func Advent01Similarity(left, right []int) int {
	shared.Logger.Info(
		"Advent 01: derive similarity",
		"left length",
		len(left),
		"right length",
		len(right),
	)
	counts := deriveCounts(right)
	similarity := 0
	for _, each := range left {
		c, ok := counts[each]
		if ok {
			similarity += (c * each)
		}
		if similarity < 0 {
			panic("Similarity is <0 (int overflow).")
		}
	}
	shared.Logger.Info("Similarity derived.", "similarity", similarity)
	return similarity
}

func deriveCounts(s []int) map[int]int {
	counts := map[int]int{}
	for _, each := range s {
		curr, ok := counts[each]
		if ok {
			counts[each] = curr + 1
		} else {
			counts[each] = 1
		}
	}
	return counts
}

func ReadAll(reader io.Reader) (s string, err error) {
	shared.Logger.Info("Read all.")

	var b []byte
	b, err = io.ReadAll(reader)
	if err != nil {
		return
	}
	s = string(b)
	return
}

func ReadLines(reader io.Reader) (lines []string, err error) {
	shared.Logger.Info("Read lines.")

	var b []byte
	b, err = io.ReadAll(reader)
	if err != nil {
		return
	}

	for _, each := range bytes.Split(b, []byte("\n")) {
		line := strings.TrimSpace(string(each))
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return
}

func ToColumns(s []string) (left []string, right []string) {
	shared.Logger.Info("Split slice content to two coluns.", "length", len(s))
	for _, each := range s {
		pieces := trim(strings.Split(each, " "))
		left = append(left, pieces[0])
		right = append(right, pieces[1])
	}
	return
}

func ToInts(s []string) (nums []int, err error) {
	shared.Logger.Info("Convert string slice to ints.", "length", len(s))
	for _, each := range s {
		var n int
		n, err = strconv.Atoi(each)
		if err != nil {
			shared.Logger.Error("Failed to convert to int.", "string", each, "err", err)
			return
		}
		nums = append(nums, n)
	}
	return
}

func trim(s []string) (trimmed []string) {
	for _, each := range s {
		aTrimmed := strings.TrimSpace(each)
		if aTrimmed == "" {
			continue
		}
		trimmed = append(trimmed, aTrimmed)
	}
	return
}

func ToIntTable(s []string) (table [][]int) {
	for _, each := range s {
		cells := strings.Fields(each)
		var row []int
		for _, c := range cells {
			n, err := strconv.Atoi(c)
			if err != nil {
				panic("Invalid number: " + c)
			}
			row = append(row, n)
		}
		if row == nil {
			panic("Empty row")
		}
		table = append(table, row)
	}
	return
}

func CountSafe(levels [][]int, dampener bool) int {
	count := 0
	for _, each := range levels {
		if len(each) < 2 {
			continue
		}
		index := deriveUnsafe(each)
		if index < 0 {
			count++
			continue
		}
		if !dampener {
			continue
		}

		for i := index + 1; i >= 0; i-- {
			trimmed := append([]int{}, each[0:i]...)
			trimmed = append(trimmed, each[i+1:]...)
			index = deriveUnsafe(trimmed)
			if index < 0 {
				shared.Logger.Debug("Safe after trimmed.", "original", each, "trimmed", trimmed)
				break
			}
		}
		if index < 0 {
			count++
		}
	}
	return count
}

// Derive index where unsafe was detected or -1 if levels are safe.
func deriveUnsafe(levels []int) int {
	asc := levels[0] < levels[1]
	for i := range len(levels) - 1 {
		first, second := levels[i], levels[i+1]
		if asc && first > second || !asc && first < second {
			return i
		}
		diff := shared.Abs(first - second)
		if diff < 1 || diff > 3 {
			return i
		}
	}
	return -1
}

func Multiply(text string, logic bool) int {
	total := 0
	skipping := false
	for {
		pair := mulPattern.FindStringIndex(text)
		if pair == nil {
			break
		}
		piece := text[pair[0]:pair[1]]
		text = text[pair[1]:]

		if strings.HasPrefix(piece, "don't") {
			skipping = true
			continue
		}
		if strings.HasPrefix(piece, "do") {
			skipping = false
			continue
		}
		if logic && skipping {
			continue
		}
		if strings.HasPrefix(piece, "mul") {
			a, b := splitMul(piece)
			total += a * b
		}
	}
	return total
}

//revive:disable-next-line:confusing-results
func splitMul(s string) (int, int) {
	separated := s[4 : len(s)-1]
	broken := strings.Split(separated, ",")
	return toInt(broken[0]), toInt(broken[1])
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic("Failed to convert to int: " + s)
	}
	return i
}

func CountInTable(table []string, word string) int {
	if len(word) == 0 {
		return 0
	}
	count := 0
	for r := 0; r < len(table); r++ {
		for c := 0; c < len(table[r]); c++ {
			if table[r][c] == word[0] {
				count += countWordsAt(table, word, r, c)
			}
		}
	}
	return count
}

func CountWordCrosses(table []string, word string) int {
	locations := findWordLocations(table, word)
	counts := map[loc]int{}
	total := 0
	for _, each := range locations {
		if count, ok := counts[each]; ok {
			if count > 1 {
				panic("wtf")
			}
			counts[each]++
			total++
			shared.Logger.Info("Found word cross.", "location", each)
		} else {
			shared.Logger.Debug("Found half of word cross.", "location", each)
			counts[each] = 1
		}
	}
	return total
}

type direction struct {
	x int
	y int
}

func countWordsAt(table []string, word string, row, col int) int {
	directions := []direction{{1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}}
	count := 0
	for _, each := range directions {
		if readTableAt(table, row, col, len(word), each) == word {
			count++
		}
	}
	return count
}

type loc struct {
	x int
	y int
}

func findWordLocations(table []string, word string) []loc {
	if len(word)%2 != 1 {
		panic("only works with odd length words: 3, 5, 7, ...")
	}
	var locations []loc
	mid := len(word) / 2
	for r := 0; r < len(table); r++ {
		for c := 0; c < len(table[r]); c++ {
			for _, each := range directions {
				if each.x == 0 || each.y == 0 {
					continue
				}
				if readTableAt(table, r, c, len(word), each) == word {
					x := r + each.x*mid
					y := c + each.y*mid
					locations = append(locations, loc{x, y})
				}
			}
		}
	}
	return locations
}

func readTableAt(table []string, row, col, count int, dir direction) string {
	result := ""
	for range count {
		if row < 0 || row >= len(table) {
			break
		}
		line := table[row]
		if col < 0 || col >= len(line) {
			break
		}
		result += line[col : col+1]
		row += dir.x
		col += dir.y
	}
	return result
}
