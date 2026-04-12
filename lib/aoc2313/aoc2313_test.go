package aoc2313

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/denarced/gent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSumReflections(t *testing.T) {
	run := func(fixSmudge bool, expected int) {
		name := gent.Tri(fixSmudge, "fix smudge", "skip smudge")
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			lines, err := inr.ReadPath("testdata/in.txt", inr.IncludeEmpty())
			req.NoError(err, "failed to read test data")

			// EXERCISE
			sum := SumReflections(lines, fixSmudge)

			// VERIFY
			req.Equal(expected, sum)
		})
	}

	run(false, 405)
	run(true, 400)
}

func TestSumReflection(t *testing.T) {
	run := func(name string, lines []string, fixSmudge bool, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			sum := sumReflection(lines, fixSmudge)

			// VERIFY
			req.Equal(expected, sum)
		})
	}

	run(
		"horizontal happy path",
		[]string{
			".#..",
			"..#.",
			"..#.",
			".#..",
			"#.#.",
		},
		false,
		200,
	)
	run(
		"vertical happy path",
		[]string{
			".####.",
			"..##..",
			".####.",
			"..##..",
		},
		false,
		3,
	)
	run(
		"horizontal smudge",
		[]string{
			"..##..",
			"..##..",
			"..#...",
			"......",
		},
		true,
		300,
	)
}

func TestDeriveMaxDistance(t *testing.T) {
	run := func(index, count, expected int) {
		name := fmt.Sprintf("%d+%d==%d", index, count, expected)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			maximus := deriveMaxDistance(index, count)

			// VERIFY
			req.Equal(expected, maximus)
		})
	}

	run(0, 2, 1)
	run(1, 3, 1)
	run(1, 4, 2)
}

type reflectionDto struct {
	lines                        []string
	start, distance, maxDistance int
	vertical                     bool
	expectedSuccess              bool
	// expectedDiffDistance         int
	fixSmudge bool
}

func (v reflectionDto) with(f func(reflectionDto) reflectionDto) reflectionDto {
	return f(v)
}

func withLines(l []string) func(reflectionDto) reflectionDto {
	return func(dto reflectionDto) reflectionDto {
		dto.lines = l
		return dto
	}
}

func withStart(i int) func(reflectionDto) reflectionDto {
	return func(dto reflectionDto) reflectionDto {
		dto.start = i
		return dto
	}
}

func withMaxDistance(i int) func(reflectionDto) reflectionDto {
	return func(dto reflectionDto) reflectionDto {
		dto.maxDistance = i
		return dto
	}
}

func withExpectedSuccess(b bool) func(reflectionDto) reflectionDto {
	return func(dto reflectionDto) reflectionDto {
		dto.expectedSuccess = b
		return dto
	}
}

// func withExpectedDiffDistance(i int) func(reflectionDto) reflectionDto {
// 	return func(dto reflectionDto) reflectionDto {
// 		dto.expectedDiffDistance = i
// 		return dto
// 	}
// }

func withFixSmudge(b bool) func(reflectionDto) reflectionDto {
	return func(dto reflectionDto) reflectionDto {
		dto.fixSmudge = b
		return dto
	}
}

func TestCheckReflection(t *testing.T) {
	run := func(name string, dto reflectionDto) {
		if name == "" {
			name = fmt.Sprintf("%v", dto)
		}
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			ass := assert.New(t)

			ok := checkReflection(
				dto.lines,
				dto.start,
				dto.distance,
				dto.maxDistance,
				dto.vertical,
				dto.fixSmudge,
			)
			ass.Equal(dto.expectedSuccess, ok)
		})
	}

	verLines := []string{
		".##.",
		"#..#",
		"#..#",
		".##.",
	}

	dto := reflectionDto{lines: verLines, distance: 1, vertical: true}
	run(
		"vertical happy path",
		dto.with(withStart(1)).with(withMaxDistance(2)).with(withExpectedSuccess(true)),
	)
	run("vertical without reflection", dto.with(withStart(2)).with(withMaxDistance(1)))

	first := verLines[0]
	modified := first[:3] + "#"
	verLines[0] = modified
	run(
		"vertical fix smudge and diff distance",
		dto.with(withStart(1)).
			with(withMaxDistance(2)).
			with(withFixSmudge(true)).
			with(withExpectedSuccess(true)),
	)
	run(
		"vertical fix smudge without reflection",
		dto.
			with(withLines([]string{
				".##.",
				"###.",
				".###",
				".##.",
			})).
			with(withStart(1)).
			with(withMaxDistance(2)).
			with(withExpectedSuccess(false)).
			with(withFixSmudge(true)),
	)

	horLines := []string{
		".##",
		"#..",
		"#..",
		".##",
	}

	dto = reflectionDto{
		lines:    horLines,
		distance: 1,
		vertical: false,
	}
	run(
		"horizontal happy path",
		dto.with(withStart(1)).with(withMaxDistance(2)).with(withExpectedSuccess(true)),
	)
	run("horizontal without reflection", dto.with(withStart(2)).with(withMaxDistance(1)))

	horLines = []string{
		".#.#",
		"#.##",
		"#.##",
		".#..",
		"####",
	}
	run(
		"horizontal fix smudge and diff distance",
		dto.with(withLines(horLines)).
			with(withStart(1)).
			with(withMaxDistance(2)).
			with(withExpectedSuccess(true)).
			with(withFixSmudge(true)),
	)
	run(
		"horizontal fix smudge without reflection",
		dto.
			with(withLines([]string{
				"..##..",
				".###..",
				"..##..",
				"..###.",
			})).
			with(withStart(1)).
			with(withMaxDistance(2)).
			with(withExpectedSuccess(false)).
			with(withFixSmudge(true)),
	)
}
