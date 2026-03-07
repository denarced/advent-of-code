package aoc2307

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/denarced/advent-of-code/shared"
)

const (
	card2 card = iota + 2
	card3
	card4
	card5
	card6
	card7
	card8
	card9
	cardT
	cardJ
	cardQ
	cardK
	cardA
)

const (
	handHighCard handType = iota
	handOnePair
	handTwoPair
	handThree
	handFullHouse
	handFour
	handFive
)

type card int
type handType int

type game struct {
	cards    [5]card
	bid      int
	handType handType
}

func CountWinnings(lines []string) int {
	shared.Logger.Info("Count total winnings - start.")
	games := sortGames(parseLines(lines))
	shared.Logger.Info("Games parsed.", "count", len(games))
	var total int
	for i, each := range games {
		product := (i + 1) * each.bid
		shared.Logger.Info("Game proceeds counted.", "game", each, "proceeds", product)
		total += product
	}
	shared.Logger.Info("Total winnings counted.", "total", total)
	return total
}

func sortGames(games []game) []game {
	for i, each := range games {
		each.handType = deriveHandType(each.cards)
		games[i] = each
	}
	slices.SortFunc(games, func(a, b game) int {
		diff := int(a.handType) - int(b.handType)
		if diff != 0 {
			return diff
		}
		for i := range a.cards {
			diff = int(a.cards[i]) - int(b.cards[i])
			if diff != 0 {
				return diff
			}
		}
		return 0
	})
	return games
}

type handResolver func([5]card) (handType, bool)

func deriveHandType(cards [5]card) handType {
	for _, each := range []handResolver{
		isFive,
		isFour,
		isFullHouse,
		isThree,
		isTwoPair,
		isOnePair,
	} {
		if kind, ok := each(cards); ok {
			return kind
		}
	}
	return handHighCard
}

func isFive(cards [5]card) (kind handType, matches bool) {
	kind = handFive
	first := cards[0]
	for _, each := range cards[1:] {
		if first != each {
			return
		}
	}
	matches = true
	return
}

func isFour(cards [5]card) (kind handType, matches bool) {
	kind = handFour
	counts := countCards(cards)
	for _, count := range counts {
		if count == 4 {
			matches = true
			return
		}
	}
	return
}

func isFullHouse(cards [5]card) (kind handType, matches bool) {
	kind = handFullHouse
	counts := countCards(cards)
	var tripFound bool
	var pairFound bool
	for _, count := range counts {
		if count == 3 {
			tripFound = true
		}
		if count == 2 {
			pairFound = true
		}
	}
	matches = tripFound && pairFound
	return
}

func countCards(cards [5]card) map[card]int {
	counts := map[card]int{}
	for _, each := range cards {
		count := counts[each]
		counts[each] = count + 1
	}
	return counts
}

func isThree(cards [5]card) (kind handType, matches bool) {
	kind = handThree
	counts := countCards(cards)
	for _, count := range counts {
		if count == 3 {
			matches = true
			return
		}
	}
	return
}

func isTwoPair(cards [5]card) (kind handType, matches bool) {
	kind = handTwoPair
	counts := countCards(cards)
	var found int
	for _, count := range counts {
		if count == 2 {
			found++
		}
	}
	matches = found == 2
	return
}

func isOnePair(cards [5]card) (kind handType, matches bool) {
	kind = handOnePair
	counts := countCards(cards)
	for _, count := range counts {
		if count == 2 {
			matches = true
			return
		}
	}
	return
}

func parseLines(lines []string) []game {
	games := make([]game, len(lines))
	for i, each := range lines {
		games[i] = parseLine(each)
	}
	return games
}

func parseLine(line string) game {
	pieces := strings.Fields(line)
	if len(pieces) != 2 {
		shared.Logger.Error("Invalid field count on line.", "line", line, "count", len(pieces))
		panic("invalid field count on line")
	}
	hand := pieces[0]
	bid, err := strconv.Atoi(pieces[1])
	if err != nil {
		shared.Logger.Error("Failed to convert bid to int.", "bid", pieces[1], "err", err)
		panic(err)
	}
	if len(hand) != 5 {
		shared.Logger.Error("Invalid count of cards in hand.", "hand", hand)
		panic("invalid count of card in hand")
	}
	return game{
		cards: parseCards(hand),
		bid:   bid,
	}
}

func parseCards(hand string) [5]card {
	var cards [5]card
	for i, each := range []rune(hand) {
		cards[i] = toCard(each)
	}
	return cards
}

func toCard(r rune) card {
	switch unicode.ToUpper(r) {
	case 'T':
		return cardT
	case 'J':
		return cardJ
	case 'Q':
		return cardQ
	case 'K':
		return cardK
	case 'A':
		return cardA
	case '2':
		return card2
	case '3':
		return card3
	case '4':
		return card4
	case '5':
		return card5
	case '6':
		return card6
	case '7':
		return card7
	case '8':
		return card8
	case '9':
		return card9
	default:
		panic(fmt.Sprintf("unknown card: %s", string(r)))
	}
}
