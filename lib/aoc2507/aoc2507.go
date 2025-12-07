package aoc2507

import (
	"fmt"
	"slices"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

const (
	charSpace    = '.'
	charSplitter = '^'
	charStart    = 'S'
)

func CountSplits(lines []string) int {
	board := shared.NewBoard(lines)
	beans := gent.NewSet(board.FindOrDie(charStart))
	var count int
	for {
		addedSplits, next := split(board, beans)
		count += addedSplits
		if next == nil {
			break
		}
		beans = next
	}
	return count
}

func split(board *shared.Board, beans *gent.Set[shared.Loc]) (int, *gent.Set[shared.Loc]) {
	var count int
	next := gent.NewSet[shared.Loc]()
	beans.ForEachAll(func(each shared.Loc) {
		each.Y--
		c, ok := board.Get(each)
		if !ok {
			return
		}
		if c == charSpace {
			next.Add(each)
			return
		}
		if c != charSplitter {
			return
		}
		count++
		for _, boy := range []shared.Loc{
			each.Delta(shared.Loc{X: -1}),
			each.Delta(shared.Loc{X: 1}),
		} {
			d, ok := board.Get(boy)
			if !ok {
				continue
			}
			if d == '.' {
				next.Add(boy)
			}
		}
	})
	if next.Count() == 0 {
		return count, nil
	}
	return count, next
}

func CountTimelines(lines []string) int {
	board := shared.NewBoard(lines)
	beans := map[shared.Loc]int{board.FindOrDie('S'): 1}
	for {
		if !stepInTimelines(board, beans) {
			break
		}
	}
	if shared.IsDebugEnabled() {
		keys := make([]shared.Loc, len(beans))
		for each := range beans {
			keys = append(keys, each)
		}
		slices.SortFunc(keys, func(a, b shared.Loc) int {
			x := a.X - b.X
			if x != 0 {
				return x
			}
			return a.Y - b.Y
		})
		for _, each := range keys {
			shared.Logger.Debug("Result.", "key", each, "value", beans[each])
		}
	}
	total := 0
	for _, v := range beans {
		total += v
	}
	return total
}

func stepInTimelines(board *shared.Board, beans map[shared.Loc]int) bool {
	var moved bool
	keys := make([]shared.Loc, len(beans))
	for loc := range beans {
		keys = append(keys, loc)
	}
	for _, each := range keys {
		stepped := each.Delta(shared.Loc{Y: -1})
		c, ok := board.Get(stepped)
		if !ok {
			continue
		}

		moved = true
		if c == charSpace {
			value := beans[each]
			if curr, ok := beans[stepped]; ok {
				value += curr
			}
			beans[stepped] = value
			delete(beans, each)
			continue
		}

		if c != charSplitter {
			panic(fmt.Sprintf("logic failure: %v - %v", each, c))
		}
		for _, boy := range []shared.Loc{
			each.Delta(shared.Loc{X: -1}),
			each.Delta(shared.Loc{X: 1})} {
			value := beans[each]
			if curr, ok := beans[boy]; ok {
				value += curr
			}
			beans[boy] = value
		}
		delete(beans, each)
	}
	return moved
}
