package shared

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	RealEast              = Direction{X: 1, Y: 0}
	RealSouth             = Direction{X: 0, Y: -1}
	RealWest              = Direction{X: -1, Y: 0}
	RealNorth             = Direction{X: 0, Y: 1}
	RealPrimaryDirections = []Direction{
		RealEast,
		RealSouth,
		RealWest,
		RealNorth,
	}
	RealSouthEast        = Direction{X: 1, Y: -1}
	RealSouthWest        = Direction{X: -1, Y: -1}
	RealNorthWest        = Direction{X: -1, Y: 1}
	RealNorthEast        = Direction{X: 1, Y: 1}
	RealMiddleDirections = []Direction{
		RealSouthEast,
		RealSouthWest,
		RealNorthWest,
		RealNorthEast,
	}

	DirNorth   = Direction{X: 0, Y: -1}
	DirEast    = Direction{X: 1, Y: 0}
	DirSouth   = Direction{X: 0, Y: 1}
	DirWest    = Direction{X: -1, Y: 0}
	Directions = []Direction{
		DirEast,
		{X: 1, Y: -1},
		DirNorth,
		{X: -1, Y: -1},
		DirWest,
		{X: -1, Y: 1},
		DirSouth,
		{X: 1, Y: 1},
	}
)

func Abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func Die(err error, message string) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Die. Message: \"%s\". Error: %s.\n", message, err)
	//revive:disable-next-line:deep-exit
	os.Exit(2)
}

// If "must" isn't true, panic with message "message".
func Assert(must bool, message string) {
	if must {
		return
	}
	Logger.Error("Drop the axe: Assert failed.", "message", message)
	panic(message)
}

func Or[T any](o bool, yes, no T) T {
	if o {
		return yes
	}
	return no
}

func ReadAll(reader io.Reader) (s string, err error) {
	Logger.Info("Read all.")

	var b []byte
	b, err = io.ReadAll(reader)
	if err != nil {
		return
	}
	s = string(b)
	return
}

func ReadLines(reader io.Reader) (lines []string, err error) {
	Logger.Info("Read lines.")

	var b []byte
	b, err = io.ReadAll(reader)
	if err != nil {
		return
	}

	for _, each := range bytes.Split(b, []byte("\n")) {
		line := strings.TrimSpace(string(each))
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return
}

func ToColumns(s []string) (left []string, right []string) {
	Logger.Info("Split slice content to two coluns.", "length", len(s))
	for _, each := range s {
		pieces := trim(strings.Split(each, " "))
		left = append(left, pieces[0])
		right = append(right, pieces[1])
	}
	return
}

func trim(s []string) (trimmed []string) {
	for _, each := range s {
		aTrimmed := strings.TrimSpace(each)
		if aTrimmed == "" {
			continue
		}
		trimmed = append(trimmed, aTrimmed)
	}
	return
}

type number interface {
	uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | float32 | float64
}

func Max[T number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func ToIntTable(s []string) (table [][]int) {
	for _, each := range s {
		cells := strings.Fields(each)
		var row []int
		for _, c := range cells {
			n, err := strconv.Atoi(c)
			if err != nil {
				panic("Invalid number: " + c)
			}
			row = append(row, n)
		}
		if row == nil {
			panic("Empty row")
		}
		table = append(table, row)
	}
	return
}

type Direction struct {
	X int
	Y int
}

func (v Direction) TurnRight() Direction {
	if v == DirNorth {
		return DirEast
	}
	if v == DirEast {
		return DirSouth
	}
	if v == DirSouth {
		return DirWest
	}
	if v == DirWest {
		return DirNorth
	}
	panic("no direction")
}

func (v Direction) TurnRealRight() Direction {
	switch v {
	case RealEast:
		return RealSouth
	case RealSouth:
		return RealWest
	case RealWest:
		return RealNorth
	case RealNorth:
		return RealEast
	default:
		panic("no direction")
	}
}

type Loc struct {
	X int
	Y int
}

func ParseLoc(s string) Loc {
	var pieces []string
	if strings.Contains(s, " ") {
		pieces = strings.Fields(s)
	} else {
		pieces = strings.Split(s, "x")
	}
	if len(pieces) != 2 {
		panic(fmt.Sprintf("Invalid Loc string: %s.", s))
	}
	ints, err := ToInts(pieces)
	Die(err, "ParseLoc -> ToInts")
	return Loc{X: ints[0], Y: ints[1]}
}

func (v Loc) ToString() string {
	return fmt.Sprintf("%dx%d", v.X, v.Y)
}

func (v Loc) Delta(delta Loc) Loc {
	return Loc{X: v.X + delta.X, Y: v.Y + delta.Y}
}

func (v Loc) Rev() Loc {
	return Loc{X: -v.X, Y: -v.Y}
}

func ToInts(s []string) (nums []int, err error) {
	for _, each := range s {
		var n int
		n, err = strconv.Atoi(each)
		if err != nil {
			Logger.Error("Failed to convert to int.", "string", each, "err", err)
			return
		}
		nums = append(nums, n)
	}
	return
}

func StripPadding(lines []string) []string {
	stripped := make([]string, 0, len(lines))
	for _, each := range lines {
		stripped = append(stripped, strings.ReplaceAll(each, " ", ""))
	}
	return stripped
}

type Board struct {
	ReadOnly bool
	grid     [][]rune
}

func NewBoard(lines []string) *Board {
	brd := new(Board)
	width := 0
	brd.grid = createGrid(lines)
	height := len(brd.grid)
	if height > 0 {
		width = len(brd.grid[0])
	}
	Logger.Info("Board created, max coordinates.", "width", width, "height", height)
	return brd
}

func createGrid(lines []string) [][]rune {
	grid := make([][]rune, len(lines))
	lineCount := len(lines)
	for y := range lineCount {
		line := lines[Abs(y-lineCount+1)]
		grid[y] = []rune(line)
	}
	return grid
}

// Loc is in proper x-y coordinates.
//
// 2|
// 1|
// 0|
// ..---
// ..012
type BoardIterCb func(loc Loc, c rune) bool

func (v *Board) Iter(cb BoardIterCb) {
	for y := range len(v.grid) {
		width := len(v.grid[y])
		for x := range width {
			if !cb(Loc{X: x, Y: y}, v.grid[y][x]) {
				return
			}
		}
	}
}

func (v *Board) Get(loc Loc) (c rune, ok bool) {
	x := loc.X
	y := loc.Y
	if x < 0 || y < 0 {
		return
	}
	if y >= len(v.grid) {
		return
	}
	line := v.grid[y]
	if x >= len(line) {
		return
	}
	c = line[x]
	ok = true
	return
}

func (v *Board) GetOrDie(loc Loc) rune {
	c, ok := v.Get(loc)
	if !ok {
		panic(fmt.Sprintf("Can't find %v.", loc))
	}
	return c
}

func (v *Board) Set(loc Loc, c rune) {
	if v.ReadOnly {
		panic("Board is read-only.")
	}
	y := loc.Y
	if y < 0 || len(v.grid) <= y {
		Logger.Error("Illegal location for Y.", "loc", loc, "c", c)
		panic("Illegal location for Y.")
	}
	line := v.grid[y]
	x := loc.X
	if x < 0 || len(line) <= x {
		Logger.Error("Illegal location for X.", "loc", loc, "c", c)
		panic("Illegal location for X.")
	}
	line[x] = c
}

func (v *Board) NextTo(loc Loc, c rune) []Loc {
	locs := []Loc{}
	for _, xd := range []int{-1, 0, 1} {
		for _, yd := range []int{-1, 0, 1} {
			if xd == 0 && yd == 0 || xd != 0 && yd != 0 {
				continue
			}
			near := Loc{X: loc.X + xd, Y: loc.Y + yd}
			atC, ok := v.Get(near)
			if ok && atC == c {
				locs = append(locs, near)
			}
		}
	}
	return locs
}

func (v *Board) CountArea() int {
	height := len(v.grid)
	if height == 0 {
		return 0
	}
	width := len(v.grid[0])
	return width * height
}

func (v *Board) GetLines() []string {
	lines := []string{}
	for y := len(v.grid) - 1; y >= 0; y-- {
		line := string(v.grid[y])
		lines = append(lines, line)
	}
	return lines
}

func (v *Board) FindOrDie(c rune) Loc {
	var found *Loc
	v.Iter(func(eachLoc Loc, eachC rune) bool {
		if eachC == c {
			found = &eachLoc
			return false
		}
		return true
	})
	if found == nil {
		panic(fmt.Sprintf("Can't find %s.", string(c)))
	}
	return *found
}

func (v *Board) Copy() *Board {
	grid := make([][]rune, 0, len(v.grid))
	for _, line := range v.grid {
		grid = append(grid, append([]rune{}, line...))
	}

	return &Board{
		ReadOnly: v.ReadOnly,
		grid:     grid,
	}
}

func (v *Board) GetWidth() int {
	height := len(v.grid)
	if height == 0 {
		return 0
	}
	return len(v.grid[0])
}

func (v *Board) GetHeight() int {
	return len(v.grid)
}

func Pow(b, e int) int {
	if e == 0 {
		return 1
	}
	if e == 1 {
		return b
	}
	res := b
	for range e - 1 {
		res *= b
	}
	return res
}

type Pair[T any] struct {
	First  T
	Second T
}

func NewPair[T any](first, second T) Pair[T] {
	return Pair[T]{First: first, Second: second}
}

func (v Pair[T]) String() string {
	return fmt.Sprintf("%v-%v", v.First, v.Second)
}

func ParseIntOrDie(s string) int {
	i, err := strconv.Atoi(s)
	Die(err, s)
	return i
}

func OrPanic2[T any](value T, err error) func(msg string) T {
	return func(msg string) T {
		if err == nil {
			return value
		}
		Logger.Error("Expected nil error. Got !nil so panicking.", "message", msg, "err", err)
		panic(msg)
	}
}

func DigitLength(i int) int {
	if i < 0 {
		Logger.Error("Invalid value for DigitLength. Must be >=0.", "value", i)
		panic("Invalid value for DigitLength.")
	}
	if i == 0 {
		return 1
	}
	return int(math.Log10(float64(i))) + 1
}
