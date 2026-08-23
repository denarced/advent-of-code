package aoc2206

import (
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

func DeriveStartOfPackage(line string, seekPackage bool) int {
	m := map[rune]int{}
	space := []rune(strings.TrimSpace(line))
	shared.Logger.Info(
		"Derive start of package.",
		"total character count", len(space),
		"seek package", seekPackage)
	overCount := 0
	length := 14
	if seekPackage {
		length = 4
	}
	var i int
	for i = range length {
		m[space[i]]++
		if m[space[i]] == 2 {
			overCount++
		}
	}
	shared.Logger.Debug("Initial mapping done.", "over count", overCount)
	i++
	defer func() {
		shared.Logger.Info("Start of package found.", "character count", i)
	}()
	if overCount == 0 {
		return i
	}
	for ; i < len(space); i++ {
		dropped := space[i-length]
		if m[dropped] == 2 {
			overCount--
		}
		m[dropped]--

		added := space[i]
		if m[added] == 1 {
			overCount++
		}
		m[added]++

		if overCount == 0 {
			return i + 1
		}
	}
	return i + 1
}
