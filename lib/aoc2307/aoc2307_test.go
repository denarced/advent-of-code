package aoc2307

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestCountWinnings(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "failed to read test data")
	req.Equal(6440, CountWinnings(lines))
}

func TestSortGames(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	createGame := func(hand string) game {
		return game{cards: parseCards(hand)}
	}
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
	sorted := sortGames(games)

	withHandType := func(aGame game, kind handType) game {
		aGame.handType = kind
		return aGame
	}
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
