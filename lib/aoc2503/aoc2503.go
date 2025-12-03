package aoc2503

import (
	"fmt"
	"strconv"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func DeriveMaxJoltageSum(lines []string, count int) int64 {
	banks := toBanks(lines)
	var sum int64
	for _, each := range banks {
		joltage := deriveMaxJoltage(each, count)
		sum += joltage
		shared.Logger.Info("Bank processed.", "bank", each, "joltage", joltage, "sum", sum)
	}
	return sum
}

func deriveMaxJoltage(bank []int, count int) int64 {
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Derive max joltage.", "bank", bank, "count", count)
	}
	if count < 2 {
		var maximum int64 = -1
		for _, each := range bank {
			maximum = max(maximum, int64(each))
		}
		return maximum
	}
	var maxMajor int64
	var maxJoltage int64
	maxIndex := len(bank) - count + 1
	for i := range maxIndex {
		if len(bank[i+1:]) < count-1 {
			break
		}
		d := int64(bank[i])
		if d > maxMajor {
			maxMajor = d
			maxJoltage = 0
			rest := deriveMaxJoltage(bank[i+1:], count-1)
			cand := gent.OrPanic2(
				strconv.ParseInt(fmt.Sprintf("%d%d", d, rest), 10, 64),
			)(
				"combine numbers",
			)
			if shared.IsDebugEnabled() {
				shared.Logger.Debug("Candidate.", "candidate", cand, "rest", rest)
			}
			if cand > maxJoltage {
				maxJoltage = cand
			}
		}
	}
	shared.Logger.Info("Max joltage derived.", "joltage", maxJoltage, "major", maxMajor)
	if maxJoltage > 0 {
		return maxJoltage
	}
	return maxMajor
}

func toBanks(lines []string) [][]int {
	var banks [][]int
	for _, each := range lines {
		var bank []int
		for _, c := range each {
			value := int(c - '0')
			bank = append(bank, value)
		}
		if bank != nil {
			if len(bank) < 2 {
				panic(fmt.Sprintf("invalid bank: %s", each))
			}
			banks = append(banks, bank)
		}
	}
	return banks
}
