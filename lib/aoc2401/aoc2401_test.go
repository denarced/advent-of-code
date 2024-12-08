package aoc2401

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestAdvent01Distance(t *testing.T) {
	run := func(name string, left, right []int, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual := Distance(left, right)
			req.Equal(expected, actual)
		})
	}

	run("empty", []int{}, []int{}, 0)
	run("happy path", []int{11, 7, 1}, []int{2, 11, 5}, 3)
}

func TestAdvent01Similarity(t *testing.T) {
	run := func(name string, left, right []int, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			similarity := Similarity(left, right)
			req.Equal(expected, similarity)
		})
	}

	run("empty", []int{}, []int{}, 0)
	run("empty right", []int{1, 2, 3}, []int{}, 0)
	run("empty left", []int{}, []int{1, 2, 2, 3}, 0)
	run("no shared", []int{1, 2, 2, 3, 3, 3}, []int{6, 6, 6}, 0)
	run("one shared", []int{3, 1, 3}, []int{3, 4, 5}, 6)
	run("multiple shared", []int{2, 3, 4}, []int{2, 2, 3, 3, 3, 5}, 4+9)
}
