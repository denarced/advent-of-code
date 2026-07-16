package aoc2323

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/denarced/gent"
	"github.com/stretchr/testify/require"
)

func TestFindLongestPath(t *testing.T) {
	lines := gent.OrPanic2(inr.ReadPath("testdata/in.txt"))("read test data")
	run := func(name string, expected int, f func([]string) int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			pathLength := f(lines)

			// VERIFY
			req.Equal(expected, pathLength)
		})
	}

	run("downhill", 94, FindLongestPath)
	run("downhill and uphill", 154, FindLongestPathWithGraph)
}

func BenchmarkFindLongestPath(b *testing.B) {
	shared.InitNullLogging()
	lines := gent.OrPanic2(inr.ReadPath("testdata/in.txt"))("read test data")
	b.ResetTimer()
	for range b.N {
		FindLongestPath(lines)
	}
}

func BenchmarkFindLongestPathWithGraph(b *testing.B) {
	shared.InitNullLogging()
	lines := gent.OrPanic2(inr.ReadPath("testdata/in.txt"))("read test data")
	b.ResetTimer()
	for range b.N {
		FindLongestPathWithGraph(lines)
	}
}

func TestParseGraph(t *testing.T) {
	lines := []string{
		"#.#######", // #0#######
		"#.....###", // #2.3..###
		"#.#.#.###", // #.#.#.###
		"#.#.#.###", // #.#.#.###
		"#...#...#", // #4.5#...#
		"#.#.###.#", // #.#.###.#
		"#.#.###.#", // #.#.###.#
		"#...###.#", // #...###.#
		"#######.#", // #######1#
	}

	t.Run("parse graph", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)

		// EXERCISE
		aGraph := parseGraph(lines)

		expected := graph{
			edges: []edge{
				{fromIndex: 0, toIndex: 2, length: 1, startDir: shared.RealSouth},
				{fromIndex: 2, toIndex: 3, length: 2, startDir: shared.RealEast},
				{fromIndex: 2, toIndex: 4, length: 3, startDir: shared.RealSouth},
				{fromIndex: 3, toIndex: 1, length: 11, startDir: shared.RealEast},
				{fromIndex: 3, toIndex: 5, length: 3, startDir: shared.RealSouth},
				{fromIndex: 4, toIndex: 5, length: 2, startDir: shared.RealEast},
				{fromIndex: 4, toIndex: 5, length: 8, startDir: shared.RealSouth},
			},
			vertices: []vertice{
				{loc: shared.Loc{X: 1, Y: 8}}, // 0
				{loc: shared.Loc{X: 7, Y: 0}}, // 1
				{loc: shared.Loc{X: 1, Y: 7}}, // 2
				{loc: shared.Loc{X: 3, Y: 7}}, // 3
				{loc: shared.Loc{X: 1, Y: 4}}, // 4
				{loc: shared.Loc{X: 3, Y: 4}}, // 5
			},
		}
		// VERIFY
		req.Equal(expected.edges, aGraph.edges, "edges")
		req.Equal(expected.vertices, aGraph.vertices, "vertices")
	})

	t.Run("miniloop", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)

		// EXERCISE
		result := parseGraph([]string{
			"#.#####", // #0#####
			"#.#...#", // #.#...#
			"#.#.#.#", // #.#.#.#
			"#.....#", // #..23.#
			"####.##", // ####.##
			"####..#", // ####..#
			"#####.#", // #####2#
		})

		expected := graph{
			vertices: []vertice{
				{loc: shared.Loc{X: 1, Y: 6}},
				{loc: shared.Loc{X: 5, Y: 0}},
				{loc: shared.Loc{X: 3, Y: 3}},
				{loc: shared.Loc{X: 4, Y: 3}},
			},
			edges: []edge{
				{fromIndex: 0, toIndex: 2, length: 5, startDir: shared.RealSouth},
				{fromIndex: 2, toIndex: 3, length: 1, startDir: shared.RealEast},
				{fromIndex: 2, toIndex: 3, length: 7, startDir: shared.RealNorth},
				{fromIndex: 3, toIndex: 1, length: 4, startDir: shared.RealSouth},
			},
		}
		// VERIFY
		req.Equal(expected.vertices, result.vertices, "vertices")
		req.Equal(expected.edges, result.edges, "edges")
	})

	t.Run("noname", func(t *testing.T) {
		lines := []string{
			"#.#########",
			"#.#.....###",
			"#...#.#.###",
			"###.#.#.###",
			"###...#...#",
			"#########.#",
		}

		t.Run("parseGraph", func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE
			aGraph := parseGraph(lines)

			expected := graph{
				vertices: []vertice{
					{loc: shared.Loc{X: 1, Y: 5}},
					{loc: shared.Loc{X: 9, Y: 0}},
					{loc: shared.Loc{X: 3, Y: 3}},
					{loc: shared.Loc{X: 5, Y: 4}},
				},
				edges: []edge{
					{fromIndex: 0, toIndex: 2, length: 4, startDir: shared.RealSouth},
					{fromIndex: 3, toIndex: 2, length: 7, startDir: shared.RealSouth},
					{fromIndex: 3, toIndex: 2, length: 3, startDir: shared.RealWest},
					{fromIndex: 3, toIndex: 1, length: 8, startDir: shared.RealEast},
				},
			}
			// VERIFY
			req.Equal(expected.vertices, aGraph.vertices, "vertices")
			req.Equal(expected.edges, aGraph.edges, "edges")
		})

		t.Run("FindLongestPathWithGraph", func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)

			// EXERCISE & VERIFY
			req.Equal(19, FindLongestPathWithGraph(lines))
		})
	})
}

func TestFindLongestPathWithGraph(t *testing.T) {
	t.Run("small board", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)

		lines := []string{
			"#.#######", // #0#######
			"#.....###", // #2.3..###
			"#.#.#.###", // #.#.#.###
			"#.#.#.###", // #.#.#.###
			"#...#...#", // #4.5#...#
			"#.#.###.#", // #.#.###.#
			"#.#.###.#", // #.#.###.#
			"#...###.#", // #...###.#
			"#######.#", // #######1#
		}
		// EXERCISE & VERIFY
		req.Equal(26, FindLongestPathWithGraph(lines))
	})

	t.Run("cross", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)

		// EXERCISE & VERIFY
		req.Equal(
			29,
			FindLongestPathWithGraph([]string{
				"#.######", // #0######
				"#...####", // #2..####
				"#.#.#...", // #.#.#...
				"#.#.#.#.", // #.#.#.#.
				"#.......", // #4.3567.
				"#.##.#.#", // #.##.#.#
				"#.##.#.#", // #.##.#.#
				"#....#.#", // #....#.#
				"######.#", // ######1#
			}))
	})

	t.Run("grid", func(t *testing.T) {
		shared.InitTestLogging(t)
		req := require.New(t)

		// EXERCISE & VERIFY
		req.Equal(
			30,
			FindLongestPathWithGraph([]string{
				"#.#########",
				"#.........#",
				"#.#.#.#.#.#",
				"#.........#",
				"#.#.#.#.#.#",
				"#.........#",
				"#########.#",
			}))
	})
}

func TestExtractDirection(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)

	link := shared.AddLink(nil, shared.Loc{X: 0, Y: 3})
	link = shared.AddLink(link, shared.Loc{X: 1, Y: 3})
	link = shared.AddLink(link, shared.Loc{X: 2, Y: 3})
	link = shared.AddLink(link, shared.Loc{X: 2, Y: 4})
	link = shared.AddLink(link, shared.Loc{X: 2, Y: 5})

	// EXERCISE
	startDir := extractDirection(link, true)
	endDir := extractDirection(link, false)

	req.Equal(shared.RealEast, startDir)
	req.Equal(shared.RealSouth, endDir)
}
