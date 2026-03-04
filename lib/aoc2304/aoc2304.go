package aoc2304

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

func SumPoints(lines []string) int {
	cards, err := parseLines(lines)
	if err != nil {
		shared.Logger.Error("Failed parse lines.", "err", err)
		panic(err)
	}
	grandTotal := 0
	for _, each := range cards {
		points := 0
		shared.Logger.Debug("Resolving a card.", "ID", each.ID)
		each.yours.ForEachAll(func(num int) {
			if each.winners.Has(num) {
				if points == 0 {
					points++
				} else {
					points *= 2
				}
				shared.Logger.Debug("Match found.", "num", num, "total", points)
			}
		})
		shared.Logger.Info("Card counted.", "total", points, "ID", each.ID)
		grandTotal += points
	}
	return grandTotal
}

type card struct {
	ID      int
	winners *gent.Set[int]
	yours   *gent.Set[int]
}

func parseLines(lines []string) ([]card, error) {
	var cards []card
	for _, each := range lines {
		pieces := strings.SplitN(each, ":", 2)
		var id int
		_, err := fmt.Sscanf(pieces[0], "Card %d", &id)
		if err != nil {
			return nil, fmt.Errorf("failed to parse card ID - %w", err)
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
			return nil, errors.New("invalid number of section pieces")
		}
		strings.Fields(sectionPieces[0])
		winners, winnerErr := parseNumbers(sectionPieces[0])
		yours, yourErr := parseNumbers(sectionPieces[1])
		if err := errors.Join(winnerErr, yourErr); err != nil {
			return nil, fmt.Errorf("failed to parse card numbers - %w", err)
		}
		cards = append(cards, card{
			ID:      id,
			winners: winners,
			yours:   yours,
		})
	}
	return cards, nil
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
