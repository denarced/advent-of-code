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

			squad := NewFiringSquad(lines)
			tracker := new(SignalTracker)
			squad.SignalCb = tracker.Add

			// EXERCISE
			squad.Fire()

			// VERIFY
			req.Equal(expected, tracker.LowCount*tracker.HighCount)
		})
	}
	run("in1.txt", 32_000_000)
	run("in2.txt", 11_687_500)
}

func BenchmarkFire(b *testing.B) {
	shared.InitNullLogging()
	req := require.New(b)
	filen := "in2.txt"
	lines, err := inr.ReadPath(filepath.Join("testdata", filen))
	req.NoErrorf(err, "failed to read %s", filen)

	squad := NewFiringSquad(lines)
	b.ResetTimer()
	for range b.N {
		squad.Fire()
	}
}

func TestFindTracked(t *testing.T) {
	req := require.New(t)

	// EXERCISE
	components := map[string][]string{
		"one": {"two"},
		"two": {"three", "four"},
	}
	callers, soughtPulse := FindTracked(components, "one", Low)

	req.Equal([]string{"three", "four"}, callers)
	req.Equal(Low, soughtPulse)
}
