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

func copySlice[S ~[]T, T any](s S) S {
	copied := make(S, len(s))
	copy(copied, s)
	return copied
}

type condPair struct {
	damaged, operational int
}

type condCounter struct {
	target               condPair
	status               condPair
	seenOperationalCount int
}

func (v *condCounter) add(c condition, count int) bool {
	switch c {
	case condDamaged:
		if v.status.damaged+count > v.target.damaged {
			return false
		}
		v.status.damaged += count
	case condOperational:
		if v.status.operational+count > v.target.operational {
			return false
		}
		v.status.operational += count
	default:
		return false
	}
	return true
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
	doHypothesize(cand, groups, cb, 0, counter, pool)
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

//revive:disable-next-line:cyclomatic,cognitive-complexity,function-length
func doHypothesize(
	aSpringPtr *spring,
	groups []int,
	cb func(spring),
	i int,
	counter *condCounter,
	pool *sync.Pool,
) {
	defer pool.Put(aSpringPtr)
	aSpring := *aSpringPtr
	for ; i < len(aSpring); i++ {
		c := aSpring[i]
		if c == condOperational {
			counter.seenOperationalCount++
			continue
		} else if c == condUnknown {
			if len(groups) == 0 {
				if !counter.add(condOperational, 1) {
					return
				}
				aSpring[i] = condOperational
				counter.seenOperationalCount++
				continue
			}
			var options []condition
			if (counter.target.operational - counter.seenOperationalCount - (len(groups) - 1)) > 0 {
				options = append(options, condOperational)
			}
			options = append(options, condDamaged)
			for _, each := range options {
				altered := *counter
				if !altered.add(each, 1) {
					continue
				}
				cand := copySpring(pool, aSpringPtr)
				(*cand)[i] = each
				doHypothesize(cand, groups, cb, i, &altered, pool)
			}
			return
		} else if c == condDamaged {
			if len(groups) == 0 {
				return
			}
			end := i + groups[0]
			if end > len(aSpring) {
				return
				//revive:disable-next-line:max-control-nesting
			} else if end < len(aSpring) {
				switch aSpring[end] {
				case condDamaged:
					return
				case condUnknown:
					if !counter.add(condOperational, 1) {
						return
					}
					aSpring[end] = condOperational
					counter.seenOperationalCount++
				case condOperational:
					counter.seenOperationalCount++
				default:
					panic("unknown cond, x81")
				}
			}
			success, added := fillSafelyWithDamaged(aSpring[i:end])
			if !success {
				return
			}
			if !counter.add(condDamaged, added) {
				return
			}
			i = end
			groups = groups[1:]
		}
	}
	if len(groups) == 0 {
		cb(aSpring)
	}
}

func fillSafelyWithDamaged(aSpring spring) (success bool, added int) {
	success = true
mainLoop:
	for i, c := range aSpring {
		switch c {
		case condOperational:
			success = false
			break mainLoop
		case condUnknown:
			added++
			aSpring[i] = condDamaged
		case condDamaged:
			// Nothing to do.
		default:
			panic("unknown cond")
		}
	}
	return
}

func pow(base, exp int) int {
	if exp == 0 {
		return 1
	}
	result := base
	for range exp - 1 {
		result *= base
	}
	return result
}
