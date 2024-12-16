package shared

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
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

type Set[T comparable] struct {
	m map[T]int
}

func NewSet[T comparable](s []T) *Set[T] {
	m := map[T]int{}
	for _, each := range s {
		m[each] = 0
	}
	return &Set[T]{m: m}
}

func (v *Set[T]) Add(item T) {
	v.m[item] = 0
}

func (v *Set[T]) Has(item T) bool {
	_, ok := v.m[item]
	return ok
}

func (v *Set[T]) Count() int {
	return len(v.m)
}

func (v *Set[T]) Copy() *Set[T] {
	m := make(map[T]int, len(v.m))
	for key, val := range v.m {
		m[key] = val
	}
	return &Set[T]{
		m: m,
	}
}

func (v *Set[T]) Iter(cb func(item T) bool) {
	for each := range v.m {
		if !cb(each) {
			break
		}
	}
}

func (v *Set[T]) ToSlice() []T {
	s := make([]T, 0, len(v.m))
	for each := range v.m {
		s = append(s, each)
	}
	return s
}

func (v *Set[T]) Clear() {
	v.m = map[T]int{}
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

func (v Loc) Delta(x, y int) Loc {
	return Loc{X: v.X + x, Y: v.Y + y}
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

func MapValues[T any, U any](s []T, f func(v T) U) []U {
	var result []U
	for _, each := range s {
		result = append(result, f(each))
	}
	return result
}

func FilterValues[T any](s []T, f func(v T) bool) []T {
	var result []T
	for _, each := range s {
		if f(each) {
			result = append(result, each)
		}
	}
	return result
}

func StripPadding(lines []string) []string {
	stripped := make([]string, 0, len(lines))
	for _, each := range lines {
		stripped = append(stripped, strings.ReplaceAll(each, " ", ""))
	}
	return stripped
}

type Board struct {
	lines []string
	MaxX  int
	MaxY  int
}

func NewBoard(lines []string) *Board {
	brd := &Board{lines: lines}
	if len(lines) == 0 {
		return brd
	}
	brd.MaxY = len(lines) - 1
	brd.MaxX = len([]rune(lines[0])) - 1
	Logger.Info("Board created, max coordinates.", "x", brd.MaxX, "y", brd.MaxY)
	return brd
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
	lineCount := len(v.lines)
	for y := 0; y < lineCount; y++ {
		line := v.lines[Abs(y-lineCount+1)]
		runes := []rune(line)
		for x := 0; x < len(runes); x++ {
			if !cb(Loc{X: x, Y: y}, runes[x]) {
				return
			}
		}
	}
}

func (v *Board) Get(loc Loc) (c rune, ok bool) {
	x := loc.X
	if x < 0 || v.MaxX < x {
		return
	}
	y := loc.Y
	if y < 0 || v.MaxY < y {
		return
	}
	line := v.lines[Abs(y-len(v.lines)+1)]
	c = []rune(line)[x]
	ok = true
	return
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
	return (v.MaxX + 1) * (v.MaxY + 1)
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
