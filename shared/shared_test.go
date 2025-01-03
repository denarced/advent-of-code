package shared

import (
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
