package aoc2304

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func countPoints(each card) (points int, count int) {
	each.yours.ForEachAll(func(num int) {
		if each.winners.Has(num) {
			count++
			if points == 0 {
				points++
			} else {
				points *= 2
			}
		}
	})
	return
}

type treeMap struct {
	m    map[int]int
	keys []int
}

func newTreeMap(maxID int) *treeMap {
	keys := make([]int, maxID)
	m := make(map[int]int, maxID)
	for i := range maxID {
		keys[i] = i + 1
		m[i+1] = 1
	}
	return &treeMap{
		m:    m,
		keys: keys,
	}
}

func (v *treeMap) pick() int {
	return v.keys[0]
}

func (v *treeMap) inc(id int) {
	current, ok := v.m[id]
	if !ok {
		panic("new IDs should be impossible")
	}
	v.m[id] = current + 1
}

func (v *treeMap) dec(id int) {
	current, ok := v.m[id]
	if !ok {
		shared.Logger.Error("ID doesn't exist, impossible.", "ID", id)
		panic("illegal state: id must exist")
	}
	current--
	if current > 0 {
		v.m[id] = current
		return
	}

	var index int
	for i, each := range v.keys {
		if each == id {
			index = i
			break
		}
	}
	v.keys = append(v.keys[:index], v.keys[index+1:]...)
	delete(v.m, id)
}

func SumPoints(lines []string, spawn bool) int {
	cards, maxCardID, err := parseLines(lines)
	if err != nil {
		shared.Logger.Error("Failed parse lines.", "err", err)
		panic(err)
	}
	if !spawn {
		var total int
		for _, id := range sortedIntKeys(cards) {
			each := cards[id]
			points, _ := countPoints(each)
			shared.Logger.Info("Card counted.", "ID", each.ID, "points", points)
			total += points
		}
		return total
	}

	links := map[int][]int{}
	for key, card := range cards {
		_, count := countPoints(card)
		targets := make([]int, 0, count)
		for id := key + 1; id <= min(maxCardID, key+count); id++ {
			targets = append(targets, id)
		}
		links[key] = targets
	}

	stack := newTreeMap(maxCardID)
	var total int
	for len(stack.m) > 0 {
		id := stack.pick()
		targets := links[id]
		for _, each := range targets {
			stack.inc(each)
		}
		stack.dec(id)
		total++
	}
	return total
}

type card struct {
	ID      int
	winners *gent.Set[int]
	yours   *gent.Set[int]
}

func parseLines(lines []string) (map[int]card, int, error) {
	cards := map[int]card{}
	var maxID int
	for _, each := range lines {
		pieces := strings.SplitN(each, ":", 2)
		var id int
		_, err := fmt.Sscanf(pieces[0], "Card %d", &id)
		if err != nil {
			return nil, maxID, fmt.Errorf("failed to parse card ID - %w", err)
		}

		sectionPieces := strings.Split(pieces[1], "|")
		if len(sectionPieces) != 2 {
			shared.Logger.Error(
				"Invalid number of card section pieces.",
				"line",
				each,
				"pieces",
				sectionPieces,
			)
			return nil, maxID, errors.New("invalid number of section pieces")
		}
		strings.Fields(sectionPieces[0])
		winners, winnerErr := parseNumbers(sectionPieces[0])
		yours, yourErr := parseNumbers(sectionPieces[1])
		if err := errors.Join(winnerErr, yourErr); err != nil {
			return nil, maxID, fmt.Errorf("failed to parse card numbers - %w", err)
		}
		cards[id] = card{
			ID:      id,
			winners: winners,
			yours:   yours,
		}
		maxID = max(maxID, id)
	}
	return cards, maxID, nil
}

func parseNumbers(s string) (*gent.Set[int], error) {
	fields := strings.Fields(s)
	nums := gent.NewSet[int]()
	for _, each := range fields {
		num, err := strconv.Atoi(each)
		if err != nil {
			return nil, fmt.Errorf("failed to parse a number in a field - %w", err)
		}
		nums.Add(num)
	}
	return nums, nil
}

func sortedIntKeys[T shared.Number, U any](m map[T]U) []T {
	keys := make([]T, 0, len(m))
	for each := range m {
		keys = append(keys, each)
	}
	slices.Sort(keys)
	return keys
}
