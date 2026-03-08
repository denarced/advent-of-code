package aoc2307

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

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

func (v handType) String() string {
	switch v {
	case handHighCard:
		return "highCard"
	case handOnePair:
		return "onePair"
	case handTwoPair:
		return "twoPair"
	case handThree:
		return "threeOfKind"
	case handFullHouse:
		return "fullHouse"
	case handFour:
		return "fourOfKind"
	case handFive:
		return "fiveOfKind"
	default:
		return "unknown"
	}
}

type game struct {
	cards    [5]card
	bid      int
	handType handType
}

func CountWinnings(lines []string, useJokers bool) int {
	shared.Logger.Info("Count total winnings - start.")
	games := sortGames(parseLines(lines), useJokers)
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

func countJokers(cards [5]card) int {
	var count int
	for _, each := range cards {
		if each == cardJ {
			count++
		}
	}
	return count
}

func convertCardToInt(aCard card, useJokers bool) int {
	nominal := int(aCard)
	if !useJokers || aCard != cardJ {
		return nominal
	}
	return 1
}

func sortGames(games []game, useJokers bool) []game {
	for i, each := range games {
		each.handType = deriveHandType(each.cards)
		if useJokers && countJokers(each.cards) > 0 {
			each.handType = deriveHighestHand(each)
		}
		games[i] = each
	}
	slices.SortFunc(games, func(a, b game) int {
		diff := int(a.handType) - int(b.handType)
		if diff != 0 {
			return diff
		}
		for i := range a.cards {
			diff = convertCardToInt(a.cards[i], useJokers) - convertCardToInt(b.cards[i], useJokers)
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
	if len(hand) != 5 {
		shared.Logger.Error("Invalid hand.", "hand", hand)
		panic("invalid hand")
	}
	var cards [5]card
	for i, each := range []rune(hand) {
		cards[i] = toCard(each)
	}
	return cards
}

func toCard(r rune) card {
	letters := []rune("23456789TJQKA")
	cards := []card{
		card2,
		card3,
		card4,
		card5,
		card6,
		card7,
		card8,
		card9,
		cardT,
		cardJ,
		cardQ,
		cardK,
		cardA,
	}
	for i, each := range letters {
		if each == r {
			return cards[i]
		}
	}
	panic(fmt.Sprintf("unknown card: %s", string(r)))
}

func (v card) String() string {
	intVal := int(v)
	if 2 <= intVal && intVal <= 9 {
		return fmt.Sprint(intVal)
	}
	switch v {
	case cardT:
		return "T"
	case cardJ:
		return "J"
	case cardQ:
		return "Q"
	case cardK:
		return "K"
	case cardA:
		return "A"
	default:
		return "unknown"
	}
}

func deriveHighestHand(aGame game) handType {
	jokerCount := countJokers(aGame.cards)
	commonCount := countMostCommonNonJoker(aGame.cards)
	combinedCount := jokerCount + commonCount
	if combinedCount == 5 {
		return handFive
	}
	if combinedCount == 4 {
		return handFour
	}
	if len(countCards(aGame.cards)) == 3 {
		return handFullHouse
	}
	if combinedCount == 3 {
		return handThree
	}
	if jokerCount == 1 && commonCount == 1 {
		return handOnePair
	}
	return aGame.handType
}

func countMostCommonNonJoker(cards [5]card) int {
	counts := countCards(cards)
	var maximus int
	for c, each := range counts {
		if c != cardJ {
			maximus = max(maximus, each)
		}
	}
	return maximus
}
