package aoc2323

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/denarced/gent"
	"github.com/stretchr/testify/require"
)

func TestFindLongestPath(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	lines, err := inr.ReadPath("testdata/in.txt")
	req.NoError(err, "test data")
	// EXERCISE & VERIFY
	req.Equal(94, FindLongestPath(lines))
}

func BenchmarkFindLongestPath(b *testing.B) {
	shared.InitNullLogging()
	lines := gent.OrPanic2(inr.ReadPath("testdata/in.txt"))("read test data")
	b.ResetTimer()
	for range b.N {
		FindLongestPath(lines)
	}
}
