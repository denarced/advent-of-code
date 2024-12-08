package shared

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToIntTable(t *testing.T) {
	run := func(name string, lines []string, expected [][]int) {
		t.Run(name, func(t *testing.T) {
			InitTestLogging(t)
			req := require.New(t)
			actual := ToIntTable(lines)
			req.Equal(expected, actual)
		})
	}

	run("nil", nil, nil)
	run("empty", []string{}, nil)
	run("1x1", []string{" 33 "}, [][]int{{33}})
	run("3x2", []string{"1 2", "2 3", "3 4"}, [][]int{{1, 2}, {2, 3}, {3, 4}})
}
