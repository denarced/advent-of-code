package aoc2503

import (
	"fmt"
	"sync"

	"github.com/denarced/advent-of-code/shared"
)

func DeriveMaxJoltageSum(lines []string, count int) int64 {
	banks := toBanks(lines)
	var sum int64
	var wg sync.WaitGroup
	for _, each := range banks {
		wg.Add(1)
		go func(bank []int) {
			defer wg.Done()
			joltage := deriveMaxJoltage(bank, count)
			sum += joltage
			shared.Logger.Info("Bank processed.", "bank", bank, "joltage", joltage, "sum", sum)
		}(each)
	}
	wg.Wait()
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
			cand := appendInts(d, rest)
			if shared.IsDebugEnabled() {
				shared.Logger.Debug("Candidate.", "candidate", cand, "rest", rest)
			}
			if cand > maxJoltage {
				maxJoltage = cand
			}
		}
	}
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

func appendInts(a, b int64) int64 {
	var m int64 = 1
	for temp := b; temp > 0; temp /= 10 {
		m *= 10
	}
	return a*m + b
}
