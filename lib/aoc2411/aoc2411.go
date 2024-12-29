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

type spottedStone struct {
	value int
	spots int
}

type spottedValue struct {
	stone spottedStone
	value int
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
	acc    int
	state  int
}

type stoneCache struct {
	m  map[spottedStone]int
	mu sync.Mutex
}

func (v *stoneCache) set(stone spottedStone, count int) {
	v.mu.Lock()
	if _, exists := v.m[stone]; !exists {
		v.m[stone] = count
	}
	v.mu.Unlock()
}

func (v *stoneCache) get(stone spottedStone) (value int, exists bool) {
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
	cacheCh := make(chan spottedValue)
	cache := &stoneCache{
		m: map[spottedStone]int{},
	}
	go func() {
		for each := range cacheCh {
			cache.set(each.stone, each.value)
		}
	}()

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
		go walkIntoStone(each, blinks, cache, cacheCh, resultCh, &wg)
	}
	wg.Wait()
	close(cacheCh)
	close(resultCh)
	count := <-totalCh
	shared.Logger.Info("Stones counted.", "count", count, "cache size", cache.size())
	return count
}

func walkIntoStone(
	value int,
	blinks int,
	cache *stoneCache,
	cacheWrite chan<- spottedValue,
	resultWrite chan<- int,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	root := &node{
		value:  value,
		blinks: blinks,
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
		if current.blinks <= 0 {
			shared.Logger.Debug("Leaf reached, no blinks.", "value", current.value)
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
			current = current.parent
			continue
		}

		shared.Logger.Debug("Processing branch.", "blinks", current.blinks, "value", current.value)
		if len(current.kids) == 0 {
			if cached, exists := cache.get(spottedStone{
				value: current.value,
				spots: current.blinks,
			}); exists {
				shared.Logger.Debug(
					"Cache hit.",
					"blinks",
					current.blinks,
					"value",
					current.value,
					"cached",
					cached,
				)
				current.state = stateResolved
				current.kids = nil
				current.acc = cached
				if current.parent != nil {
					current.parent.acc += current.acc
				}
				resultWrite <- cached
				current = current.parent
				continue
			}

			shared.Logger.Debug("Generate kids.", "value", current.value, "blinks", current.blinks)
			first, second, ok := transform(current.value)
			left := &node{
				value:  first,
				blinks: current.blinks - 1,
				parent: current,
				kids:   []*node{},
				acc:    0,
				state:  stateUnresolved,
			}
			kids := []*node{left}
			if ok {
				kids = append(kids, &node{
					value:  second,
					blinks: current.blinks - 1,
					parent: current,
					kids:   []*node{},
					acc:    0,
					state:  stateUnresolved,
				})
			}
			current.kids = kids
			current = left
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
		shared.Assert(
			current.acc > 0,
			"impossible to not have kids, they were already processed so acc should be >0",
		)
		// This parent node is done.
		current.state = stateResolved
		current.kids = nil
		cachedValue := spottedValue{
			stone: spottedStone{
				value: current.value,
				spots: current.blinks,
			},
			value: current.acc,
		}
		cacheWrite <- cachedValue
		if current.parent != nil {
			current.parent.acc += current.acc
		}
		current = current.parent
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
