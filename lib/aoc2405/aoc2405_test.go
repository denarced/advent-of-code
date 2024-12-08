package aoc2405

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestSumCorrectMiddlePageNumbers(t *testing.T) {
	shared.InitTestLogging(t)
	// 143 is from the problem description.
	require.Equal(t, 143, SumCorrectMiddlePageNumbers(advent05Lines()))
}

func advent05Lines() []string {
	// Example values from problem description.
	return []string{
		"47|53",
		"97|13",
		"97|61",
		"97|47",
		"75|29",
		"61|13",
		"75|53",
		"29|13",
		"97|29",
		"53|29",
		"61|53",
		"97|53",
		"61|29",
		"47|13",
		"75|47",
		"97|75",
		"47|61",
		"75|61",
		"47|29",
		"75|13",
		"53|13",
		"75,47,61,53,29",
		"97,61,53,29,13",
		"75,29,13",
		"75,97,47,61,53",
		"61,13,29",
		"97,13,75,29,47",
	}
}

func TestSumIncorrectMiddlePageNumbers(t *testing.T) {
	shared.InitTestLogging(t)
	// 123 is from the problem description.
	require.Equal(t, 123, SumIncorrectMiddlePageNumbers(advent05Lines()))
}
