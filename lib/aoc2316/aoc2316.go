package aoc2316

import (
	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

const (
	tileSpace              = '.'
	tileSplitterVertical   = '|'
	tileSplitterHorizontal = '-'
	tileMirrorDown         = '\\'
	tileMirrorUp           = '/'
)

type bean struct {
	loc shared.Loc
	dir shared.Direction
}

func FindMaxEnergizedTileCount(lines []string) int {
	shared.Logger.Info("Find start location with most energized tiles.")
	brd := shared.NewBoard(lines)
	maxCount := -1
	try := func(x, y int, dir shared.Direction) {
		count := CountEnergizedTilesWithStart(brd, shared.Loc{X: x, Y: y}, dir)
		maxCount = max(maxCount, count)
	}
	for _, x := range []int{-1, brd.GetWidth()} {
		for y := range brd.GetHeight() {
			dir := shared.RealWest
			if x < 0 {
				dir = shared.RealEast
			}
			try(x, y, dir)
		}
	}
	for _, y := range []int{-1, brd.GetHeight()} {
		for x := range brd.GetWidth() {
			dir := shared.RealSouth
			if y < 0 {
				dir = shared.RealNorth
			}
			try(x, y, dir)
		}
	}
	shared.Logger.Info("Maximum tile count found.", "count", maxCount)
	return maxCount
}

func CountEnergizedTiles(lines []string) int {
	brd := shared.NewBoard(lines)
	// Need to start outside the board because the first cell might be a mirror. The logic would in
	// that case if this starts with X=0.
	return CountEnergizedTilesWithStart(
		brd,
		shared.Loc{X: -1, Y: brd.GetHeight() - 1},
		shared.RealEast,
	)
}

func CountEnergizedTilesWithStart(
	brd *shared.Board,
	startLoc shared.Loc,
	startDir shared.Direction,
) int {
	shared.Logger.Info(
		"Count energized tiles.",
		"width", brd.GetWidth(),
		"height", brd.GetHeight(),
		"start location", startLoc,
		"start direction", startDir)
	aBean := bean{loc: startLoc, dir: startDir}
	energonized := gent.NewSet[bean]()
	var backlog []bean
	pickFromBacklog := func() (bean, bool) {
		if len(backlog) == 0 {
			shared.Logger.Debug("Nothing in backlog, stopping.")
			return bean{}, false
		}
		aBean = backlog[0]
		backlog = backlog[1:]
		shared.Logger.Debug(
			"Take bean from backlog.",
			"bean", aBean,
			"remaining in backlog", len(backlog))
		return aBean, true
	}
	first := true
	for {
		shared.Logger.Debug("Move.", "bean", aBean)
		// Start cell is illegal so it shouldn't be added.
		if !first && !energonized.Add(aBean) {
			shared.Logger.Debug("Been there, done that.", "bean", aBean)
			var ok bool
			aBean, ok = pickFromBacklog()
			if !ok {
				break
			}
			continue
		}
		first = false
		nextLoc := aBean.loc.Delta(shared.Loc(aBean.dir))
		c, ok := brd.Get(nextLoc)
		if !ok {
			shared.Logger.Debug(
				"Next location outside board.",
				"previous", aBean.loc,
				"next", nextLoc,
			)
			var ok bool
			aBean, ok = pickFromBacklog()
			if !ok {
				break
			}
			continue
		}
		switch c {
		case tileSpace:
			aBean.loc = nextLoc
		case tileSplitterVertical, tileSplitterHorizontal:
			aBean.loc = nextLoc
			if c == tileSplitterHorizontal &&
				(aBean.dir == shared.RealEast || aBean.dir == shared.RealWest) {
				continue
			}
			if c == tileSplitterVertical &&
				(aBean.dir == shared.RealNorth || aBean.dir == shared.RealSouth) {
				continue
			}
			firstDir := aBean.dir.TurnRealLeft()
			secDir := aBean.dir.TurnRealRight()
			aBean.dir = firstDir
			aBean.loc = nextLoc

			backlog = append(backlog, bean{
				loc: nextLoc,
				dir: secDir})
		case tileMirrorUp:
			switch aBean.dir {
			case shared.RealNorth, shared.RealSouth:
				aBean.dir = aBean.dir.TurnRealRight()
			case shared.RealEast, shared.RealWest:
				aBean.dir = aBean.dir.TurnRealLeft()
			default:
				panic("unknown direction for mirror up")
			}
			aBean.loc = nextLoc
		case tileMirrorDown:
			switch aBean.dir {
			case shared.RealNorth, shared.RealSouth:
				aBean.dir = aBean.dir.TurnRealLeft()
			case shared.RealEast, shared.RealWest:
				aBean.dir = aBean.dir.TurnRealRight()
			default:
				panic("unknown direction for mirror down")
			}
			aBean.loc = nextLoc
		default:
			panic("unknown kind of tile: " + string(c))
		}
	}

	locations := gent.NewSet[shared.Loc]()
	energonized.ForEachAll(func(each bean) {
		locations.Add(each.loc)
	})

	count := locations.Count()
	shared.Logger.Info("Energized tiles counted.", "count", count)
	return count
}
