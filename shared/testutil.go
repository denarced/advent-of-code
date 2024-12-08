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
