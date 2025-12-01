package aoc2501

import (
	"strconv"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func SolvePassword(lines []string, allZeroes bool) int {
	pointer := 50
	count := 0
	rotations := splitValues(lines)
	for _, rot := range rotations {
		shared.Logger.Debug(
			"Rotate.",
			"rotation",
			rot,
			"count",
			count,
			"pointer",
			pointer,
		)
		for _, value := range rot.values {
			previous := pointer
			pointer += rot.m * value
			if allZeroes && value != 0 {
				flippedOver := previous > 0 && pointer <= 0
				flippedUnder := previous <= 99 && pointer > 99
				flowedOver := previous == 0 && (pointer == -100 || pointer == 100)
				if flippedOver || flippedUnder || flowedOver {
					count++
				}
			}
			pointer = (pointer + 100) % 100
		}
		if !allZeroes && pointer == 0 {
			count++
		}
	}
	return count
}

type rotation struct {
	m      int
	values []int
}

// SplitValues split values into pieces of 100 (max).
// {"L201"} -> [{m: -1, values: [100, 100, 1]}]
// {"R100"} -> [{m: 1, values: [100]}]
func splitValues(lines []string) []rotation {
	var rotations []rotation
	for _, each := range lines {
		var rot rotation
		rot.m = 1
		if each[0] == 'L' {
			rot.m = -1
		}
		totalValue := gent.OrPanic2(strconv.Atoi(each[1:]))("convert rotation value")
		if totalValue == 0 {
			rot.values = []int{0}
			rotations = append(rotations, rot)
			continue
		}
		for totalValue != 0 {
			if totalValue < 0 {
				panic("illegal value")
			}
			value := min(totalValue, 100)
			totalValue -= value
			rot.values = append(rot.values, value)
		}
		rotations = append(rotations, rot)
	}
	return rotations
}
