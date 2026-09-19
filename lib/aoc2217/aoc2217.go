package aoc2217

import (
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

const (
	boardWidth = 7
	cellSpace  = ' '
)

var cellChars = []rune{'■', '△', '●', '□', '▲', '○'}

func logBoard(message string, board [][]rune) {
	for i := len(board) - 1; i >= 0; i-- {
		shared.Logger.Debug(
			message,
			"l", string(board[i]))
	}
}

func countBoard(board [][]rune) int {
	var boardTotal int
	for _, each := range board {
		var lineTotal int
		for _, b := range each {
			if b != cellSpace {
				lineTotal++
			}
		}
		boardTotal += lineTotal
	}
	return boardTotal
}

func MeasureTowerHeight(lines []string) int {
	var dirIndex int
	topY := 0
	board := [][]rune{}
	dirs := parseLines(lines)
	pieceGenerator := createPieceGenerator()
	var totalCount int
	for pieceIndex := range 2022 {
		piece := pieceGenerator(topY + 3)
		totalCount += len(piece)
		for {
			currDir := dirs[dirIndex%len(dirs)]
			dirIndex++
			movedPiece, _ := movePiece(board, piece, shared.Loc{X: gent.Tri(currDir == '<', -1, 1)})
			dropped, ok := movePiece(board, movedPiece, shared.Loc{Y: -1})
			if !ok {
				board = addToBoard(board, movedPiece, pieceIndex)
				topY = len(board)
				break
			}
			piece = dropped
		}
	}
	if shared.IsDebugEnabled() {
		logBoard("End.", board)
	}
	if boardTotalCount := countBoard(board); totalCount != boardTotalCount {
		shared.Logger.Error(
			"Atom count mismatch.",
			"atom count in added rocks", totalCount,
			"atom count on board", boardTotalCount)
		panic("atom count mismatch")
	}
	height := len(board)
	shared.Logger.Info("Measure done.", "height", height)
	return height
}

func createPieceGenerator() func(int) []shared.Loc {
	pieces := [][]shared.Loc{
		// #####
		{
			{X: 2, Y: 0},
			{X: 3, Y: 0},
			{X: 4, Y: 0},
			{X: 5, Y: 0},
		},
		//  #
		// ###
		//  #
		{
			{X: 3, Y: 0},
			{X: 2, Y: 1},
			{X: 3, Y: 1},
			{X: 4, Y: 1},
			{X: 3, Y: 2},
		},
		//   #
		//   #
		// ###
		{
			{X: 2, Y: 0},
			{X: 3, Y: 0},
			{X: 4, Y: 0},
			{X: 4, Y: 1},
			{X: 4, Y: 2},
		},
		// #
		// #
		// #
		// #
		{
			{X: 2, Y: 0},
			{X: 2, Y: 1},
			{X: 2, Y: 2},
			{X: 2, Y: 3},
		},
		// ##
		// ##
		{
			{X: 2, Y: 0},
			{X: 3, Y: 0},
			{X: 2, Y: 1},
			{X: 3, Y: 1},
		},
	}
	var pieceIndex int
	return func(y int) []shared.Loc {
		current := pieces[pieceIndex]
		piece := make([]shared.Loc, len(current))
		for i, each := range current {
			piece[i] = each.Delta(shared.Loc{Y: y})
		}
		pieceIndex++
		if pieceIndex >= len(pieces) {
			pieceIndex = 0
		}
		return piece
	}
}

func parseLines(lines []string) []rune {
	for _, each := range lines {
		trimmed := strings.TrimSpace(each)
		if len(trimmed) > 0 {
			return []rune(trimmed)
		}
	}
	return nil
}

func movePiece(board [][]rune, piece []shared.Loc, delta shared.Loc) ([]shared.Loc, bool) {
	movedPiece := make([]shared.Loc, len(piece))
	for i, each := range piece {
		next := each.Delta(delta)
		if !isValid(board, next) {
			return piece, false
		}
		movedPiece[i] = next
	}
	return movedPiece, true
}

func isValid(board [][]rune, loc shared.Loc) bool {
	if loc.X < 0 || boardWidth <= loc.X {
		return false
	}
	if loc.Y < 0 {
		return false
	}
	if loc.Y < len(board) {
		if board[loc.Y][loc.X] != cellSpace {
			return false
		}
	}
	return true
}

func addToBoard(board [][]rune, piece []shared.Loc, pieceIndex int) [][]rune {
	var maxY int
	for _, loc := range piece {
		maxY = max(maxY, loc.Y)
	}
	added := maxY + 1 - len(board)
	for range added {
		board = append(board, createLine())
	}
	cellChar := cellChars[pieceIndex%len(cellChars)]
	for _, loc := range piece {
		if board[loc.Y][loc.X] != cellSpace {
			shared.Logger.Error(
				"Cell already occupied.",
				"char", string(board[loc.Y][loc.X]),
				"to use char", string(cellChar),
				"loc", loc)
			panic("already occupied")
		}
		board[loc.Y][loc.X] = cellChar
	}
	return board
}

func createLine() []rune {
	line := make([]rune, boardWidth)
	for i := range line {
		line[i] = cellSpace
	}
	return line
}
