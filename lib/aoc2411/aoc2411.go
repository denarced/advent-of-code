package aoc2411

import (
	"fmt"
	"math"
	"slices"

	"github.com/denarced/advent-of-code/shared"
)

type spottedStone struct {
	value int
	spots int
}

func CountStones(blinkCount int, stoneValues []int) int {
	shared.Logger.Info("Count stones.", "blink count", blinkCount, "stone count", len(stoneValues))

	stones := []spottedStone{}
	stoneToCount := map[spottedStone]int{}
	for _, each := range stoneValues {
		stones = append(stones, spottedStone{value: each, spots: blinkCount})
		if blinkCount > 2 {
			countNodes(buildPartialTree(each, blinkCount), stoneToCount)
		}
	}
	shared.Logger.Info("Tree built.", "size", len(stoneToCount))
	cycle := 10_000_000
	round := cycle
	totalCount := 0
	halfBlink := blinkCount / 2
	getCachedCount := func(each *spottedStone) int {
		if 0 < each.spots && each.spots < halfBlink {
			if count, ok := stoneToCount[*each]; ok {
				return count
			}
		}
		return 0
	}
mainLoop:
	for len(stones) > 0 {
		round--
		each := stones[0]
		stones = stones[1:]
		if cachedCount := getCachedCount(&each); cachedCount > 0 {
			totalCount += cachedCount
			continue
		}
		for range each.spots {
			if cachedCount := getCachedCount(&each); cachedCount > 0 {
				totalCount += cachedCount
				continue mainLoop
			}

			first, second, cloned := transform(each.value)
			each.spots--
			each.value = first
			if cloned {
				// Prepend so newer stones with less steps are processed first. That way memory
				// usage is reduced significantly because we remove stones quicker from "stones"
				// slice.
				stones = slices.Insert(stones, 0, spottedStone{value: second, spots: each.spots})
			}
		}
		totalCount++
		if round <= 0 {
			round = cycle
			shared.Logger.Info("Sort.", "total count", totalCount, "stone count", len(stones))
		}
	}
	return totalCount
}

func transform(stone int) (first int, second int, cloned bool) {
	if stone == 0 {
		first = 1
		return
	}
	first, second, cloned = splitStone(stone)
	if !cloned {
		first = 2024 * stone
		if first < 2024 {
			panic(fmt.Sprintf("int overflow: %d*%d==%d", 2024, stone, first))
		}
	}
	return
}

func splitStone(stone int) (first int, second int, ok bool) {
	length := int(math.Log10(float64(stone))) + 1
	if length%2 != 0 {
		return
	}
	div := shared.Pow(10, length/2)
	ok = true
	first = stone / div
	second = stone % div
	return
}

type node struct {
	value  int
	blinks int
	parent *node
	kids   []*node
}

func countNodes(nod *node, pre map[spottedStone]int) (sum int) {
	if nod == nil {
		return
	}
	if len(nod.kids) == 0 {
		sum = 1
		return
	}
	for _, each := range nod.kids {
		sum += countNodes(each, pre)
	}
	spotted := spottedStone{
		value: nod.value,
		spots: nod.blinks,
	}
	if pre != nil {
		pre[spotted] = sum
	}
	return
}

func buildPartialTree(stone, blinks int) *node {
	var half *node
	{
		current := &node{value: stone, blinks: blinks}
		for current.blinks > blinks/2 {
			first, _, _ := transform(current.value)
			left := &node{value: first, blinks: current.blinks - 1, parent: current}
			current.kids = []*node{left}
			current = left
		}
		half = current
	}

	stack := []*node{half}
	for len(stack) > 0 {
		each := stack[0]
		stack = stack[1:]
		if each.blinks <= 0 {
			if len(each.kids) > 0 {
				panic("what? kids?")
			}
			continue
		}
		first, second, cloned := transform(each.value)
		kids := []*node{}
		left := &node{value: first, blinks: each.blinks - 1, parent: each}
		stack = append(stack, left)
		kids = append(kids, left)
		if cloned {
			right := &node{value: second, blinks: each.blinks - 1, parent: each}
			kids = append(kids, right)
			stack = append(stack, right)
		}
		each.kids = kids
	}
	return half
}
