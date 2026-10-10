package aoc2217

import (
	"math"
	"slices"
	"strings"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/gent"
)

const (
	boardWidth = 7
	cellSpace  = ' '
	cellChar   = '#'
)

func MeasureTowerHeight(lines []string, rockCount int) int {
	shared.Logger.Info("Measure tower.", "rock count", rockCount)
	var dirIndex int
	board := createPieceBoard()
	dirs := parseLines(lines)
	pieceGenerator := createPieceGenerator()
	var totalCount int
	hive := createHiveMind()
	for pieceIndex := 0; pieceIndex < rockCount; pieceIndex++ {
		pieceID, piece := pieceGenerator(board.getHeight() + 3)
		totalCount += len(piece)
		for {
			currDir := dirs[dirIndex%len(dirs)]
			dirIndex++
			movedPiece, _ := movePiece(board, piece, shared.Loc{X: gent.Tri(currDir == '<', -1, 1)})
			dropped, ok := movePiece(board, movedPiece, shared.Loc{Y: -1})
			if !ok {
				board = addToBoard(board, movedPiece)
				break
			}
			piece = dropped
		}
		if hive == nil {
			continue
		}

		board.buryHidden()
		shape := board.feelShape()
		if existing, ok := hive.exists(shape, shapeSpec{
			ID:         pieceID,
			dirIndex:   dirIndex % len(dirs),
			pieceIndex: pieceIndex,
			maxY:       board.getMaxY(),
		}); ok {
			currMaxY := board.getMaxY()
			yIncrease := currMaxY - existing.maxY
			pieceIndexInc := pieceIndex - existing.pieceIndex
			var totalY int
			var lastI int
			for i := pieceIndex + pieceIndexInc; i < rockCount; i += pieceIndexInc {
				totalY += yIncrease
				lastI = i
			}
			if lastI > 0 {
				shared.Logger.Info(
					"Done warping.",
					"piece index", pieceIndex,
					"current max Y", currMaxY,
					"new piece index", lastI,
					"new max Y", currMaxY+totalY)
				pieceIndex = lastI
				board.setExtraHeight(totalY)
			}
			hive = nil
		}
	}
	height := board.getHeight()
	shared.Logger.Info("Measure done.", "height", height)
	return height
}

func createPieceGenerator() func(int) (int, []shared.Loc) {
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
	return func(y int) (int, []shared.Loc) {
		current := pieces[pieceIndex]
		piece := make([]shared.Loc, len(current))
		for i, each := range current {
			piece[i] = each.Delta(shared.Loc{Y: y})
		}
		currentIndex := pieceIndex
		pieceIndex++
		if pieceIndex >= len(pieces) {
			pieceIndex = 0
		}
		return currentIndex, piece
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

func movePiece(board *pieceBoard, piece []shared.Loc, delta shared.Loc) ([]shared.Loc, bool) {
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

func isValid(board *pieceBoard, loc shared.Loc) bool {
	if loc.X < 0 || boardWidth <= loc.X {
		return false
	}
	if loc.Y < 0 {
		return false
	}
	if loc.Y < board.getHeight() {
		if board.get(loc) != cellSpace {
			return false
		}
	}
	return true
}

func addToBoard(board *pieceBoard, piece []shared.Loc) *pieceBoard {
	var maxY int
	for _, loc := range piece {
		maxY = max(maxY, loc.Y)
	}
	added := maxY + 1 - board.getHeight()
	for range added {
		board.appendLine(createLine())
	}
	board.set(piece)
	return board
}

func createLine() []rune {
	line := make([]rune, boardWidth)
	for i := range line {
		line[i] = cellSpace
	}
	return line
}

type pieceBoard struct {
	xToWhys []map[int]struct{}

	wanderBook  *guestbook
	visibleBook *guestbook
}

type guestbook struct {
	exes [][]int
}

func createGuestbook() *guestbook {
	return &guestbook{exes: make([][]int, boardWidth)}
}

func (v *guestbook) clear() {
	for i := range v.exes {
		v.exes[i] = v.exes[i][:0]
	}
}

func (v *guestbook) add(loc shared.Loc) {
	i, exists := slices.BinarySearch(v.exes[loc.X], loc.Y)
	if exists {
		return
	}
	v.exes[loc.X] = slices.Insert(v.exes[loc.X], i, loc.Y)
}

func (v *guestbook) has(loc shared.Loc) bool {
	_, ok := slices.BinarySearch(v.exes[loc.X], loc.Y)
	return ok
}

func createPieceBoard() *pieceBoard {
	board := new(pieceBoard)
	board.xToWhys = make([]map[int]struct{}, boardWidth)
	board.wanderBook = createGuestbook()
	board.visibleBook = createGuestbook()
	for i := range boardWidth {
		board.xToWhys[i] = map[int]struct{}{-1: struct{}{}}
		board.wanderBook.exes[i] = []int{}
		board.visibleBook.exes[i] = []int{}
	}
	return board
}

func (v *pieceBoard) buryHidden() {
	v.wanderBook.clear()
	v.visibleBook.clear()
	v.wander()
	v.findDropped()
}

func (v *pieceBoard) wander() {
	wanderers := []shared.Loc{{X: boardWidth / 2, Y: v.getHeight()}}
	for len(wanderers) > 0 {
		curr := wanderers[0]
		wanderers = wanderers[1:]
		v.wanderBook.add(curr)
		if v.get(curr) != cellSpace {
			v.visibleBook.add(curr)
		} else {
			for _, delta := range []shared.Loc{{X: -1}, {X: 1}, {Y: -1}} {
				each := curr.Delta(delta)
				if each.X < 0 ||
					boardWidth <= each.X ||
					v.visibleBook.has(each) ||
					v.wanderBook.has(each) {
					continue
				}
				wanderers = append(wanderers, each)
			}
		}
	}
}

func (v *pieceBoard) findDropped() {
	for x, whys := range v.xToWhys {
		var toRemove []int
		for y := range whys {
			cand := shared.Loc{X: x, Y: y}
			if !v.visibleBook.has(cand) {
				toRemove = append(toRemove, y)
			}
		}
		for _, each := range toRemove {
			delete(whys, each)
		}
	}
}

func (v *pieceBoard) feelShape() uint64 {
	yToX := map[int][]int{}
	for x, whys := range v.xToWhys {
		for y := range whys {
			if yToX[y] == nil {
				yToX[y] = make([]int, boardWidth)
			}
			yToX[y][x] = shared.Pow(2, x)
		}
	}
	whys := make([]int, len(yToX))
	var i int
	for y, exes := range yToX {
		slices.Sort(exes)
		whys[i] = y
		i++
	}
	slices.Sort(whys)
	signs := make([]int, len(whys))
	i = 0
	for _, y := range whys {
		var signature int
		for _, sprite := range yToX[y] {
			signature += sprite
		}
		signs[i] = signature
		i++
	}
	return shared.Hash(signs)
}

func (v *pieceBoard) getHeight() int {
	return v.getMaxY() + 1
}

func (v *pieceBoard) getMaxY() int {
	maximus := math.MinInt
	for _, whys := range v.xToWhys {
		for y := range whys {
			maximus = max(maximus, y)
		}
	}
	return maximus
}

func (v *pieceBoard) setExtraHeight(extra int) {
	for _, whys := range v.xToWhys {
		sorted := make([]int, 0, len(whys))
		for y := range whys {
			sorted = append(sorted, y)
		}
		slices.Sort(sorted)
		clear(whys)
		for i := len(sorted) - 1; i >= 0; i-- {
			whys[sorted[i]+extra] = struct{}{}
		}
	}
}

func (v *pieceBoard) get(loc shared.Loc) rune {
	_, ok := v.xToWhys[loc.X][loc.Y]
	if ok {
		return cellChar
	}
	return cellSpace
}

func (v *pieceBoard) set(added []shared.Loc) {
	for _, atom := range added {
		v.xToWhys[atom.X][atom.Y] = struct{}{}
	}
}

func (*pieceBoard) appendLine([]rune) {}

type shapeSpec struct {
	ID         int
	dirIndex   int
	pieceIndex int
	maxY       int
}

type hiveMind struct {
	seen map[uint64]shapeSpec
}

func createHiveMind() *hiveMind {
	return &hiveMind{seen: map[uint64]shapeSpec{}}
}

func (v *hiveMind) exists(shape uint64, spec shapeSpec) (shapeSpec, bool) {
	existing, ok := v.seen[shape]
	v.seen[shape] = spec
	if !ok {
		return shapeSpec{}, false
	}
	return existing, existing.dirIndex == spec.dirIndex && existing.ID == spec.ID
}
