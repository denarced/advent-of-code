package shared

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
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
