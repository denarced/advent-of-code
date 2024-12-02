// Package aoc2024 contains implementation for 2024 Advent of Code solutions.
package aoc2024

import (
	"bytes"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
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
	for i := 0; i < len(left); i++ {
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
