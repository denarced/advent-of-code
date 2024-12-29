package aoc2411

import (
	"fmt"
	"math"
	"sync"

	"github.com/denarced/advent-of-code/shared"
)

const (
	stateUnresolved = iota
	stateResolved
)

type blinkStone struct {
	value  int
	blinks int
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
	stone  blinkStone
	parent *node
	kids   []*node
	acc    int
	state  int
}

type stoneCache struct {
	m  map[blinkStone]int
	mu sync.Mutex
}

func (v *stoneCache) set(stone blinkStone, count int) {
	v.mu.Lock()
	if _, exists := v.m[stone]; !exists {
		v.m[stone] = count
	}
	v.mu.Unlock()
}

func (v *stoneCache) get(stone blinkStone) (value int, exists bool) {
	v.mu.Lock()
	value, exists = v.m[stone]
	v.mu.Unlock()
	return
}

func (v *stoneCache) size() int {
	v.mu.Lock()
	defer v.mu.Unlock()
	return len(v.m)
}

func CountStones(values []int, blinks int) int {
	cache := &stoneCache{
		m: map[blinkStone]int{},
	}

	resultCh := make(chan int)
	totalCh := make(chan int)
	go func() {
		total := 0
		for each := range resultCh {
			total += each
		}
		totalCh <- total
	}()

	var wg sync.WaitGroup
	for _, each := range values {
		wg.Add(1)
		go walkIntoStone(each, blinks, cache, resultCh, &wg)
	}
	wg.Wait()
	close(resultCh)
	count := <-totalCh
	shared.Logger.Info("Stones counted.", "count", count, "cache size", cache.size())
	return count
}

func walkIntoStone(
	value int,
	blinks int,
	cache *stoneCache,
	resultWrite chan<- int,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	root := &node{
		stone: blinkStone{
			value:  value,
			blinks: blinks,
		},
		parent: nil,
		kids:   []*node{},
		acc:    0,
		state:  stateUnresolved,
	}
	current := root
	for current != nil {
		shared.Assert(current.state != stateResolved, "resolved state is impossible")
		shared.Logger.Debug("Iterate within a stone path.", "current acc", current.acc)
		// At the bottom.
		if current.stone.blinks <= 0 {
			current = processLeaf(current, resultWrite)
			continue
		}

		shared.Logger.Debug("Processing branch.", "stone", current.stone)
		if len(current.kids) == 0 {
			if cached, exists := cache.get(current.stone); exists {
				shared.Logger.Debug(
					"Cache hit.",
					"stone", current.stone,
					"cached", cached,
				)
				current = processBranchWithCache(current, cached, resultWrite)
				continue
			}

			shared.Logger.Debug("Generate kids.", "stone", current.stone)
			current = processBranch(current)
			continue
		}

		shared.Assert(
			len(current.kids) > 0,
			"must have kids, all cases without them were already covered",
		)
		unresolved := findUnresolved(current.kids)
		if unresolved != nil {
			current = unresolved
			continue
		}
		current = finishBranch(current, cache)
	}
}

func findUnresolved(nodes []*node) *node {
	for _, each := range nodes {
		if each.state == stateUnresolved {
			return each
		}
	}
	return nil
}

func processLeaf(current *node, resultWrite chan<- int) *node {
	shared.Logger.Debug("Leaf reached, no blinks.", "value", current.stone.value)
	resultWrite <- 1
	current.state = stateResolved
	shared.Assert(len(current.kids) == 0, "leaf nodes don't have kids")
	// Just for consistency, it should be impossible to actually have kids within a leaf
	// node.
	current.kids = nil
	current.acc = 1
	if current.parent != nil {
		current.parent.acc += current.acc
	}
	return current.parent
}

func processBranchWithCache(current *node, cached int, resultWrite chan<- int) *node {
	current.state = stateResolved
	current.kids = nil
	current.acc = cached
	if current.parent != nil {
		current.parent.acc += current.acc
	}
	resultWrite <- cached
	return current.parent
}

func processBranch(current *node) *node {
	first, second, ok := transform(current.stone.value)
	left := &node{
		stone: blinkStone{
			value:  first,
			blinks: current.stone.blinks - 1,
		},
		parent: current,
		kids:   []*node{},
		acc:    0,
		state:  stateUnresolved,
	}
	kids := []*node{left}
	if ok {
		kids = append(kids, &node{
			stone: blinkStone{
				value:  second,
				blinks: current.stone.blinks - 1,
			},
			parent: current,
			kids:   []*node{},
			acc:    0,
			state:  stateUnresolved,
		})
	}
	current.kids = kids
	return left
}

func finishBranch(current *node, cache *stoneCache) *node {
	shared.Assert(
		current.acc > 0,
		"impossible to not have kids, they were already processed so acc should be >0",
	)
	// This parent node is done.
	current.state = stateResolved
	current.kids = nil
	cache.set(current.stone, current.acc)
	if current.parent != nil {
		current.parent.acc += current.acc
	}
	return current.parent
}
