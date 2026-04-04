package aoc2312

import (
	"fmt"
	"strconv"
	"strings"

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

func SumPermutations(lines []string) int {
	rows := parseLines(lines)
	var count int
	for _, each := range rows {
		count += countPermutations(each)
	}
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

func countPermutations(row springRow) int {
	shared.Logger.Info("Count permutations.", "row", row)
	var counter int
	cb := func(_ spring) {
		counter++
	}
	hypothesize(row.springs, row.groups, cb)
	return counter
}

func copySlice[S ~[]T, T any](s S) S {
	copied := make(S, len(s))
	copy(copied, s)
	return copied
}

func hypothesize(aSpring spring, groups []int, cb func(spring)) {
	cand := copySlice(aSpring)
	doHypothesize(cand, groups, cb, 0)
}

func doHypothesize(aSpring spring, groups []int, cb func(spring), i int) {
	for ; i < len(aSpring); i++ {
		c := aSpring[i]
		if c == condOperational {
			continue
		} else if c == condUnknown {
			for _, each := range []condition{condOperational, condDamaged} {
				cand := copySlice(aSpring)
				cand[i] = each
				doHypothesize(cand, groups, cb, i)
			}
			return
		} else if c == condDamaged {
			if len(groups) == 0 {
				return
			}
			end := i + groups[0]
			if end > len(aSpring) {
				return
			} else if end < len(aSpring) {
				switch aSpring[end] {
				case condDamaged:
					return
				case condUnknown, condOperational:
					aSpring[end] = condOperational
				default:
					panic("unknown cond, x81")
				}
			}
			if !fillSafelyWithDamaged(aSpring[i:end]) {
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

func fillSafelyWithDamaged(aSpring spring) bool {
	for i, c := range aSpring {
		switch c {
		case condOperational:
			return false
		case condUnknown, condDamaged:
			aSpring[i] = condDamaged
		default:
			panic("unknown cond")
		}
	}
	return true
}
