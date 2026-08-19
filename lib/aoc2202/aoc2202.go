package aoc2202

import (
	"fmt"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

func DeriveTotalScore(lines []string) int {
	var score int
	for _, each := range lines {
		trimmed := strings.TrimSpace(each)
		if trimmed == "" {
			continue
		}
		pieces := strings.Fields(trimmed)
		if len(pieces) != 2 {
			panic(
				fmt.Sprintf(
					"should have 2 pieces on a line, got %d for line \"%s\"",
					len(pieces),
					trimmed))
		}
		pieceScore := derivePieceScore(pieces[1])
		winScore := deriveWinScore(pieces[0], pieces[1])
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

func derivePieceScore(piece string) int {
	switch piece {
	case "X":
		return 1
	case "Y":
		return 2
	case "Z":
		return 3
	default:
		panic("unknown piece: " + piece)
	}
}

func deriveWinScore(their, mine string) int {
	first := int(their[0] - 'A')
	second := int(mine[0] - 'X')
	loss, even, win := 0, 3, 6
	if first == second {
		return even
	}
	// 0 rock, 1 paper, 2 scissors
	if first == 0 {
		if second == 1 {
			return win
		}
		return loss
	}
	if first == 1 {
		if second == 0 {
			return loss
		}
		return win
	}
	if second == 0 {
		return win
	}
	return loss
}
