package aoc2507

import (
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
	var next []shared.Loc
	beans.ForEachAll(func(each shared.Loc) {
		each.Y--
		c, ok := board.Get(each)
		if !ok {
			return
		}
		if c == charSpace {
			next = append(next, each)
			return
		}
		if c == charSplitter {
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
					next = append(next, boy)
				}
			}
		}
	})
	if len(next) == 0 {
		return count, nil
	}
	return count, gent.NewSet(next...)
}
