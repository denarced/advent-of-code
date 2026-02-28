package aoc2301

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	expectedSuccess expectedStatus = iota
	expectedFailure
)

type expectedStatus int

func TestParseDigit(t *testing.T) {
	for _, each := range []struct {
		value    string
		expected int
	}{
		{"a1b3z2cd", 12},
		{"a3cd", 33},
	} {
		t.Run(each.value, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, each.expected, parseDigit(each.value, true))
		})
	}
}

func TestParseDigitIn(t *testing.T) {
	run := func(
		name, value string,
		expected int,
		prefix seekTarget,
		justDigits bool,
		shouldFind expectedStatus) {
		t.Run(name, func(t *testing.T) {
			ass := assert.New(t)
			d, found := parseDigitIn(value, prefix, justDigits)
			if found && shouldFind == expectedFailure {
				ass.Fail("should have failed")
			} else if !found && shouldFind == expectedSuccess {
				ass.Fail("should have succeeded")
			}
			ass.Equal(expected, d)
		})
	}
	run("pre-digit plus word", "1two", 1, prefixTarget, false, expectedSuccess)
	run("suffix-digit plus prefix-digit", "1two2", 2, suffixTarget, false, expectedSuccess)
	run("prefix-digit", "1", 1, prefixTarget, false, expectedSuccess)
	run("suffix-digit", "6", 6, suffixTarget, false, expectedSuccess)
	run("prefix-word", "seven", 7, prefixTarget, false, expectedSuccess)
	run("suffix-word", "six", 6, suffixTarget, false, expectedSuccess)
	run("prefix-word plus garble", "two#dood#", 2, prefixTarget, false, expectedSuccess)
	run("suffix-word with garble", "two#dood#", 0, suffixTarget, false, expectedFailure)
	run("prefix overlapping", "eightwo", 8, prefixTarget, false, expectedSuccess)
	run("suffix overlapping", "eightwo", 2, suffixTarget, false, expectedSuccess)
	run("justDigits prefix word", "one", 0, prefixTarget, true, expectedFailure)
	run("justDigits prefix digit", "1one", 1, prefixTarget, true, expectedSuccess)
}
