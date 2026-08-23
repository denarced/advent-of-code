package aoc2206

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/denarced/advent-of-code/shared/inr"
	"github.com/denarced/gent"
	"github.com/stretchr/testify/require"
)

func TestDeriveStartOfPackage(t *testing.T) {
	lines, err := inr.ReadPath("testdata/in.txt")
	require.NoError(t, err, "read test data")

	run := func(i int, line string, expected int, seekPackage bool) {
		name := fmt.Sprintf(
			"%d-%s-%d-%s",
			i,
			line[:4],
			expected,
			gent.Tri(seekPackage, "package", "message"),
		)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			req := require.New(t)
			req.Equal(expected, DeriveStartOfPackage(line, seekPackage))
		})
	}

	for i, each := range lines {
		fields := strings.Fields(each)
		expected := gent.OrPanic2(strconv.Atoi(fields[1]))("failed convert expected")
		run(i, fields[0], expected, true)

		if len(fields) == 2 {
			continue
		}
		expected = gent.OrPanic2(strconv.Atoi(fields[2]))("failed to convert expected")
		run(i, fields[0], expected, false)
	}
}
