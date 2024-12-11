package aoc2410

import (
	"fmt"
	"sync"

	"github.com/denarced/advent-of-code/shared"
)

func DeriveSumOfTrailheadScores(lines []string, ratings bool) int {
	if len(lines) == 0 {
		return 0
	}
	ch := make(chan trail)
	resCh := make(chan int)
	go func() {
		trails := shared.NewSet([]trail{})
		count := 0
		for each := range ch {
			trails.Add(each)
			count++
		}
		if ratings {
			resCh <- count
		} else {
			resCh <- trails.Count()
		}
		close(resCh)
	}()
	var wg sync.WaitGroup
	brd := shared.NewBoard(lines)
	brd.Iter(func(loc shared.Loc, c rune) {
		if c == '0' {
			wg.Add(1)
			go blaze(loc, loc, &wg, brd, ch)
		}
	})
	wg.Wait()
	close(ch)
	return <-resCh
}

type trail struct {
	start shared.Loc
	end   shared.Loc
}

func blaze(startLoc, currLoc shared.Loc, wg *sync.WaitGroup, brd *shared.Board, ch chan<- trail) {
	defer wg.Done()

	current, ok := brd.Get(currLoc)
	if !ok {
		panic(fmt.Sprintf("!ok should've been impossible: %v.", currLoc))
	}
	if current == '9' {
		t := trail{start: startLoc, end: currLoc}
		shared.Logger.Debug("Found trail.", "trail", t)
		ch <- t
		return
	}
	near := brd.NextTo(currLoc, increment(current))
	for _, each := range near {
		wg.Add(1)
		go blaze(startLoc, each, wg, brd, ch)
	}
}

func increment(c rune) rune {
	return c + 1
}
