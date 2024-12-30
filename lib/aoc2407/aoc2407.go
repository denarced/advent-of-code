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
	permGen := newPermutationGenerator()
	for _, each := range dtos {
		sumTotal += each.sum
		if isValid(each, withConcat, permGen) {
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

func isValid(dto calibrationDto, withConcat bool, permGen *permutationGenerator) bool {
	if len(dto.parts) == 1 {
		return dto.sum == dto.parts[0]
	}
	base := 2
	if withConcat {
		base = 3
	}
	for _, each := range permGen.generate(permutationSpec{length: len(dto.parts) - 1, base: base}) {
		sum := deriveSum(dto, each)
		if sum == dto.sum {
			return true
		}
	}
	return false
}

func deriveSum(dto calibrationDto, operators []int) int {
	res, tail := dto.parts[0], dto.parts[1:]
	for i, each := range tail {
		if res > dto.sum {
			return res
		}
		switch operators[i] {
		case 0:
			res += each
		case 1:
			res *= each
		case 2:
			res = concat(res, each)
		default:
			panic(fmt.Sprintf("Unknown operator: %d.", operators[i]))
		}
	}
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
	length := shared.DigitLength(b)
	mul := shared.Pow(10, length)
	return a*mul + b
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

type permutationSpec struct {
	length int
	base   int
}

type permutationGenerator struct {
	m map[permutationSpec][][]int
}

func newPermutationGenerator() *permutationGenerator {
	return &permutationGenerator{
		m: map[permutationSpec][][]int{},
	}
}

func (v *permutationGenerator) generate(spec permutationSpec) [][]int {
	if existing, exists := v.m[spec]; exists {
		return existing
	}
	perm := generatePermutations(spec.length, spec.base)
	v.m[spec] = perm
	return perm
}
