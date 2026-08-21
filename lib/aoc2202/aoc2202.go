package aoc2202

import (
	"fmt"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

const (
	Rock Piece = iota
	Paper
	Scissors
)

const (
	Lose TargetResult = iota
	Tie
	Win
)

type Piece int
type TargetResult int

func (v Piece) String() string {
	switch v {
	case Rock:
		return "rock"
	case Paper:
		return "paper"
	case Scissors:
		return "scissors"
	default:
		return "unknown Piece"
	}
}

func (v Piece) Wins() Piece {
	switch v {
	case Rock:
		return Scissors
	case Paper:
		return Rock
	default:
		return Paper
	}
}

func (v Piece) LoseTo() Piece {
	switch v {
	case Rock:
		return Paper
	case Paper:
		return Scissors
	default:
		return Rock
	}
}

func (v TargetResult) String() string {
	switch v {
	case Lose:
		return "lose"
	case Tie:
		return "tie"
	case Win:
		return "win"
	default:
		return "unknown TargetResult"
	}
}

func DeriveTotalScore(lines []string, declareEnd bool) int {
	var score int
	for _, each := range lines {
		trimmed := strings.TrimSpace(each)
		if trimmed == "" {
			continue
		}
		pieces := toPieces(strings.Fields(trimmed))
		if len(pieces) != 2 {
			panic(
				fmt.Sprintf(
					"should have 2 pieces on a line, got %d for line \"%s\"",
					len(pieces),
					trimmed))
		}
		myPiece := derivePiece(pieces, declareEnd)
		winScore := deriveWinScore(pieces[0], myPiece)
		pieceScore := myPiece.Score()
		score += pieceScore + winScore
		if shared.IsDebugEnabled() {
			shared.Logger.Debug(
				"Round solved.",
				"pieces", pieces,
				"piece score", pieceScore,
				"win score", winScore)
		}
	}
	shared.Logger.Info("Total score derived.", "score", score)
	return score
}

func derivePiece(pieces []Piece, declareEnd bool) Piece {
	if !declareEnd {
		return pieces[1]
	}
	targetResult := TargetResult(pieces[1])
	switch targetResult {
	case Lose:
		return pieces[0].Wins()
	case Tie:
		return pieces[0]
	case Win:
		return pieces[0].LoseTo()
	default:
		panic(fmt.Sprintf("unknown TargetResult: %d", targetResult))
	}
}

func deriveWinScore(first, second Piece) int {
	loss, even, win := 0, 3, 6
	if first == second {
		return even
	}
	if first.LoseTo() == second {
		return win
	}
	if first.Wins() == second {
		return loss
	}
	panic("impossible state")
}

func toPieces(s []string) []Piece {
	pieces := make([]Piece, len(s))
	for i, each := range s {
		pieces[i] = toPiece(each)
	}
	return pieces
}

func toPiece(s string) Piece {
	switch s {
	case "A", "X":
		return Rock
	case "B", "Y":
		return Paper
	case "C", "Z":
		return Scissors
	default:
		panic("unknown piece: " + s)
	}
}

func (v Piece) Score() int {
	switch v {
	case Rock:
		return 1
	case Paper:
		return 2
	case Scissors:
		return 3
	default:
		panic(fmt.Sprintf("unknown piece: %d", v))
	}
}
