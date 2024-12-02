package aoc2024

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestToColumns(t *testing.T) {
	run := func(name string, lines, expectedLeft, expectedRight []string) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
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

func TestAdvent01Distance(t *testing.T) {
	run := func(name string, left, right []int, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual := Advent01Distance(left, right)
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
			similarity := Advent01Similarity(left, right)
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

func TestToInts(t *testing.T) {
	run := func(name string, strings []string, expected []int, errMessage string) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			actual, err := ToInts(strings)
			if errMessage == "" {
				req.Nil(err)
				req.Equal(expected, actual)
			} else {
				req.ErrorContains(err, errMessage)
			}
		})
	}

	run("empty", []string{}, nil, "")
	run("happy path", []string{"-1", "0", "1"}, []int{-1, 0, 1}, "")
	run("failure", []string{"e"}, nil, "invalid syntax")
}
