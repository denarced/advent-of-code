package aoc2506

import (
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/stretchr/testify/require"
)

func TestCalculate(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines, err := inr.ReadPath("testdata/in.txt", inr.NoTrim())
	req.NoError(err, "read test data")
	require.Equal(t, int64(4_277_556), Calculate(lines, true))
	require.Equal(t, int64(3_263_827), Calculate(lines, false))
}

func TestSplitByEmptyColumn(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	lines := []string{
		"123 234 345",
		" 45  56  67",
		"  6   7   8",
		"+   *   +  ",
	}
	table := splitByEmptyColumns(lines)
	req.Equal(
		[][]string{
			{
				"123",
				" 45",
				"  6",
				"+  ",
			},
			{
				"234",
				" 56",
				"  7",
				"*  ",
			},
			{
				"345",
				" 67",
				"  8",
				"+  ",
			},
		},
		table)
}

func TestSplitWithIndexes(t *testing.T) {
	shared.InitTestLogging(t)
	req := require.New(t)
	req.Equal(
		[]string{"0123", "5", "789a"},
		splitWithIndexes(
			"0123456789a",
			[]int{4, 6},
		),
	)
}

func TestPivot(t *testing.T) {
	result := pivot([]string{
		"123",
		" 45",
		"  6",
		"  7",
	})
	require.Equal(
		t,
		[]string{
			"3567",
			"24  ",
			"1   ",
		},
		result)
}
