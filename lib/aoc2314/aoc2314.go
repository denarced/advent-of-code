package aoc2314

import (
	"github.com/denarced/advent-of-code/shared"
)

const (
	roundRock  rune = 'O'
	emptySpace rune = '.'
)

func CountTotalLoad(lines []string, cycleCount int) int {
	shared.Logger.Info("Count total load.")
	brd := shared.NewBoard(lines)
	if cycleCount <= 0 {
		moveRocks(brd, shared.RealNorth)
	} else {
		hashPool := make([]int, 3*brd.GetArea())
		hashTrack := make([]uint64, 0, 1100)
		var loopLength int
		loopCounter := 3
		for cycleCount > 0 {
			cycleCount--
			for _, each := range []shared.Direction{
				shared.RealNorth,
				shared.RealWest,
				shared.RealSouth,
				shared.RealEast,
			} {
				moveRocks(brd, each)
			}
			if loopCounter <= 0 {
				continue
			}
			rockHash := hashRocks(brd, hashPool)
			hashTrack = append(hashTrack, rockHash)
			if distance, found := findLoop(hashTrack); found {
				if loopLength == 0 {
					loopLength = distance
				} else {
					if loopLength != distance {
						panic("loop length is not supposeed to change")
					}
					loopCounter--
					if loopCounter <= 0 {
						shared.Logger.Info(
							"Loop length confirmed, about to run out cycles.",
							"cycle count", cycleCount,
							"loop length", loopLength)
						var c int
						for c = cycleCount; c >= 0; c -= loopLength {
							cycleCount = c
						}
						shared.Logger.Info("Cycles consumed.", "cycle count", cycleCount)
					}
				}
			}
			count := len(hashTrack)
			if count >= 1100 {
				hashTrack = hashTrack[count-1000 : count]
			}
		}
	}
	var weight int
	shared.Logger.Info("Count weight.")
	brd.Iter(func(loc shared.Loc, c rune) (keepGoing bool) {
		keepGoing = true
		if c != roundRock {
			return
		}
		weight += loc.Y + 1
		return
	})
	shared.Logger.Info("Weight counted.", "weight", weight)
	return weight
}

func moveRocks(brd *shared.Board, direction shared.Direction) {
	var moveCount, rockCount int
	feedRocks(brd, direction, func(each shared.Loc) {
		var moved bool
		dest := each
		for {
			cand := dest.Delta(shared.Loc(direction))
			if c, ok := brd.Get(cand); ok && c == emptySpace {
				dest = cand
				moved = true
				moveCount++
				continue
			}
			break
		}
		if !moved {
			return
		}
		rockCount++
		brd.Set(dest, roundRock)
		brd.Set(each, emptySpace)
		if shared.IsDebugEnabled() {
			shared.Logger.Debug("Move rock.", "from", each, "to", dest)
		}
	})
	if shared.IsDebugEnabled() {
		shared.Logger.Info(
			"Rocks moved.",
			"rock count", rockCount,
			"move count", moveCount,
			"direction", direction)
	}
}

//revive:disable-next-line:cyclomatic
func feedRocks(brd *shared.Board, direction shared.Direction, cb func(shared.Loc)) {
	call := func(loc shared.Loc) {
		if c, ok := brd.Get(loc); ok && c == roundRock {
			cb(loc)
		}
	}
	switch direction {
	case shared.RealEast:
		for y := range brd.GetHeight() {
			for x := brd.GetWidth() - 1; x >= 0; x-- {
				call(shared.Loc{X: x, Y: y})
			}
		}
	case shared.RealSouth:
		for x := range brd.GetWidth() {
			for y := range brd.GetHeight() {
				call(shared.Loc{X: x, Y: y})
			}
		}
	case shared.RealWest:
		for y := range brd.GetHeight() {
			for x := range brd.GetWidth() {
				call(shared.Loc{X: x, Y: y})
			}
		}
	case shared.RealNorth:
		for x := range brd.GetWidth() {
			for y := brd.GetHeight() - 1; y >= 0; y-- {
				call(shared.Loc{X: x, Y: y})
			}
		}
	default:
		shared.Logger.Error("Failure! Not a valid direction.", "direction", direction)
		panic("invalid direction")
	}
}

func hashRocks(brd *shared.Board, pool []int) uint64 {
	var i int
	brd.Iter(func(loc shared.Loc, c rune) bool {
		pool[i] = loc.X
		i++
		pool[i] = loc.Y
		i++
		pool[i] = int(c)
		i++
		return true
	})
	return shared.Hash(pool)
}

func findLoop(track []uint64) (length int, found bool) {
	latest := track[len(track)-1]
	for j := len(track) - 2; j >= 0; j-- {
		if track[j] == latest {
			length = len(track) - 1 - j
			found = true
			return
		}
	}
	return
}
