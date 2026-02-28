package aoc2302

import (
	"fmt"
	"strings"

	"github.com/denarced/advent-of-code/shared"
)

const (
	KindBlue Kind = iota
	KindRed
	KindGreen
)

func DeriveGameCountSum(lines []string, limits map[Kind]int) int {
	games := parseLines(lines)
	var sum int
	for _, each := range games {
		if checkFeasibility(each, limits) {
			sum += each.ID
		}
	}
	return sum
}

func DerivePowerSum(lines []string) int {
	games := parseLines(lines)
	var sum int
	for _, each := range games {
		minimum := deriveMinimum(each)
		sum += multiplyMapValues(minimum)
	}
	return sum
}

type Kind int

type game struct {
	ID   int
	sets []map[Kind]int
}

func parseLines(lines []string) []game {
	games := make([]game, len(lines))
	for i, each := range lines {
		g, err := parseLine(each)
		if err != nil {
			panic(err)
		}
		games[i] = g
	}
	return games
}

func parseLine(line string) (game, error) {
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Parse line.", "line", line)
	}

	linePieces := strings.SplitN(line, ":", 2)
	var gameID int
	rest := linePieces[1]
	_, err := fmt.Sscanf(linePieces[0], "Game %d", &gameID)
	if err != nil {
		return game{}, fmt.Errorf("failed to parse game ID - %w", err)
	}
	if shared.IsDebugEnabled() {
		shared.Logger.Debug("Game parsed.", "ID", gameID, "rest", rest, "line", line)
	}

	setPieces := strings.Split(rest, ";")
	gameSets := []map[Kind]int{}
	for _, eachSet := range setPieces {
		kindPieces := strings.Split(eachSet, ",")
		aGameSet := map[Kind]int{}
		for _, eachKind := range kindPieces {
			var count int
			var kindStr string
			_, err = fmt.Sscanf(strings.TrimSpace(eachKind), "%d %s", &count, &kindStr)
			if err != nil {
				return game{}, fmt.Errorf("failed to scan kind and count - %w", err)
			}
			aGameSet[toKind(kindStr)] = count
		}
		gameSets = append(gameSets, aGameSet)
	}
	return game{ID: gameID, sets: gameSets}, nil
}

func checkFeasibility(g game, limits map[Kind]int) bool {
	for _, each := range g.sets {
		for kind, count := range each {
			if count > limits[kind] {
				return false
			}
		}
	}
	return true
}

func toKind(s string) Kind {
	switch s {
	case "blue":
		return KindBlue
	case "red":
		return KindRed
	case "green":
		return KindGreen
	default:
		panic(fmt.Sprint("No such kind:", s))
	}
}

func deriveMinimum(g game) map[Kind]int {
	vals := map[Kind]int{}
	for _, eachSet := range g.sets {
		for kind, count := range eachSet {
			vals[kind] = max(count, vals[kind])
		}
	}
	return vals
}

func multiplyMapValues[T comparable, V shared.Number](vals map[T]V) V {
	var product V = 1
	for _, each := range vals {
		product *= each
	}
	return product
}
