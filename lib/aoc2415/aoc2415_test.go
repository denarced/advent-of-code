package aoc2415

import (
	"fmt"
	"strings"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestCountCoordinateSum(t *testing.T) {
	run := func(name string, lines []string, doubled bool, expected int) {
		suffix := "not doubled"
		if doubled {
			suffix = "doubled"
		}
		t.Run(fmt.Sprintf("%s - %s", name, suffix), func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, CountCoordinateSum(lines, doubled))
		})
	}

	run("empty", []string{}, false, 0)
	run(
		"example",
		[]string{
			"#######",
			"#...O..",
			"#.@....",
			"#......",
			"#......",
		},
		false,
		104)
	run(
		"example",
		[]string{
			"#######",
			"#...#.#",
			"#.....#",
			"#..OO@#",
			"#..O..#",
			"#.....#",
			"#######",
			"",
			"<vv<<^^<<^^",
		},
		true,
		105+207+306)
}

func TestWalk(t *testing.T) {
	run := func(
		name string,
		doubled bool,
		boardLines []string,
		directions string,
		expectedLines []string) {
		suffix := "not doubled"

		if doubled {
			suffix = "doubled"
		}
		t.Run(fmt.Sprintf("%s - %s", name, suffix), func(t *testing.T) {
			shared.InitTestLogging(t)

			brd := shared.NewBoard(dropSpaces(boardLines))
			// EXERCISE
			walk(brd, []rune(directions), doubled)

			// VERIFY
			actual := brd.GetLines()

			ex := strings.Join(dropSpaces(expectedLines), "\n")
			ac := strings.Join(actual, "\n")
			if ex != ac {
				fmt.Println("Expected:")
				fmt.Println(ex)
				fmt.Println("Actual:")
				fmt.Println(ac)
				t.Fail()
			}
		})
	}

	run(
		"basic steps",
		false,
		[]string{
			"....",
			"..@.",
			"....",
		},
		"<<^>>>vv",
		[]string{
			"....",
			"....",
			"...@",
		},
	)
	run(
		"push vertically",
		false,
		[]string{
			".@..",
			".O..",
			".O..",
			"....",
			".#..",
		},
		"vvv",
		[]string{
			"....",
			".@..",
			".O..",
			".O..",
			".#..",
		},
	)
	run(
		"push horizontally",
		false,
		[]string{
			"@O.O.#",
			".....#",
			".....#",
			".....#",
			".....#",
		},
		">>",
		[]string{
			"..@OO#",
			".....#",
			".....#",
			".....#",
			".....#",
		},
	)

	{
		directions := "<^^>>>vv<v>>v<<"
		initial := []string{
			"########",
			"#..O.O.#",
			"##@.O..#",
			"#...O..#",
			"#.#.O..#",
			"#...O..#",
			"#......#",
			"########",
		}
		for i := 0; i < len(directions); i++ {
			dirs := directions[0 : i+1]
			run(
				fmt.Sprintf("example %d %s", i, dirs),
				false,
				append([]string{}, initial...),
				dirs,
				exampleStates()[i],
			)
		}
	}

	run(
		"basic blocking",
		true,
		[]string{
			"#######",
			"..#....",
			"#.@[].#",
			".......",
			"#######",
		},
		"<<v>>^^>>",
		[]string{
			"#######",
			"..#[]..",
			"#....@#",
			".......",
			"#######",
		},
	)

	run(
		"gabriel",
		true,
		[]string{
			/* 7 */ "# # # # # # # #",
			/* 6 */ "# . . . . . . #",
			/* 5 */ "# . . . . . . #",
			/* 4 */ "# [ ] . . . . #",
			/* 3 */ "# . . . . . # #",
			/* 2 */ "# . [ ] . . . #",
			/* 1 */ "# . [ ] . # @ #",
			/* 0 */ "# # # # # # # #",
			//       0 1 2 3 4 5 6 7
		},
		"<<^^<vv<vv>><<^<^<^>^v>>^>>v<v",
		[]string{
			/* 7 */ "# # # # # # # #",
			/* 6 */ "# [ ] . . . . #",
			/* 5 */ "# . . . . @ . #",
			/* 4 */ "# . . . . [ ] #",
			/* 3 */ "# . . . . . # #",
			/* 2 */ "# . . . . . . #",
			/* 1 */ "# [ ] . . # . #",
			/* 0 */ "# # # # # # # #",
			//       0 1 2 3 4 5 6 7
		})
	run(
		"hamud",
		true,
		[]string{
			/* 7 */ "# # # # # # # # #",
			/* 6 */ "# . . . . . . . #",
			/* 5 */ "# . # . . . . . #",
			/* 4 */ "# . . . . [ ] @ #",
			/* 3 */ "# . . [ ] . . . #",
			/* 2 */ "# [ ] . . . . . #",
			/* 1 */ "# . . . . . . . #",
			/* 0 */ "# # # # # # # # #",
			//       0 1 2 3 4 5 6 7 8
		},
		"<<v<v<v<<^",
		[]string{
			/* 7 */ "# # # # # # # # #",
			/* 6 */ "# . . . . . . . #",
			/* 5 */ "# . # [ ] . . . #",
			/* 4 */ "# . [ ] . . . . #",
			/* 3 */ "# [ ] . . . . . #",
			/* 2 */ "# @ . . . . . . #",
			/* 1 */ "# . . . . . . . #",
			/* 0 */ "# # # # # # # # #",
			//       0 1 2 3 4 5 6 7 8
		})
	run(
		"bubble",
		true,
		[]string{
			/* 7 */ "# # # # # # # # #",
			/* 6 */ "# . . . . . . . #",
			/* 5 */ "# . # [ ] . . . #",
			/* 4 */ "# . [ ] . . . . #",
			/* 3 */ "# [ ] . . . . . #",
			/* 2 */ "# @ . . . . . . #",
			/* 1 */ "# . . . . . . . #",
			/* 0 */ "# # # # # # # # #",
			//       0 1 2 3 4 5 6 7 8
		},
		">>^>^>^^<<v>>v<<^^<vv<vvvv",
		[]string{
			/* 7 */ "# # # # # # # # #",
			/* 6 */ "# . . . . . . . #",
			/* 5 */ "# . # . . . . . #",
			/* 4 */ "# @ . . . . . . #",
			/* 3 */ "# [ ] . . . . . #",
			/* 2 */ "# . [ ] . . . . #",
			/* 1 */ "# [ ] . . . . . #",
			/* 0 */ "# # # # # # # # #",
			//       0 1 2 3 4 5 6 7 8
		},
	)
	run(
		"squeeze",
		true,
		[]string{
			/* 3 */ "# # # # # # # # #",
			/* 2 */ "# . [ ] . [ ] @ #",
			/* 1 */ "# [ ] . [ ] . . #",
			/* 0 */ "# # # # # # # # #",
			//       0 1 2 3 4 5 6 7 8
		},
		"<<^v>v<",
		[]string{
			/* 3 */ "# # # # # # # # #",
			/* 2 */ "# [ ] [ ] . . . #",
			/* 1 */ "# [ ] [ ] @ . . #",
			/* 0 */ "# # # # # # # # #",
			//       0 1 2 3 4 5 6 7 8
		},
	)
	run(
		"large example",
		true,
		[]string{
			"####################",
			"##....[]....[]..[]##",
			"##............[]..##",
			"##..[][]....[]..[]##",
			"##....[]@.....[]..##",
			"##[]##....[]......##",
			"##[]....[]....[]..##",
			"##..[][]..[]..[][]##",
			"##........[]......##",
			"####################",
		},
		strings.Join([]string{
			"<vv>^<v^>v>^vv^v>v<>v^v<v<^vv<<<^><<><>>v<vvv<>^v^>^<<<><<v<<<v^vv^v>^",
			"vvv<<^>^v^^><<>>><>^<<><^vv^^<>vvv<>><^^v>^>vv<>v<<<<v<^v>^<^^>>>^<v<v",
			"><>vv>v^v^<>><>>>><^^>vv>v<^^^>>v^v^<^^>v^^>v^<^v>v<>>v^v^<v>v^^<^^vv<",
			"<<v<^>>^^^^>>>v^<>vvv^><v<<<>^^^vv^<vvv>^>v<^^^^v<>^>vvvv><>>v^<<^^^^^",
			"^><^><>>><>^^<<^^v>>><^<v>^<vv>>v>>>^v><>^v><<<<v>>v<v<v>vvv>^<><<>^><",
			"^>><>^v<><^vvv<^^<><v<<<<<><^v<<<><<<^^<v<^^^><^>>^<v^><<<^>>^v<v^v<v^",
			">^>>^v>vv>^<<^v<>><<><<v<<v><>v<^vv<<<>^^v^>^^>>><<^v>>v^v><^^>>^<>vv^",
			"<><^^>^^^<><vvvvv^v<v<<>^v<v>v<<^><<><<><<<^^<<<^<<>><<><^^^>^^<>^>v<>",
			"^^>vv<^v^v<vv>^<><v<^v>^^^>>>^^vvv^>vvv<>>>^<^>>>>>^<<^v>^vvv<>^<><<v>",
			"v^^>>><<^^<>>^v^<v^vv<>v^<<>^<^v^v><^<<<><<^<v><v<>vv>>v><v^<vv<>v^<<^",
		}, ""),
		[]string{
			"####################",
			"##[].......[].[][]##",
			"##[]...........[].##",
			"##[]........[][][]##",
			"##[]......[]....[]##",
			"##..##......[]....##",
			"##..[]............##",
			"##..@......[].[][]##",
			"##......[][]..[]..##",
			"####################",
		})
	run(
		"overwrite",
		true,
		[]string{
			".........",
			"..[]#....",
			"...[]....",
			"....@....",
		},
		"^",
		[]string{
			".........",
			"..[]#....",
			"...[]....",
			"....@....",
		})
}

func exampleStates() [][]string {
	return [][]string{
		{
			"########",
			"#..O.O.#",
			"##@.O..#",
			"#...O..#",
			"#.#.O..#",
			"#...O..#",
			"#......#",
			"########",
		}, {
			"########",
			"#.@O.O.#",
			"##..O..#",
			"#...O..#",
			"#.#.O..#",
			"#...O..#",
			"#......#",
			"########",
		}, {
			"########",
			"#.@O.O.#",
			"##..O..#",
			"#...O..#",
			"#.#.O..#",
			"#...O..#",
			"#......#",
			"########",
		}, {
			"########",
			"#..@OO.#",
			"##..O..#",
			"#...O..#",
			"#.#.O..#",
			"#...O..#",
			"#......#",
			"########",
		}, {
			"########",
			"#...@OO#",
			"##..O..#",
			"#...O..#",
			"#.#.O..#",
			"#...O..#",
			"#......#",
			"########",
		}, {
			"########",
			"#...@OO#",
			"##..O..#",
			"#...O..#",
			"#.#.O..#",
			"#...O..#",
			"#......#",
			"########",
		}, {
			"########",
			"#....OO#",
			"##..@..#",
			"#...O..#",
			"#.#.O..#",
			"#...O..#",
			"#...O..#",
			"########",
		}, {
			"########",
			"#....OO#",
			"##..@..#",
			"#...O..#",
			"#.#.O..#",
			"#...O..#",
			"#...O..#",
			"########",
		}, {
			"########",
			"#....OO#",
			"##.@...#",
			"#...O..#",
			"#.#.O..#",
			"#...O..#",
			"#...O..#",
			"########",
		}, {
			"########",
			"#....OO#",
			"##.....#",
			"#..@O..#",
			"#.#.O..#",
			"#...O..#",
			"#...O..#",
			"########",
		}, {
			"########",
			"#....OO#",
			"##.....#",
			"#...@O.#",
			"#.#.O..#",
			"#...O..#",
			"#...O..#",
			"########",
		}, {
			"########",
			"#....OO#",
			"##.....#",
			"#....@O#",
			"#.#.O..#",
			"#...O..#",
			"#...O..#",
			"########",
		}, {
			"########",
			"#....OO#",
			"##.....#",
			"#.....O#",
			"#.#.O@.#",
			"#...O..#",
			"#...O..#",
			"########",
		}, {
			"########",
			"#....OO#",
			"##.....#",
			"#.....O#",
			"#.#O@..#",
			"#...O..#",
			"#...O..#",
			"########",
		}, {
			"########",
			"#....OO#",
			"##.....#",
			"#.....O#",
			"#.#O@..#",
			"#...O..#",
			"#...O..#",
			"########",
		}}
}

func TestDouble(t *testing.T) {
	shared.InitTestLogging(t)
	require.Equal(
		t,
		[]string{
			"##########",
			"##[][]####",
			"##@...[]##",
			"##......##",
			"##########",
		},
		double([]string{
			"#####",
			"#OO##",
			"#@.O#",
			"#...#",
			"#####",
		}))
}

func dropSpaces(lines []string) []string {
	s := []string{}
	for _, each := range lines {
		s = append(s, strings.ReplaceAll(each, " ", ""))
	}
	return s
}

func TestMoveBoxes(t *testing.T) {
	run := func(name string, lines []string, direction rune, expected []shared.Pair[shared.Loc]) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)

			brd := shared.NewBoard(dropSpaces(lines))
			init := findRobot(brd)
			// EXERCISE
			actual := deriveMovedBoxes(brd, init, direction)

			// VERIFY
			req := require.New(t)
			if expected == nil {
				req.Nil(actual)
			} else {
				req.ElementsMatch(stringifyLocPairs(expected), stringifyLocPairs(actual))
			}
		})
	}

	run(
		"^blocked",
		[]string{
			"#######",
			"#.[]..#",
			"#..@..#",
			"#.....#",
			"#######",
		},
		'^',
		nil)
	run(
		"[]<",
		[]string{
			/* 3 */ "# # # # # # #",
			/* 2 */ "# . . . . . #",
			/* 1 */ "# . [ ] @ . #",
			/* 0 */ "# # # # # # #",
			//       0 1 2 3 4 5 6
		},
		'<',
		[]shared.Pair[shared.Loc]{
			// Robot.
			newMove("4x1 -> 3x1"),
			// Box.
			newMove("3x1 -> 2x1"),
			newMove("2x1 -> 1x1"),
		})
	run(
		">[][]",
		[]string{
			/* 3 */ "# # # # # # # # #",
			/* 2 */ "# . . . . . . . #",
			/* 1 */ "# @ [ ] [ ] . . #",
			/* 0 */ "# # # # # # # # #",
			//       0 1 2 3 4 5 6 7 8
		},
		'>',
		[]shared.Pair[shared.Loc]{
			// Robot.
			newMove("1x1 -> 2x1"),
			// Box.
			newMove("2x1 -> 3x1"),
			newMove("3x1 -> 4x1"),
			newMove("4x1 -> 5x1"),
			newMove("5x1 -> 6x1"),
		})
	run(
		"^T",
		[]string{
			/* 6 */ "# # # # # # # # #",
			/* 5 */ "# # . . . . . . #",
			/* 4 */ "# . . . . . . . #",
			/* 3 */ "# [ ] [ ] . . . #",
			/* 2 */ "# . [ ] . . . . #",
			/* 1 */ "# . . @ . . . . #",
			/* 0 */ "# # # # # # # # #",
			//       0 1 2 3 4 5 6 7 8
		},
		'^',
		[]shared.Pair[shared.Loc]{
			// Robot.
			newMove("3x1 -> 3x2"),
			// First box.
			newMove("2x2 -> 2x3"),
			newMove("3x2 -> 3x3"),
			// First box on upper line.
			newMove("1x3 -> 1x4"),
			newMove("2x3 -> 2x4"),
			// Second box on upper line.
			newMove("3x3 -> 3x4"),
			newMove("4x3 -> 4x4"),
		})
}

func newMove(move string) shared.Pair[shared.Loc] {
	fields := strings.Fields(move)
	from := shared.ParseLoc(fields[0])
	to := shared.ParseLoc(fields[2])
	return shared.NewPair(from, to)
}

func stringifyLocPairs(pairs []shared.Pair[shared.Loc]) []string {
	stringified := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		stringified = append(
			stringified,
			fmt.Sprintf("%s -> %s", pair.First.ToString(), pair.Second.ToString()),
		)
	}
	return stringified
}

func TestCountGps(t *testing.T) {
	run := func(name string, doubled bool, lines []string, expected int) {
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			brd := shared.NewBoard(lines)
			actual := countGps(brd, doubled)
			req.Equal(expected, actual)
		})
	}

	run(
		"example",
		true,
		[]string{
			"####################",
			"##[].......[].[][]##",
			"##[]...........[].##",
			"##[]........[][][]##",
			"##[]......[]....[]##",
			"##..##......[]....##",
			"##..[]............##",
			"##..@......[].[][]##",
			"##......[][]..[]..##",
			"####################",
		},
		9021)
}
