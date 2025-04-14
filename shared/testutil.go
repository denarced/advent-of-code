package shared

import (
	"testing"

	"github.com/denarced/gent"
	"github.com/stretchr/testify/require"
)

func DiffLocSets(t *testing.T, expected, actual *gent.Set[Loc]) {
	stringify := func(l Loc) string {
		return l.ToString()
	}
	require.ElementsMatch(
		t,
		gent.Map(expected.ToSlice(), stringify),
		gent.Map(actual.ToSlice(), stringify),
	)
}
