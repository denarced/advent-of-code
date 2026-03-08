package aoc2307

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/denarced/gent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountWinnings(t *testing.T) {
	run := func(useJokers bool, expected int) {
		name := gent.Tri(useJokers, "with", "without") + " jokers"
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			lines, err := inr.ReadPath("testdata/in.txt")
			req.NoError(err, "failed to read test data")
			req.Equal(expected, CountWinnings(lines, useJokers))
		})
	}
	run(false, 6440)
	run(true, 5905)
}

func createGame(hand string) game {
	return game{cards: parseCards(hand)}
}
func withHandType(aGame game, kind handType) game {
	aGame.handType = kind
	return aGame
}

func TestSortGamesWithoutJokers(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	pair := createGame("22345")
	highCard := createGame("9TJQK")
	fiveAces := createGame("AAAAA")
	fiveKings := createGame("KKKKK")
	games := []game{
		pair,
		fiveAces,
		fiveKings,
		highCard,
	}

	// EXERCISE
	sorted := sortGames(games, false)

	// VERIFY
	req.Equal(
		[]game{
			withHandType(highCard, handHighCard),
			withHandType(pair, handOnePair),
			withHandType(fiveKings, handFive),
			withHandType(fiveAces, handFive),
		},
		sorted,
	)
}

func TestSortGamesWithJokers(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lowTrips := createGame("J2234")
	highTrips := createGame("22J34")
	games := []game{
		highTrips,
		lowTrips,
	}

	// EXERCISE
	sorted := sortGames(games, true)

	// VERIFY
	req.Equal(
		[]game{
			// Three of a kind starting with J should have lower value.
			withHandType(lowTrips, handThree),
			withHandType(highTrips, handThree),
		},
		sorted,
	)
}

func TestDeriveHighestHand(t *testing.T) {
	run := func(name string, aGame game, expected handType) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			actual := deriveHighestHand(aGame)
			require.Equal(t, fmt.Sprint(expected), fmt.Sprint(actual))
		})
	}
	run(
		"T55J5 to four of a kind",
		game{
			cards:    parseCards("T55J5"),
			handType: handThree,
		},
		handFour)
	run(
		"6JJJJ to five of a kind",
		game{
			cards:    parseCards("6JJJJ"),
			handType: handFour,
		},
		handFive)
	run(
		"2233J to full house",
		game{
			cards:    parseCards("2233J"),
			handType: handTwoPair,
		},
		handFullHouse)
	run(
		"234JJ to three of a kind",
		game{
			cards:    parseCards("234JJ"),
			handType: handOnePair,
		},
		handThree)
	run(
		"2234J to three of a kind",
		game{
			cards:    parseCards("2234J"),
			handType: handOnePair,
		},
		handThree)
	// Doesn't appear to be possible to switch jokers to two pairs. It always turns into trips.
	run(
		"2345J to one pair",
		game{
			cards:    parseCards("2345J"),
			handType: handHighCard,
		},
		handOnePair)
}

func TestCards(t *testing.T) {
	var tests = []struct {
		aCard card
		str   string
	}{
		{card2, "2"},
		{card9, "9"},
		{cardT, "T"},
		{cardA, "A"},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			ass := assert.New(t)
			ass.Equal(tt.str, fmt.Sprint(tt.aCard))
			ass.Equal(tt.aCard, toCard(rune(tt.str[0])))
		})
	}

	ass := assert.New(t)
	for _, each := range []int{-1, 0, 1, 15} {
		ass.Equal("unknown", fmt.Sprint(card(each)))
	}
}
