package shared

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToIntTable(t *testing.T) {
	run := func(name string, lines []string, expected [][]int) {
		t.Run(name, func(t *testing.T) {
			InitTestLogging(t)
			req := require.New(t)
			actual := ToIntTable(lines)
			req.Equal(expected, actual)
		})
	}

	run("nil", nil, nil)
	run("empty", []string{}, nil)
	run("1x1", []string{" 33 "}, [][]int{{33}})
	run("3x2", []string{"1 2", "2 3", "3 4"}, [][]int{{1, 2}, {2, 3}, {3, 4}})
}

func TestToColumns(t *testing.T) {
	run := func(name string, lines, expectedLeft, expectedRight []string) {
		t.Run(name, func(t *testing.T) {
			InitTestLogging(t)
			req := require.New(t)
			left, right := ToColumns(lines)
			req.Equal(expectedLeft, left)
			req.Equal(expectedRight, right)
		})
	}

	run("empty", []string{}, nil, nil)
	run("space", []string{"abc efg"}, []string{"abc"}, []string{"efg"})
	run("two spaces", []string{"313  666"}, []string{"313"}, []string{"666"})
}

func TestDigitLength(t *testing.T) {
	run := func(i, expected int) {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			InitTestLogging(t)
			req := require.New(t)
			actual := DigitLength(i)
			req.Equal(expected, actual)
		})
	}

	run(0, 1)
	run(1, 1)
	run(10, 2)
	run(19, 2)
	run(999, 3)
	run(1000, 4)
}

func TestModForIndex(t *testing.T) {
	run := func(dividend, divisor, expected int) {
		name := fmt.Sprintf("%d/%d", dividend, divisor)
		t.Run(name, func(t *testing.T) {
			req := require.New(t)

			// EXERCISE & VERIFY
			req.Equal(expected, ModForIndex(dividend, divisor))
		})
	}

	run(0, 0, 0)

	run(-3, 2, 1)
	run(-2, 2, 0)
	run(-1, 2, 1)
	run(0, 2, 0)
	run(1, 2, 1)
	run(2, 2, 0)
	run(3, 2, 1)
	run(4, 2, 0)
}

func TestDeriveGreeatestCommonDivisor(t *testing.T) {
	for _, each := range [][3]int{
		{1, 3, 5},
		{2, 2, 4},
		{5, 15, 20},
		{21, 252, 105},
	} {
		t.Run(fmt.Sprintf("%d and %d", each[1], each[2]), func(t *testing.T) {
			req := require.New(t)
			req.Equal(each[0], DeriveGreatestCommonDivisor(each[1], each[2]), "euclidean")
		})
	}
}

func TestDeriveLeastCommonMultiple(t *testing.T) {
	for _, each := range [][]int{
		{2, 1, 2},
		{12, 4, 6},
		{15, 3, 5},
		{60, 5, 2, 3, 4, 5},
	} {
		t.Run(fmt.Sprintf("%d and %d", each[1], each[2]), func(t *testing.T) {
			req := require.New(t)
			req.Equal(each[0], DeriveLeastCommonMultiple(each[1], each[2], each[3:]...))
		})
	}
}
