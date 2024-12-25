package aoc2417

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

const (
	comboRegisterA = 4
	comboRegisterB = 5
	comboRegisterC = 6
)

func TestDeriveOutput(t *testing.T) {
	run := func(name string, lines []string, expected string) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			// EXERCISE
			actual := DeriveOutput(lines)

			// VERIFY
			require.Equal(t, expected, actual)
		})
	}

	run("empty", []string{}, "")
	run(
		"example",
		[]string{
			"Register A: 729",
			"Register B: 0",
			"Register C: 0",
			"Program: 0,1,5,4,3,0",
		},
		"4,6,3,5,6,3,5,2,1,0")
}

func TestNewProcessor(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	actual := newProcessor([]string{
		"Register A: 729",
		"Register B: 1001",
		"Register C: 666",
		"Program: 0,1,5,4,3,0",
	})
	req.Equal([]int{729, 1001, 666}, actual.registers)
	req.Equal(0, actual.index)
	req.Equal([]int{0, 1, 5, 4, 3, 0}, actual.feed)
}

func TestProcess(t *testing.T) {
	run := func(name string, initial *processor, expected *processor) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			finished := initial.process()

			// VERIFY
			req.True(finished)
			req.Equal(expected, initial)
		})
	}

	run(
		"adv with literal operand",
		&processor{registers: []int{20}, feed: []int{0, 2}},
		&processor{registers: []int{5}, index: 2, feed: []int{0, 2}},
	)
	run(
		"adv with register C",
		&processor{registers: []int{50, 0, 4}, feed: []int{0, comboRegisterC}},
		&processor{
			registers: []int{50 / (2 * 2 * 2 * 2), 0, 4},
			index:     2,
			feed:      []int{0, comboRegisterC},
		},
	)

	// 6 = 110
	// 3 = 011
	// xor
	// 5 = 101
	run(
		"bxl",
		&processor{registers: []int{0, 3}, feed: []int{1, 6}},
		&processor{registers: []int{0, 5}, index: 2, feed: []int{1, 6}},
	)
	run(
		"bst with literal operand",
		&processor{registers: []int{0, 0}, feed: []int{2, 3}},
		&processor{registers: []int{0, 3}, index: 2, feed: []int{2, 3}},
	)
	run(
		"bst with register A",
		&processor{registers: []int{54_000_001, 0}, feed: []int{2, comboRegisterA}},
		&processor{
			registers: []int{54_000_001, 54_000_001 % 8},
			index:     2,
			feed:      []int{2, comboRegisterA},
		},
	)
	run(
		"jnz - do nothing",
		&processor{registers: []int{0, 3}, feed: []int{3, 6}},
		&processor{registers: []int{0, 3}, index: 2, feed: []int{3, 6}},
	)
	run(
		"jnz - jump",
		&processor{registers: []int{1, 0}, feed: []int{3, 6}},
		&processor{registers: []int{1, 0}, index: 6, feed: []int{3, 6}},
	)
	run(
		"bxc",
		&processor{registers: []int{0, 402, 700}, feed: []int{4, 0}},
		&processor{registers: []int{0, 402 ^ 700, 700}, index: 2, feed: []int{4, 0}},
	)
	run(
		"out with literal operand",
		&processor{feed: []int{5, 1}},
		&processor{output: []int{1}, index: 2, feed: []int{5, 1}},
	)
	run(
		"out with register B",
		&processor{registers: []int{0, 64}, feed: []int{5, comboRegisterB}},
		&processor{
			output:    []int{64 % 8},
			index:     2,
			registers: []int{0, 64},
			feed:      []int{5, comboRegisterB},
		},
	)
	run(
		"bdv with literal operand",
		&processor{registers: []int{33, 0}, feed: []int{6, 3}},
		&processor{registers: []int{33, 33 / (2 * 2 * 2)}, index: 2, feed: []int{6, 3}},
	)
	run(
		"bdv with register B",
		&processor{registers: []int{33, 5}, feed: []int{6, comboRegisterB}},
		&processor{registers: []int{33, 1}, index: 2, feed: []int{6, comboRegisterB}},
	)
	run(
		"cdv with literal operand",
		&processor{registers: []int{33, 0, 0}, feed: []int{7, 3}},
		&processor{registers: []int{33, 0, 33 / (2 * 2 * 2)}, index: 2, feed: []int{7, 3}},
	)
	run(
		"cdv with register C",
		&processor{registers: []int{9, 0, 1}, feed: []int{7, comboRegisterC}},
		&processor{registers: []int{9, 0, 4}, index: 2, feed: []int{7, comboRegisterC}},
	)
}
