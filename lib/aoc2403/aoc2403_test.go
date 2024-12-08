package aoc2403

import (
	"fmt"
	"testing"

	"github.com/denarced/advent-of-code/shared"
	"github.com/stretchr/testify/require"
)

func TestMultiply(t *testing.T) {
	run := func(text string, logic bool, expected int) {
		name := fmt.Sprintf("%slogic: %s", shared.Or(logic, "", "!"), text)
		t.Run(name, func(t *testing.T) {
			shared.InitTestLogging(t)
			require.Equal(t, expected, Multiply(text, logic))
		})
	}

	run("empty", false, 0)
	run("empty", true, 0)
	run("#mul(2,3)", false, 6)
	run("--mul(3,4)mul(6,2)", false, 24)
	run("##mul(a,3)-mul(,3)", false, 0)
	run("mul(2,3)do()mul(3,4)don't()mul(5,2)", true, 2*3+3*4)
	run("don't()mulmulmul(2,3)mul(23,3)do()mul(3,4)", true, 3*4)
}
