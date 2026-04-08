package aoc2312

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

const (
	condOperational condition = iota
	condDamaged
	condUnknown
)

type condition int

type spring []condition

func (v spring) String() string {
	result := make([]rune, len(v))
	for i, c := range v {
		switch c {
		case condDamaged:
			result[i] = '#'
		case condOperational:
			result[i] = '.'
		case condUnknown:
			result[i] = '?'
		default:
			panic("unknown condition: " + strconv.Itoa(int(c)))
		}
	}
	return string(result)
}

type springRow struct {
	springs spring
	groups  []int
}

func (v springRow) String() string {
	return fmt.Sprintf(
		"%s %s",
		v.springs,
		strings.Join(
			gent.Map(
				v.groups,
				func(i int) string {
					return strconv.Itoa(i)
				},
			),
			","))
}

func SumPermutations(lines []string, mul int) int {
	rows := parseLines(lines)
	var count int
	ch := make(chan int)
	closer := make(chan bool)
	go func() {
		for i := range ch {
			count += i
		}
		closer <- false
	}()
	var wg sync.WaitGroup
	for _, each := range rows {
		wg.Add(1)
		go func(row springRow) {
			defer wg.Done()
			ch <- countPermutations(row, mul)
		}(each)
	}
	wg.Wait()
	close(ch)
	<-closer
	return count
}

func parseLines(lines []string) []springRow {
	rows := make([]springRow, len(lines))
	for i, each := range lines {
		rows[i] = parseLine(each)
	}
	return rows
}

func parseLine(line string) springRow {
	fields := strings.Fields(line)
	if len(fields) != 2 {
		shared.Logger.Error("Line is invalid, not 2 fields in in.", "line", line)
		panic("invalid line")
	}
	var row springRow
	row.groups = gent.Map(
		strings.Split(fields[1], ","),
		func(s string) int {
			i, err := strconv.Atoi(s)
			if err != nil {
				shared.Logger.Error("Invalid line with non-int in group.", "line", line, "err", err)
				panic(err)
			}
			return i
		})
	row.springs = parseSpring(fields[0])
	return row
}

func parseSpring(s string) spring {
	aSpring := make(spring, len(s))
	for i, c := range s {
		switch c {
		case '?':
			aSpring[i] = condUnknown
		case '.':
			aSpring[i] = condOperational
		case '#':
			aSpring[i] = condDamaged
		default:
			panic("unknown character for spring condition: " + string(c))
		}
	}
	return aSpring
}

func createCounter() (*int, func(spring)) {
	var count int
	return &count, func(_ spring) {
		count++
	}
}

func countPermutations(row springRow, mul int) int {
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Count permutations.", "row", row, "multiplier", mul)
	}
	count, counter := createCounter()
	alpha := time.Now()
	hypothesize(multiplySpring(row.springs, mul), multiplyGroups(row.groups, mul), counter)
	shared.Logger.Info(
		"Permutations counted.",
		"count", *count,
		"row", row,
		"multiplier", mul,
		"duration", time.Since(alpha),
	)
	return *count
}

func multiplySpring(s spring, mul int) spring {
	if mul == 1 {
		return s
	}
	result := make(spring, mul*len(s)+mul-1)
	for n := range mul {
		i := n*len(s) + n
		copy(result[i:], s)
		if n < (mul - 1) {
			result[i+len(s)] = condUnknown
		}
	}
	return result
}

func multiplyGroups(groups []int, mul int) []int {
	if mul <= 1 {
		return groups
	}
	result := make([]int, mul*len(groups))
	for i := range mul {
		copy(result[i*len(groups):], groups)
	}
	return result
}

type condPair struct {
	damaged, operational int
}

type condCounter struct {
	target               condPair
	status               condPair
	seenOperationalCount int
}

func (v *condCounter) can(c condition, count int) bool {
	switch c {
	case condDamaged:
		return v.status.damaged+count <= v.target.damaged
	case condOperational:
		return v.status.operational+count <= v.target.operational
	default:
		return false
	}
}

func (v *condCounter) add(c condition, count int) {
	if !v.can(c, count) {
		panic("can't add")
	}
	if c == condDamaged {
		v.status.damaged += count
	} else if c == condOperational {
		v.status.operational += count
	} else {
		panic("impossible: invalid condition type to add")
	}
}

func createCondCounter(aSpring spring, groups []int) *condCounter {
	counter := &condCounter{}
	counter.target.damaged = sumInts(groups)
	counter.target.operational = len(aSpring) - counter.target.damaged
	for _, each := range aSpring {
		switch each {
		case condDamaged:
			counter.status.damaged++
		case condOperational:
			counter.status.operational++
		default:
			// Nothing to do.
		}
	}
	return counter
}

func createPool(springSize int) *sync.Pool {
	return &sync.Pool{
		New: func() any {
			s := make(spring, springSize)
			return &s
		},
	}
}

func hypothesize(aSpring spring, groups []int, cb func(spring)) {
	pool := createPool(len(aSpring))
	cand := copySpring(pool, &aSpring)
	counter := createCondCounter(*cand, groups)

	hypoMach := createHypothesisMachine(cb, pool)
	hypoMach.hypothesize(cand, groups, 0, counter)
}

func copySpring(pool *sync.Pool, aSpring *spring) *spring {
	next, ok := pool.Get().(*spring)
	if !ok {
		panic("pool type check failed")
	}
	n := copy(*next, *aSpring)
	if n != len(*aSpring) {
		panic("count of copied items doesn't match")
	}
	return next
}

func sumInts(ints []int) (sum int) {
	for _, each := range ints {
		sum += each
	}
	return sum
}

type hypothesisMachine struct {
	cb   func(spring)
	pool *sync.Pool
}

func createHypothesisMachine(cb func(spring), pool *sync.Pool) *hypothesisMachine {
	return &hypothesisMachine{cb: cb, pool: pool}
}

func (v *hypothesisMachine) hypothesize(
	aSpringPtr *spring,
	groups []int,
	i int,
	counter *condCounter,
) {
	defer v.pool.Put(aSpringPtr)
	aSpring := *aSpringPtr
	for ; i < len(aSpring); i++ {
		c := aSpring[i]
		switch c {
		case condOperational:
			counter.seenOperationalCount++
			continue
		case condUnknown:
			if len(groups) == 0 {
				if !counter.can(condOperational, 1) {
					return
				}
				counter.add(condOperational, 1)
				aSpring[i] = condOperational
				counter.seenOperationalCount++
				continue
			}

			// Operational
			if (counter.target.operational-counter.seenOperationalCount-(len(groups)-1)) > 0 &&
				counter.can(condOperational, 1) {
				altered := *counter
				altered.add(condOperational, 1)
				cand := copySpring(v.pool, aSpringPtr)
				(*cand)[i] = condOperational
				v.hypothesize(cand, groups, i, &altered)
			}

			// Damaged
			unknownIndexes, ok := v.validateDamageToAdd(i, groups, aSpringPtr)
			if !ok || !counter.can(condDamaged, len(unknownIndexes)) {
				return
			}
			cand := copySpring(v.pool, aSpringPtr)
			altered := *counter
			altered.add(condDamaged, len(unknownIndexes))
			for _, j := range unknownIndexes {
				(*cand)[j] = condDamaged
			}
			nextIndex := i + groups[0]
			if i+groups[0] < len(aSpring) && aSpring[i+groups[0]] == condUnknown {
				nextIndex++
				(*cand)[i+groups[0]] = condOperational
				if !altered.can(condOperational, 1) {
					return
				}
				altered.add(condOperational, 1)
			}
			v.hypothesize(cand, groups[1:], nextIndex, &altered)
			return
		case condDamaged:
			unknownIndexes, ok := v.validateDamageToAdd(i, groups, aSpringPtr)
			if !ok {
				return
			}
			if !counter.can(condDamaged, len(unknownIndexes)) {
				return
			}
			counter.add(condDamaged, len(unknownIndexes))
			for _, j := range unknownIndexes {
				aSpring[j] = condDamaged
			}
			nextIndex := i + groups[0] - 1
			if i+groups[0] < len(aSpring) && aSpring[i+groups[0]] == condUnknown {
				nextIndex++
				aSpring[i+groups[0]] = condOperational
			}
			groups = groups[1:]
			i = nextIndex
		default:
			panic("unknown type of cond")
		}
	}
	if len(groups) == 0 {
		v.cb(aSpring)
	}
}

func (v *hypothesisMachine) validateDamageToAdd(
	i int,
	groups []int,
	aSpringPtr *spring,
) (indexes []int, ok bool) {
	if len(groups) == 0 {
		return
	}
	if i+groups[0] > len(*aSpringPtr) {
		return
	}
	if i+groups[0] < len(*aSpringPtr) && (*aSpringPtr)[i+groups[0]] == condDamaged {
		return
	}
	indexes = make([]int, 0, groups[0])
	for j := i; j < i+groups[0]; j++ {
		switch (*aSpringPtr)[j] {
		case condOperational:
			indexes = nil
			return
		case condUnknown:
			indexes = append(indexes, j)
		default:
			// Nothing to do.
		}
	}
	ok = true
	return
}
