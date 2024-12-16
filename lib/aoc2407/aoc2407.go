package aoc2407

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

func DeriveCalibrationSum(lines []string, withConcat bool) int {
	dtos := toCalibrationDtos(lines)
	total := 0
	sumTotal := 0
	for _, each := range dtos {
		sumTotal += each.sum
		if isValid(each, withConcat) {
			total += each.sum
		}
	}
	shared.Logger.Info("Sums.", "total", sumTotal, "valid", total, "with concat", withConcat)
	return total
}

type calibrationDto struct {
	sum   int
	parts []int
}

func toCalibrationDtos(lines []string) []calibrationDto {
	trimmed := shared.FilterValues(
		shared.MapValues(
			lines,
			func(s string) string {
				return strings.TrimSpace(s)
			}),
		func(s string) bool {
			return s != ""
		})
	return shared.MapValues(
		trimmed,
		func(s string) calibrationDto {
			sides := shared.MapValues(
				strings.Split(s, ":"),
				func(s string) string {
					return strings.TrimSpace(s)
				})
			parts, err := shared.ToInts(strings.Fields(sides[1]))
			shared.Die(err, "right side to ints")
			sum, err := strconv.Atoi(sides[0])
			shared.Die(err, "left side to int")
			return calibrationDto{sum: sum, parts: parts}
		})
}

func isValid(dto calibrationDto, withConcat bool) bool {
	if len(dto.parts) == 1 {
		return dto.sum == dto.parts[0]
	}
	base := 2
	if withConcat {
		base = 3
	}
	for _, each := range generatePermutations(len(dto.parts)-1, base) {
		shared.Logger.Debug("Try permutation.", "dto", dto, "permutation", each)
		sum := deriveSum(dto.parts, each)
		if sum == dto.sum {
			shared.Logger.Info("Valid calibration.", "dto", dto, "operators", each)
			return true
		}
		shared.Logger.Debug("Invalid calibration.", "dto", dto, "operators", each)
	}
	shared.Logger.Debug("Invalid calibration.", "dto", dto)
	return false
}

func deriveSum(nums, operators []int) int {
	res, tail := nums[0], nums[1:]
	for i := 0; i < len(tail); i++ {
		switch operators[i] {
		case 0:
			res += tail[i]
		case 1:
			res *= tail[i]
		case 2:
			res = concat(res, tail[i])
		default:
			panic(fmt.Sprintf("Unknown operator: %d.", operators[i]))
		}
	}
	shared.Logger.Debug("Derived.", "result", res, "nums", nums, "operators", operators)
	return res
}

func generatePermutations(length, base int) [][]int {
	strs := generateZeroPaddedNumericStrings(length, base)
	return shared.MapValues(
		strs,
		func(s string) []int {
			nums, err := shared.ToInts(strings.Split(s, ""))
			if err != nil {
				panic(err)
			}
			return nums
		})
}

func concat(a, b int) int {
	i, err := strconv.Atoi(fmt.Sprintf("%d%d", a, b))
	if err != nil {
		panic(err)
	}
	return i
}

func generateZeroPaddedNumericStrings(length, base int) []string {
	max := shared.Pow(base, length)
	strs := make([]string, 0, max)
	for i := 0; i < max; i++ {
		s := strconv.FormatInt(int64(i), base)
		missing := length - len(s)
		if missing > 0 {
			s = strings.Repeat("0", missing) + s
		}
		strs = append(strs, s)
	}
	return strs
}
