package aoc2320

import (
	"path/filepath"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestCountSignalProductFromLines(t *testing.T) {
	run := func(filen string, expected int) {
		t.Run(filen, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			lines, err := inr.ReadPath(filepath.Join("testdata", filen))
			req.NoErrorf(err, "failed to read %s", filen)
			// EXERCISE & VERIFY
			req.Equal(expected, CountSignalProductFromLines(lines))
		})
	}
	run("in1.txt", 32_000_000)
	run("in2.txt", 11_687_500)
}

func BenchmarkCountSignalProductFromLines(b *testing.B) {
	shared.InitNullLogging()
	req := require.New(b)
	filen := "in2.txt"
	lines, err := inr.ReadPath(filepath.Join("testdata", filen))
	req.NoErrorf(err, "failed to read %s", filen)

	b.ResetTimer()
	for range b.N {
		CountSignalProductFromLines(lines)
	}
}
