package shared

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func DiffLocationSets(t *testing.T, expected, actual *Set[Location]) {
	stringify := func(l Location) string {
		return l.ToString()
	}
	require.ElementsMatch(
		t,
		MapValues(expected.ToSlice(), stringify),
		MapValues(actual.ToSlice(), stringify),
	)
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
