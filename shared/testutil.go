package shared

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func DiffLocSets(t *testing.T, expected, actual *Set[Loc]) {
	stringify := func(l Loc) string {
		return l.ToString()
	}
	require.ElementsMatch(
		t,
		MapValues(expected.ToSlice(), stringify),
		MapValues(actual.ToSlice(), stringify),
	)
}
